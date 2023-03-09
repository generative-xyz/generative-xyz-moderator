package helpers

import (
	"crypto/md5"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/big"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/utils"
	"rederinghub.io/utils/identicon"
)

func GenerateMd5String(input string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(input)))
}

func ToDoc(v interface{}) (*bson.D, error) {
	data, err := bson.Marshal(v)
	if err != nil {
		return nil, err
	}

	doc := &bson.D{}
	err = bson.Unmarshal(data, doc)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func Transform(from interface{}, to interface{}) error {
	bytes, err := bson.Marshal(from)
	if err != nil {
		return err
	}

	err = bson.Unmarshal(bytes, to)
	if err != nil {
		return err
	}

	return nil
}

func JsonTransform(from interface{}, to interface{}) error {
	bytes, err := json.Marshal(from)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, to)
	if err != nil {
		return err
	}

	return nil
}

func GenerateKey(key string) string {
	key = strings.ToUpper(key)
	key = strings.ReplaceAll(key, " ", "_")
	return key
}

func GenerateSlug(key string) string {
	key = strings.ReplaceAll(key, " ", "-")
	key = strings.ReplaceAll(key, "#", "-")
	key = strings.ReplaceAll(key, "@", "")
	key = strings.ReplaceAll(key, `%`, "")
	key = strings.ReplaceAll(key, `?`, "")
	key = strings.ReplaceAll(key, `(`, "")
	key = strings.ReplaceAll(key, `)`, "")
	key = strings.ReplaceAll(key, `[`, "")
	key = strings.ReplaceAll(key, `]`, "")
	key = strings.ReplaceAll(key, `{`, "")
	key = strings.ReplaceAll(key, `}`, "")
	key = strings.ReplaceAll(key, `!`, "")
	key = strings.ReplaceAll(key, `=`, "")
	//key = regexp.MustCompile(`[^a-zA-Z0-9?:-]+`).ReplaceAllString(key, "")
	key = strings.ToLower(key)
	key = ReplaceNonUTF8(key)
	return key
}

func ReplaceNonUTF8(filename string) string  {
	re := regexp.MustCompile("[^a-zA-Z0-9./:]")
	return fmt.Sprintf(re.ReplaceAllString(filename, ""))
}

func Base64Decode(base64Str string) ([]byte, error) {
	sDec, err := b64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		return nil, err
	}
	return sDec, nil
}

func Base64DecodeRaw(base64Str string, object interface{}) error {
	base64Str = strings.ReplaceAll(base64Str, "data:application/json;base64,", "")
	sDec, err := b64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		return err
	}

	err = json.Unmarshal(sDec, &object)
	if err != nil {
		return err
	}

	return nil
}

func Base64Encode(data []byte) string {
	sDec := b64.StdEncoding.EncodeToString(data)
	return sDec
}

func ReplaceToken(token string) string {
	token = strings.ReplaceAll(token, "Bearer", "")
	token = strings.ReplaceAll(token, "bearer", "")
	token = strings.ReplaceAll(token, " ", "")
	return token
}

func HexaNumberToInteger(hexaString string) string {
	// replace 0x or 0X with empty String
	numberStr := strings.Replace(hexaString, "0x", "", -1)
	numberStr = strings.Replace(numberStr, "0X", "", -1)
	return numberStr
}

func CreateIcon(name *string) string {
	return identicon.CreateIcon(name)
}

func CreateMqttTopic(ordAddress string) string {
	return fmt.Sprintf("btc_mint_adderss_%s", ordAddress)
}

func CreateProfileLink(walletAdress string, displayName string) string {
	name := walletAdress
	if displayName != "" {
		name = displayName
	}
	link := fmt.Sprintf("%s/profile/%s", os.Getenv("DOMAIN"), walletAdress)
	return fmt.Sprintf("<%s|%s>", link, name)
}

func CreateTokenLink(projectID string, tokenID string, tokenName string) string {
	link := fmt.Sprintf("%s/generative/%s/%s", os.Getenv("DOMAIN"), projectID, tokenID)
	return fmt.Sprintf("<%s|%s>", link, tokenName)
}


func CreateTokenImageLink(url string) string {
	link := fmt.Sprintf("%s", url)
	return fmt.Sprintf("<%s|%s>", link, "Review")
}


func CreateProjectLink(projectID string, priojectName string) string {
	link := fmt.Sprintf("%s/generative/%s", os.Getenv("DOMAIN"), projectID)
	return fmt.Sprintf("<%s|%s>", link, priojectName)
}

func ParseBigToFloat(number *big.Int) float64 {
	numStr := number.String()

	n, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		return 0
	}
	return n
}

func ParseUintToUnixTime(number uint64) *time.Time {
	t := time.Unix(int64(number), 0)
	return &t
}

func CreateBTCOrdWallet(userWallet string) string {
	now := time.Now().UTC().Unix()
	return fmt.Sprintf("%s_%s_%d", "USER", userWallet, now)
}

func GetExternalPrice(tokenSymbol string) (float64, error) {
	binanceAPI := os.Getenv("BINANCE_API")
	if binanceAPI == "" {
		binanceAPI = "https://api.binance.com"
	}
	binancePriceURL := fmt.Sprintf("%v/api/v3/ticker/price?symbol=", binanceAPI)
	var price struct {
		Symbol string `json:"symbol"`
		Price  string `json:"price"`
	}
	var jsonErr struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}
	retryTimes := 0
retry:
	retryTimes++
	if retryTimes > 2 {
		return 0, nil
	}
	tk := strings.ToUpper(tokenSymbol)
	resp, err := http.Get(binancePriceURL + tk + "USDT")
	if err != nil {
		log.Println(err)
		goto retry
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	err = json.Unmarshal(body, &price)
	if err != nil {
		err = json.Unmarshal(body, &jsonErr)
		if err != nil {
			log.Println(err)
			goto retry
		}
	}
	resp.Body.Close()
	value, err := strconv.ParseFloat(price.Price, 32)
	if err != nil {
		log.Println("getExternalPrice", tokenSymbol, err)
		return 0, nil
	}
	return value, nil
}

func CalcOrigBinaryLength(datas string) int {

	l := len(datas)

	// count how many trailing '=' there are (if any)
	eq := 0
	if l >= 2 {
		if datas[l-1] == '=' {
			eq++
		}
		if datas[l-2] == '=' {
			eq++
		}

		l -= eq
	}

	// basically:
	//
	// eq == 0 :    bits-wasted = 0
	// eq == 1 :    bits-wasted = 2
	// eq == 2 :    bits-wasted = 4

	// each base64 character = 6 bits

	// so orig length ==  (l*6 - eq*2) / 8

	return (l*3 - eq) / 4
}

func SliceStringContains(slice []string, target string) bool {
	for _, e := range slice {
		if e == target {
			return true
		}
	}
	return false
}

func StringToBTCAmount(price string) *big.Float {
	pow := math.Pow10(8)
	powBig := big.NewFloat(0).SetFloat64(pow)

	mintPrice := big.NewFloat(0)
	mintPrice.SetString(price)
	mintPrice.Mul(mintPrice, powBig)

	return mintPrice
}

func CalculateRefEarning(amount float64, percent int32) (string, string) {
	artist1Earning := amount * float64(percent)/1000
	generativeEarning :=  amount * float64(1000 - utils.PERCENT_EARNING)/1000
	referralEarning := artist1Earning * float64(percent)/1000
	return fmt.Sprintf("%d",int(referralEarning)), fmt.Sprintf("%d",int(generativeEarning))
}

func CalculateVolumEarning(amount float64, percent int32) (string, string) {
	rate := float64(percent) / 1000
	artist1Earning := amount * rate
	generativeEarning :=  amount - artist1Earning
	return fmt.Sprintf("%d",int(artist1Earning)), fmt.Sprintf("%d",int(generativeEarning))
}

func FileExtension(fileName string) string {
	str := strings.Split(fileName, "/")
	last := ""
	if len(str) > 0 {
		last = str[len(str) - 1]
	}

	lastArr := strings.Split(last, ".")
	if len(lastArr) > 0 {
		last = lastArr[len(lastArr) - 1]
	}
	return last
}
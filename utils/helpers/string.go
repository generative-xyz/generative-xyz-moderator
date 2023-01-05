package helpers

import (
	"crypto/md5"
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"rederinghub.io/utils"
	"rederinghub.io/utils/identicon"

	"go.mongodb.org/mongo-driver/bson"
)

func GenerateMd5String(input string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(input)))
}

func ToDoc (v interface{}) (*bson.D, error) {
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
		return   err
	}

	err = bson.Unmarshal(bytes, to)
	if err != nil {
		return   err
	}

	return nil
}

func JsonTransform(from interface{}, to interface{}) error {
	bytes, err := json.Marshal(from)
	if err != nil {
		return   err
	}

	err = json.Unmarshal(bytes, to)
	if err != nil {
		return   err
	}

	return nil
}

func ParseCache(cached *string, resp interface{}) error {
	if cached == nil {
		return errors.New("Cached data is empty")
	}

	bytes := []byte(*cached)
	err := json.Unmarshal(bytes, &resp)
	if err != nil  {
		return err
	}
	return nil
}
 
func GenerateKey(key string) string {
	key = strings.ToUpper(key)
	key = strings.ReplaceAll(key, " ", "_")
	return key
}

func GenerateCachedProfileKey(accessToken string) string {
	return fmt.Sprintf("%s_%s", utils.REDIS_PROFILE, GenerateMd5String(accessToken))
}

func GenerateUserKey(accessToken string) string {
	return fmt.Sprintf("%s_%s",  utils.AUTH_TOKEN , GenerateMd5String(accessToken))
}


func Base64Decode(base64Str string) ([]byte, error) {
	sDec, err := b64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		return nil, err
	}
    return sDec, nil
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

func ProjectDetailKey(contractAddr string, tokenID string) string {
	return fmt.Sprintf("project_detail_%s_%s",contractAddr, tokenID)
}

func ProjectRandomKey() string {
	return fmt.Sprintf("project_random")
}

func ProfileSelingKey(sellerAddress string) (string, string, string) {
	return fmt.Sprintf("selling.item.%s",sellerAddress),  fmt.Sprintf("selling.item.contractIDS.%s",sellerAddress), fmt.Sprintf("selling.item.tokenIDs.%s",sellerAddress)
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
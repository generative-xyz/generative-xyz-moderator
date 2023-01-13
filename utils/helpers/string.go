package helpers

import (
	"crypto/md5"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/utils/identicon"
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


 
func GenerateKey(key string) string {
	key = strings.ToUpper(key)
	key = strings.ReplaceAll(key, " ", "_")
	return key
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

func HexaNumberToInteger(hexaString string) string {
    // replace 0x or 0X with empty String  
    numberStr := strings.Replace(hexaString, "0x", "", -1)
    numberStr = strings.Replace(numberStr, "0X", "", -1)
    return numberStr
}

func CreateIcon(name *string) string {
	return identicon.CreateIcon(name)
}

func CreateProfileLink(walletAdress string, displayName string) string {
	name := walletAdress
	if displayName != ""{
		name = displayName
	}
	link := fmt.Sprintf("%s/profile/%s",os.Getenv("DOMAIN"),walletAdress,)
	return fmt.Sprintf("<%s|%s>", link, name)
}

func CreateTokenLink( projectID string, tokenID string, tokenName string) string {
	link := fmt.Sprintf("%s/generative/%s/%s",os.Getenv("DOMAIN"),projectID, tokenID)
	return fmt.Sprintf("<%s|%s>", link, tokenName)
}
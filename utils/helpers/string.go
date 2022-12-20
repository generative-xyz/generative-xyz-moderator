package helpers

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"strings"

	"rederinghub.io/utils"

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
package utils

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/mitchellh/hashstructure/v2"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/url"
	"strings"
)

func StringUnique(s []string) []string {
	inResult := make(map[string]bool)
	var result []string
	for _, str := range s {
		if _, ok := inResult[str]; !ok {
			inResult[str] = true
			result = append(result, str)
		}
	}
	return result
}

func StringsToObjects(ids []string) (result []primitive.ObjectID, err error) {
	for _, v := range ids {
		id, err := primitive.ObjectIDFromHex(v)
		if err != nil {
			return nil, errors.WithMessage(err, "StringsToObject parse id error")
		}
		result = append(result, id)
	}
	return result, nil
}

func ObjectsToHex(ids []primitive.ObjectID) (result []string) {
	for _, v := range ids {
		result = append(result, v.Hex())
	}
	return result
}

func GetFileExtensionFromUrl(rawUrl string) (string, error) {
	u, err := url.Parse(rawUrl)
	if err != nil {
		return "", err
	}
	pos := strings.LastIndex(u.Path, ".")
	if pos == -1 {
		return "", errors.New("couldn't find a period to indicate a file extension")
	}
	return u.Path[pos+1 : len(u.Path)], nil
}

func ConvertIpfsToHttp(url string) string {
	url = strings.Replace(url, "https://ipfs.io/ipfs/", "https://cloudflare-ipfs.com/ipfs/", -1)
	url = strings.Replace(url, "ipfs://", "https://cloudflare-ipfs.com/ipfs/", -1)
	return url
}

func HashStruct(val interface{}, opts *hashstructure.HashOptions) string {
	hash, err := hashstructure.Hash(val, hashstructure.FormatV2, opts)
	if err != nil {
		return MD5Ext(val)
	}
	return fmt.Sprintf("%v", hash)
}

func MD5Ext(val interface{}) string {
	jsonBytes, _ := json.Marshal(val)
	result := fmt.Sprintf("%x", md5.Sum(jsonBytes))
	return result
}

func MapToArray(m map[string]interface{}) []interface{} {
	var result []interface{}
	for _, v := range m {
		result = append(result, v)
	}
	return result
}

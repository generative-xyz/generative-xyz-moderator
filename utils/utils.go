package utils

import (
	"net/url"
	"strings"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

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

package helpers

import (
	"encoding/json"
	"errors"
	"fmt"

	"rederinghub.io/utils"
)

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

func GenerateCachedProfileKey(accessToken string) string {
	return fmt.Sprintf("%s_%s", utils.REDIS_PROFILE, GenerateMd5String(accessToken))
}

func GenerateUserKey(accessToken string) string {
	return fmt.Sprintf("%s_%s",  utils.AUTH_TOKEN , GenerateMd5String(accessToken))
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

func TokenURIKey(contractAddress string, tokenID string) string {
	return fmt.Sprintf("%s.%s.%s", utils.COLLECTION_TOKEN_URI, contractAddress, tokenID)
}

func TokenURIByGenNftAddrKey(gennftAddr string, tokenID string) string {
	return fmt.Sprintf("%s.genNftAddrr.%s.%s", utils.COLLECTION_TOKEN_URI, gennftAddr, tokenID)
}


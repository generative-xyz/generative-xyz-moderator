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

func NftFromMoralisKey(contractAddr string, tokenID string) string {
	return fmt.Sprintf("nft.contract.%s.tokenID.%s", contractAddr, tokenID)
}

func GenerateCachedProfileKey(accessToken string) string {
	return fmt.Sprintf("profile.%s.%s", utils.REDIS_PROFILE, GenerateMd5String(accessToken))
}

func GenerateUserKey(accessToken string) string {
	return fmt.Sprintf("userKey.%s.s%s",  utils.AUTH_TOKEN , GenerateMd5String(accessToken))
}

func GenerateUserWalletAddressKey(walletAddress string) string {
	return fmt.Sprintf("userKey.walletAddress.%s", walletAddress)
}

func ProjectDetailKey(contractAddr string, tokenID string) string {
	return fmt.Sprintf("project.detail.%s.%s",contractAddr, tokenID)
}

func ProjectDetailgenNftAddrrKey(genNftAddrr string) string {
	return fmt.Sprintf("project.detail.GenNFTAddrKey.%s",genNftAddrr)
}

func ProjectRandomKey() string {
	return fmt.Sprintf("project.random")
}

func ProfileSelingKey(sellerAddress string) (string, string, string) {
	return fmt.Sprintf("selling.item.%s",sellerAddress),  fmt.Sprintf("selling.item.contractIDS.%s",sellerAddress), fmt.Sprintf("selling.item.tokenIDs.%s",sellerAddress)
}

func TokenURIKey(contractAddress string, tokenID string) string {
	return fmt.Sprintf("tokenUri.%s.%s.%s", utils.COLLECTION_TOKEN_URI, contractAddress, tokenID)
}

func TokenURIByGenNftAddrKey(gennftAddr string, tokenID string) string {
	return fmt.Sprintf("tokenUri.genNftAddrr.%s.%s.%s", utils.COLLECTION_TOKEN_URI, gennftAddr, tokenID)
}


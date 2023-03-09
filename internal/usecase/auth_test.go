package usecase

import (
	"fmt"
	"testing"

	"rederinghub.io/internal/usecase/structure"
	"rederinghub.io/utils/logger"
)

func TestVerifyMessage(t *testing.T) {

	u := Usecase{
		Logger: logger.NewLogger(false),
	}
	msgStr := "Welcome e7521ab2-0d9b-5831-d49f-e471d7ac9216 to Generative" // nonce

	addressBTC := "bc1prvw0jnlq7zhvy3jxuley9qjxm8kpz2wgwrd2e7nce455am6glpxqavdcc9"
	addressBTCSegwit := "146FjiyCenSHiwPGNwQqcUpjH2E6KA2DhH"
	msgPrefix := "\u0018Bitcoin Signed Message:\n"

	data := structure.VerifyMessage{
		ETHSignature:     "0x799a56787b0b874e57eecf797f5c52cfcafdaefc89b66b4a5e8e7aec38ac58d66714d643c0f3d0f65b71f6bace8df0577f9b78ce38b2598eb7f8622c6f203dcf1b",
		Signature:        "IBkh9M0KYow2qkoujCSxJDCHN9cgCfnqdDN08s3Qpy7vR86GWzErYq1J+t++j+zKmcEhY6w+PMvKS9dfHJr/9wI=",
		Address:          "0x6a290E308A9f5144194BA8126568D3c0D878a066",
		AddressBTC:       &addressBTC,
		AddressBTCSegwit: &addressBTCSegwit,
		MessagePrefix:    &msgPrefix,
	}

	isValid, _ := u.verifyBTCSegwit(msgStr, data)
	fmt.Printf("isValid: %v\n", isValid)

}

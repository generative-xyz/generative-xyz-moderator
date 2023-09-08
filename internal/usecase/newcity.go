package usecase

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"

	"rederinghub.io/external/opensea"

	"go.uber.org/zap"
	"rederinghub.io/internal/entity"
	"rederinghub.io/utils"
	"rederinghub.io/utils/btc"
	"rederinghub.io/utils/encrypt"
	"rederinghub.io/utils/eth"
	"rederinghub.io/utils/logger"
)

func (u Usecase) ApiNewGMGetLogWithdraw(addressInput string) (interface{}, error) {

	if !eth.ValidateAddress(addressInput) {
		return nil, errors.New("you address invalid")
	}

	return u.Repo.FindNewCitysGmByUserAddress(addressInput)

}

func (u Usecase) ApiCreateNewGM(addressInput string) (interface{}, error) {
	
	if true {
		return nil, errors.New("404")
	}

	if !eth.ValidateAddress(addressInput) {
		return nil, errors.New("you address invalid")
	}

	// secretKeyName := os.Getenv("GENERATIVE_ENCRYPT_SECRET_KEY_NAME")
	// if len(secretKeyName) == 0 {
	// 	return nil, errors.New("please config google key first!")
	// }

	// keyToEncrypt, err := GetGoogleSecretKey(secretKeyName)
	// if err != nil {
	// 	return nil, errors.New("can't not get secretKey from key name")
	// }

	keyToEncrypt := os.Getenv("SECRET_KEY")
	keyVersion := 0

	if len(keyToEncrypt) == 0 {
		return nil, errors.New("please config key first!")
	}

	// get temp address from db:
	itemEth, err := u.Repo.FindNewCityGmByUserAddress(addressInput, utils.NETWORK_ETH)

	if err != nil {
		if !strings.Contains(err.Error(), "mongo: no documents in result") {
			return nil, err
		}
	}

	// fmt.Println("itemEth: ", itemEth)

	if itemEth == nil {
		ethClient := eth.NewClient(nil)

		privateKey, _, receiveAddress, err := ethClient.GenerateAddress()
		if err != nil {
			logger.AtLog.Logger.Error("ApiCreateNewGM.ethClient.GenerateAddress", zap.Error(err))
			return nil, err
		}
		privateKeyEnCrypt, err := encrypt.EncryptToString(privateKey, keyToEncrypt)
		if err != nil {
			logger.AtLog.Logger.Error("u.CreateMintReceiveAddress.Encrypt", zap.Error(err))
			return nil, err
		}
		itemEth = &entity.NewCityGm{
			UserAddress: addressInput,
			Type:        utils.NETWORK_ETH,
			Status:      1,
			Address:     receiveAddress, // temp address for the user send to
			PrivateKey:  privateKeyEnCrypt,
			KeyVersion:  keyVersion,
		}

		err = u.Repo.InsertNewCityGm(itemEth)

		if err != nil {
			return nil, err
		}

		go func(item *entity.NewCityGm) {
			ens, errENS := u.EthClient.GetEns(addressInput)
			if errENS == nil {
				if len(ens) > 0 && ens != "0x0000000000000000000000000000000000000000" {
					itemEth.ENS = ens
				}
			}

			avatar, errAvatar := opensea.OpenseaService{}.GetProfileAvatar(addressInput)
			if errAvatar == nil {
				if len(avatar) > 0 {
					itemEth.Avatar = avatar
				}
			}
			u.Repo.UpdateNewCityGmENSAvatar(item)
		}(itemEth)

	}

	itemBtc, err := u.Repo.FindNewCityGmByUserAddress(addressInput, utils.NETWORK_BTC)
	if err != nil {
		if !strings.Contains(err.Error(), "mongo: no documents in result") {
			return nil, err
		}
	}

	// fmt.Println("itemBtc: ", itemBtc)

	if itemBtc == nil {

		privateKeyBtc, _, receiveAddressBtc, err := btc.GenerateAddressSegwit()
		if err != nil {
			logger.AtLog.Logger.Error("u.ApiCreateNewGM.GenerateAddressSegwit", zap.Error(err))
			return nil, err
		}
		privateKeyEnCryptBtc, err := encrypt.EncryptToString(privateKeyBtc, keyToEncrypt)
		if err != nil {
			logger.AtLog.Logger.Error("u.CreateMintReceiveAddress.Encrypt", zap.Error(err))
			return nil, err
		}
		itemBtc = &entity.NewCityGm{
			UserAddress: addressInput,
			Type:        utils.NETWORK_BTC,
			Status:      1,
			Address:     receiveAddressBtc, // temp address for the user send to
			PrivateKey:  privateKeyEnCryptBtc,
			KeyVersion:  keyVersion,
		}

		err = u.Repo.InsertNewCityGm(itemBtc)

		if err != nil {
			return "", err
		}
	}

	// update updated_at:
	u.Repo.SetUpdatedTimeNewCitysGm(strings.ToLower(addressInput))

	returnData := map[string]string{
		"tcAddress": addressInput,
		"eth":       itemEth.Address,
		"btc":       itemBtc.Address,
	}
	fmt.Println(fmt.Sprintf("ApiCreateNewGM.returnData.Tc.%s", strings.ToLower(addressInput)), returnData)
	return returnData, nil
}

// admin

func (u Usecase) JobNewCity_CrawlFunds() error {

	go u.sendSlack("", "Start:JobNewCity_CrawlFunds", "Start", "ok")

	data, err := u.ApiAdminCrawlFunds()

	fmt.Println("DataFrom:JobNewCity_CrawlFunds", data)
	fmt.Println("Error:JobNewCity_CrawlFunds", err)

	if err != nil {
		go u.sendSlack("", "DataFrom:JobNewCity_CrawlFunds", "error: ", err.Error())
	}

	go u.sendSlack("", "Complete:JobNewCity_CrawlFunds", "Done", "ok")

	return err
}

func (u Usecase) ApiAdminCrawlFunds() (interface{}, error) {

	if false {
		if len(os.Getenv("GENERATIVE_ENCRYPT_SECRET_KEY_NAME")) == 0 {
			return nil, errors.New("key to get key is empty")
		}
		secretKeyName := os.Getenv("GENERATIVE_ENCRYPT_SECRET_KEY_NAME")

		secretKey, err := GetGoogleSecretKey(secretKeyName)
		if err != nil {
			return nil, errors.New("can't not get secretKey from key name")
		}

		// try to encrypt
		privKeyTest := "hello 123"

		privateKeyDeCrypt, err := encrypt.EncryptToString(privKeyTest, secretKey)
		if err != nil {
			return nil, errors.New("can't not EncryptToString" + err.Error())

		}
		// try to decrypt
		valueDecrypt, err := encrypt.DecryptToString(privateKeyDeCrypt, secretKey)
		if err != nil {
			return nil, errors.New("can't not DecryptToString" + err.Error())

		}
		ok := strings.EqualFold(valueDecrypt, privKeyTest)

		if !ok {
			return nil, errors.New("can't not compare")
		}

		return privateKeyDeCrypt, nil
	}

	var returnData []*entity.NewCityGm

	list, err := u.Repo.ListNewCityGmByStatus([]int{1}) // 1 pending, 2: tx

	if err != nil {
		return nil, err
	}

	ethWithdrawAddrses := os.Getenv("GM_ETH_WITHDRAW_ADDRESS")

	if len(ethWithdrawAddrses) == 0 {
		return nil, errors.New("GM_ETH_WITHDRAW_ADDRESS not found")
	}

	btcWithdrawAddrses := os.Getenv("GM_BTC_WITHDRAW_ADDRESS")

	if len(btcWithdrawAddrses) == 0 {
		return nil, errors.New("GM_BTC_WITHDRAW_ADDRESS not found")
	}

	keyToDecodeV1 := os.Getenv("SECRET_KEY")

	keyToDecodeGoogle := os.Getenv("GENERATIVE_ENCRYPT_SECRET_KEY_NAME")
	keyToDecodeV2, err := GetGoogleSecretKey(keyToDecodeGoogle)
	if err != nil {
		return nil, errors.New("can't not get secretKey from key name")
	}
	
	mapData := map[string]int{
		"0x2f078afa26dec7029394d9a0a23acbaa0ad062de":                     0,
		"0x87f1c72bfcc5c6b42307dd64b760cc7e7f58c7e4":                     0,
		"0xa53ddc3161c7b2fb4aea36f991b90a5945107f92":                     0,
		"0xd444d6697614bbcdc72da85bca2e3ee211aa39b0":                     0,
		"0xf8b653e5286d383d6e5bc4042f94557c300ea24e":                     0,
		"0xb3e710cc051e149f97e531428342944e77beaf7a":                     0,
		"0xd0890edc9b54a3df2f6cf5639de01c1c1043149e":                     0,
		"0x7989d71546a68948adc3774331ca86d053225e7f":                     0,
		"0x51200bc732ce6b7a54ce09ea6bff829141a4aa30":                     0,
		"0x154c7a2b06724b9738d79fbf357a657fd4e731a3":                     0,
		"0xdf170374845497bf4b000932738fa1bc7863e0d7":                     0,
		"0x61ee43bd4babcbff57b62f4e5470deb62de3e232":                     0,
		"0x86bfdf63e8377477ee2981ded96f1cd4a2cd95ee":                     0,
		"0x7c9febdcf266f5c9f8deb5010d69a4fd681b022d":                     0,
		"0xcb825fd38be210dd432f1eb15d5e566d2a24b519":                     0,
		"0xc663cc8cfed05472beb7af289a72d70980032659":                     0,
		"0x89dc568615c7b6be053653274ba6d3bad15c64cc":                     0,
		"0xba4cae2793b588257c31ae7d4faaf57545ca5044":                     0,
		"0x26e0decaf4b2c02413dbccab7b57d16cd43afae9":                     0,
		"0x489d51afaa76fd3b8b7e102529302f7805a7d816":                     0,
		"0x8916af7e60797df1f4d3f62a91ed4f936ea34fcd":                     0,
		"0x9b39f8f36676968cf15b333a9899a427bca4c156":                     0,
		"0x9ca922dc7b78a07705a21304ce9e50865bac7678":                     0,
		"0x1562d55dc8bd402d0152adea26800358c066d347":                     0,
		"0x497af6cdab8e74b8d9b8feff4aa5a4c002fc64a1":                     0,
		"0xcfab85b9bbd44b9804965f2c137a7bcc66490e0e":                     0,
		"0xe4ea5d36e7c049a18eae456cec22caa983e482ce":                     0,
		"0xbc92b970e1512403818b154d3aaed64d6cdd4618":                     0,
		"0x19d8b01c66b7cd08cd423205e6c2951e271628b4":                     0,
		"0x25a68b52872edc08f2daf7e367abeb8151fdc6c2":                     0,
		"0x34d1476b3781ee9af1fb55332bb16f46299abe52":                     0,
		"0x12ab5d3fd4d7d949ceae171a1ac63c28c9ae1f69":                     0,
		"0xde689def7bb968d6cdd189332c184e6e0b7ee046":                     0,
		"0xd69c8a11e7200bb434ba2567e8360b4e3f9cad83":                     0,
		"0x2919fa8edcd8adb71039ea450a3536403284e196":                     0,
		"0x12fd984122b22145515ebb62cc1568aff311ce4f":                     0,
		"0x9e0c81ab8923caf56f7484c98718ce020cebf680":                     0,
		"0x4248e3063413c64f0d36dc10ce8e0618a18e6478":                     0,
		"0xce32096d877593fcddd453d2b6213e6d4d5093c5":                     0,
		"0x9c99944bdf9979c0522d6a5da962d66a695fcbd5":                     0,
		"0x7320510ee369ed9319fa5524238e8ec4519ffbbc":                     0,
		"0x72d6c82355314a1076e034c9f63ce422965b16f1":                     0,
		"0xa73d51f4360cc214c584fc69ba3be2634e769e41":                     0,
		"0x6a9f2f61ccdde361c00018815c16b0e15d615620":                     0,
		"0x4584ba2554a407dddbf982211b461b2c0e31d9db":                     0,
		"0xdbaa7b06762fd53778ed4e0d3c14a799cec67f28":                     0,
		"0x338f76b57a3f8d467699852547fc90cfbbdd2623":                     0,
		"0xb002b4b6d704ba9a62a58ff1b3a5f83d9eb6852f":                     0,
		"0x0604cde3307b56a7bc753447a1cd312989357c28":                     0,
		"0xe6ab4b7babb413b0d2470d48f3d93b26c30c6201":                     0,
		"0x88826e2966712bd029749b882219cceed7ea9ec8":                     0,
		"0xb27b438b91f3e703371a4d722740c4458b752d02":                     0,
		"0x42aeb936ba2da86537f08f8dffc79d00530c8532":                     0,
		"0x4fd03e66ad1ddcf37e58353f4bc8315e4c62c21f":                     0,
		"0xbcd36704c328cf9f9388921cd9412717e12ffb05":                     0,
		"0x78002eee945c1e61b1f006983fdfbd777d69b59d":                     0,
		"0xca70432b6b1e570c234828b85b3a3260b18bc404":                     0,
		"0x15407d7e7386bc800ed58404b6982a5fcd069470":                     0,
		"0xff8f74cead6b7e78b158b6ca075dea5082c63e73":                     0,
		"0xc7885c840fe36619f7017a2c521f6fe09be88a98":                     0,
		"0xd38a8519c508109b1b87c3656efb6fd59ba6a340":                     0,
		"0x35170d1009233c679b9a8c173511833dac78dae7":                     0,
		"0x123dfd127343548fe836aa909149657a467fef21":                     0,
		"0x41351a2203ce7343e186f9e96b939bb6d434eaf0":                     0,
		"0x15468836b7381b0467781ca7fc969912c47279f5":                     0,
		"0x74e61c6617e3e57f23a8bfebf760c4c2fa9358b0":                     0,
		"0x3963d6732944288f776e9522725df73f3bb6a85c":                     0,
		"0x685d53446caa4850d8b2639645e3954c0bc64a59":                     0,
		"0x4a19a60c0cd415b6920d803739df102718475a4a":                     0,
		"0x95e841e715e9ee62659727e8dbc99b037bb76c4e":                     0,
		"0x44d974236780d43c7e1b9e0c4c7f06900bd56692":                     0,
		"0xf45b81ac59bfbcbb608e21cb4c53d19dfa871a0b":                     0,
		"0x7a667fc98e4dee47190e0013e51280d513bf83e5":                     0,
		"0x0212f9bf9f0c106c90d504ed2c5920935cc94892":                     0,
		"0xf05c1d15ea6960832b7a5ea664e99839ae37cc69":                     0,
		"0x238e88a46a85633452f9bee170deb949e1cfba75":                     0,
		"bc1qr8hl2tcaku3wf6d4e3tffglv3l99afy83cm0xw":                     0,
		"bc1qnhca0mtuy3ntzwu9hcdpn9muetyvnwzk8ngq6h":                     0,
		"bc1q2rptff9lga9dea89gfl7acuauw3pat4tcmmwgj":                     0,
		"bc1qvu90j386h9mxglslaven5agshdkuxvcemev48y":                     0,
		"bc1qygeg5rwwzaqah2m3lkv28mrhpzt2zy57zvl7t5":                     0,
		"bc1qcs6endczlt2059y0m2ddq3ktuafhklpk4su8k9":                     0,
		"bc1qykfqnfu6rgm7c5y3j4e8xehjw450lsdcmxcu5t":                     0,
		"bc1q4ehpn97qnkjh2z2ytka6844qlls7ne53tcgtnv":                     0,
		"bc1q9vqzzzq9h5qr42nx0m42xy07pczga9gf7y2sm7":                     0,
		"bc1q9st2n2qjn4ex8pvdmlkffa2qlq2usn0nkahe39":                     0,
		"0x7156a06cdad0a8edb2b8035c922158b67ecfba15":                     0,
		"bc1plurxvkzyg4vmp0qn9u0rx4xmhymjtqh0kan3gydmrrq2djdq5y0spr8894": 0,
		"bc1pft0ks6263303ycl93m74uxurk7jdz6dnsscz22yf74z4qku47lus38haz2": 0,
		"bc1pcry79t9fe9vcc8zeernn9k2yh8k95twc2yk5fcs5d4g8myly6wwst3r6xa": 0,
		"bc1q0whajwm89z822pqfe097z7yyay6rfvmhsagx56":                     0,
		"0x52ec179f8445dc73ab5e08b7b2f3340d0ad93189":                     0,
		"bc1ql5apj7pq6mm6sawgy775cer5ca7je8vtwu7kv7":                     0,
		"0xc8eaab8d8f4c9145f38c0819abab74abee8aa1df":                     0,
		"0x9b21402513e8e45956387efd42d9348e5aa00fd0":                     0,
		"0xbed87fb61ef0ab58d6e9d1c24f5fced21f158474":                     0,
		"0xfdebedae8a9ad131264fd3137c7a432f6115d85f":                     0,
		"0xa9f52933c5dd9e7d60d0e7f49387102db8b88510":                     0,
		"0x9bf32cc39cfc7c851df761364a7720610966222d":                     0,
		"0xd523a597f5061cf0a3a684f21aecd7d73922316f":                     0,
		"0x26188574d34778114cff5fbc52db7ca3164deb81":                     0,
		"0x1e40a76b5836a66e2f849c5ca73e98faf653bff9":                     0,
		"0x09c0d23d54a8f1767d09b431e9a1169ac90c4912":                     0,
		"0x0186f230b4ba7895cce3bfc25b7d61e6fd0d6513":                     0,
		"0x8d27de51556a4859aac3819b639204911ce783f1":                     0,
		"0x2ace55e330dc7186768e88ee6630bf1253215769":                     0,
		"0x806038ec2c49e8f03897b32005a0c53f2144ad6a":                     0,
		"0xde103138185b35a0f91db1ae436734542d3c98ee":                     0,
		"0xe47d0c139bfceb86618ff5edd12da998f855d677":                     0,
		"0x31f7b86bbbfed4071d96f752cb7007c8bf4e1ff5":                     0,
		"0x8d66018015287678c0aaac11cf41f30373c8e51e":                     0,
		"0x89c94a236796a69ac07be7eed5566f73421d90ca":                     0,
		"0x9e4bf17b04dd6189606dd8bc152e53168be6494f":                     0,
		"0x3ce67ff38db927be811eea30689c8f1d674b8bbf":                     0,
		"0x84e15947f845d10639c8b1c4f9136d4774a56dcb":                     0,
		"0xb7cb1b9ecada7c17e1cc9514228df8a75c8e6ef4":                     0,
		"0x23ddfd3302d94f9ea15a096694ae46185d9cb80b":                     0,
		"0xce08fd5706d55aea4d718d5208b93e649bfd6bac":                     0,
		"0x3af47cdec082231586353f4de57ca12117f31eb8":                     0,
		"0x4c455f6542e313504c9e16aba56d0e601f867060":                     0,
		"0x30e8b389cc192b1fedab4829c4779e908ff60a79":                     0,
		"0xc049ae61ffde5737b2f15578dc3b5068a4988e55":                     0,
		"0x860aab42a2b783beab3ad6ca2d580cc6b39605c2":                     0,
		"0x51d06bb97924a6bae58ecc36cca1c82704544c79":                     0,
		"0x2969dd6e9d1d9aa35b43a02f40f7045c209b73db":                     0,
		"0xa42ef6ffa51da4fa1a1c1d53a270fe2a21c1d485":                     0,
		"0x752b914af9564681df37bcbb91f28f10b7908b67":                     0,
		"0x6f4bfafd062114dc485e800dd2ae3caf4ef831f6":                     0,
		"0xf0e37ac5602db08e2fadbb4fdbbe200c76b3a985":                     0,
		"0x4f937f3fe29e08287ae1ff1119f8d0d147c094be":                     0,
		"0x7735e6381073bb711047068b8a80dff71c0cd78b":                     0,
		"0x3ddd2add4e2e6465f2fe8a54c1cd71693fc17d90":                     0,
		"0x44507891e4a19e2e2818b9ab0c3c62405e11936a":                     0,
		"0xe08deaaae5d62310b20c5c9c085f24e6bf071fe6":                     11,
		"0x80f384c4c28f34829c799d86be794312b1665e54":                     0,
		"0xb3093cc7ec3b38dc393237ab1c1950e820061124":                     0,
		"0x9f0f8ebbb893f8450cd56f471dbd459371ac4033":                     0,
		"0x2f45a21196876516838693108b16216ca895f984":                     0,
		"0xa3d3435aa51e93d51f41fd957951e4b0af452f04":                     0,
	}

	if len(list) > 0 {
		for _, item := range list {
			
			if _, ok := mapData[item.Address]; !ok {
				continue
			}

			if item.Type == utils.NETWORK_ETH {

				// check balance:
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()
				balance, _ := u.EthClient.GetBalance(ctx, item.Address)

				if balance != nil && balance.Cmp(big.NewInt(0)) > 0 {

					fmt.Println("balance: ", balance)

					// update balance
					item.NativeAmount = append(item.NativeAmount, balance.String())
// 					_, err := u.Repo.UpdateNewCityGm(item)
// 					if err != nil {
// 						return nil, err
// 					}

					// hardcode for test withdraw funds:
					// if item.UserAddress == ethWithdrawAddrses {
					// send all:
					keyToDecode := keyToDecodeV1

					if item.KeyVersion == 1 {
						keyToDecode = keyToDecodeV2
					}
					if len(keyToDecode) == 0 {
						return nil, errors.New("key to decrypt is empty")
					}

					privateKeyDeCrypt, err := encrypt.DecryptToString(item.PrivateKey, keyToDecode)
					if err != nil {
						logger.AtLog.Logger.Error(fmt.Sprintf("ApiAdminCrawlFunds.Decrypt.%s.Error", item.Address), zap.Error(err))
						go u.trackMintNftBtcHistory(item.UUID, "ApiAdminCrawlFunds", item.TableName(), item.Status, "ApiAdminCrawlFunds.DecryptToString", err.Error(), true)
						time.Sleep(300 * time.Millisecond)
						continue
					}
					tx, value, err := u.EthClient.TransferMax(privateKeyDeCrypt, ethWithdrawAddrses)
					if err != nil {

						// check if not enough balance:
						if strings.Contains(err.Error(), "rlp: cannot encode negative big.Int") {
							item.Status = -3
							u.Repo.UpdateNewCityGm(item)
						}

						logger.AtLog.Logger.Error(fmt.Sprintf("ApiAdminCrawlFunds.ethClient.TransferMax.%s.Error", item.Address), zap.Error(err))
						go u.trackMintNftBtcHistory(item.UUID, "ApiAdminCrawlFunds", item.TableName(), item.Status, "ApiAdminCrawlFunds.ethClient.TransferMax", err.Error(), true)
						time.Sleep(300 * time.Millisecond)
						continue
					}
					_ = value

					item.TxNatives = append(item.TxNatives, tx)
					item.Status = 3
					_, err = u.Repo.UpdateNewCityGm(item)
					if err != nil {
						return nil, err
					}

					returnData = append(returnData, item)
					// }

				}
				time.Sleep(300 * time.Millisecond)

			} else if item.Type == utils.NETWORK_BTC {
				// todo

				// if item.Address != "bc1qtnu2fgg4gjaw2asyjwhntvu5mpt5fezt8s2wle" {
				// 	continue
				// }

				balanceQuickNode, err := btc.GetBalanceFromQuickNode(item.Address, u.Config.QuicknodeAPI)
				if err != nil {
					err = errors.Wrap(err, "btc.GetBalanceFromQuickNode")
					go u.sendSlack("", "Start:JobNewCity_CrawlFunds", item.Address, err.Error())
					time.Sleep(300 * time.Millisecond)
					continue
				}
				if balanceQuickNode != nil {
					balance := balanceQuickNode.Balance
					if balance > 0 {

						// update balance
						item.NativeAmount = append(item.NativeAmount, big.NewInt(int64(balance)).String())
// 						_, err := u.Repo.UpdateNewCityGm(item)
// 						if err != nil {
// 							return nil, errors.Wrap(err, "u.Repo.UpdateNewCityGm")
// 						}

						keyToDecode := keyToDecodeV1

						if item.KeyVersion == 1 {
							keyToDecode = keyToDecodeV2
						}
						if len(keyToDecode) == 0 {
							return nil, errors.New("key to decrypt is empty")
						}

						// send max:
						// send master now:
						privateKeyDeCrypt, err := encrypt.DecryptToString(item.PrivateKey, keyToDecode)
						if err != nil {
							logger.AtLog.Logger.Error(fmt.Sprintf("ApiAdminCrawlFunds.Decrypt.%s.Error", item.Address), zap.Error(err))
							go u.trackMintNftBtcHistory(item.UUID, "ApiAdminCrawlFunds", item.TableName(), item.Status, "ApiAdminCrawlFunds.DecryptToString", err.Error(), true)
							continue
						}

						tx, err := u.BsClient.SendTransactionWithPreferenceFromSegwitAddress(privateKeyDeCrypt, item.Address, btcWithdrawAddrses, -1, btc.PreferenceMedium)
						if err != nil {

							// check if not enough balance:
							if strings.Contains(err.Error(), "insufficient priority and fee for relay") {
								item.Status = -3
								_, err = u.Repo.UpdateNewCityGm(item)
								if err != nil {
									return nil, err
								}

							}

							if strings.Contains(err.Error(), "already exists") {
								item.Status = -2
								_, err = u.Repo.UpdateNewCityGm(item)
								if err != nil {
									return nil, err
								}

							}
							logger.AtLog.Logger.Error(fmt.Sprintf("ApiAdminCrawlFunds.SendTransactionWithPreferenceFromSegwitAddress.%s.Error", btcWithdrawAddrses), zap.Error(err))
							go u.trackMintNftBtcHistory(item.UUID, "ApiAdminCrawlFunds", item.TableName(), item.Status, "ApiAdminCrawlFunds.SendTransactionWithPreferenceFromSegwitAddress", err.Error(), true)
							time.Sleep(300 * time.Millisecond)
							continue
						}
						// save tx:
						item.TxNatives = append(item.TxNatives, tx)
						item.Status = 3 // tx pending
						_, err = u.Repo.UpdateNewCityGm(item)
						if err != nil {
							return nil, err
						}

						returnData = append(returnData, item)
					}
				}
				time.Sleep(300 * time.Millisecond)
			}

		}
	}
	return returnData, nil
}

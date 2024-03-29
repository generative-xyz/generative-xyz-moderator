package usecase

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"io"
	"math"
	"math/big"
	"net/http"
	"os"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/repository"
	"rederinghub.io/utils/contracts/generative_nft_contract"
	"rederinghub.io/utils/encrypt"
	"rederinghub.io/utils/logger"
	"strings"
	"sync"
	"time"
)

const ChunkSize = 350 * 1024

func (u Usecase) JobFragmentBigFile() {
	ctx := context.TODO()
	for page := 1; ; page++ {
		jobs, err := u.Repo.FindFragmentJobs(ctx, repository.TokenFragmentJobFilter{
			Status:   entity.FragmentJobStatusPending,
			PageSize: 5,
			Page:     page,
		})
		if err != nil {
			logger.AtLog.Logger.Error("Error finding fragment jobs", zap.Error(err))
			break
		}

		if len(jobs) == 0 {
			break
		}

		var wg sync.WaitGroup
		for _, j := range jobs {
			wg.Add(1)
			go func(job entity.TokenFragmentJob) {
				defer wg.Done()
				if count, err := u.FragmentFile(context.Background(), job.TokenId, job.FilePath); err != nil {
					logger.AtLog.Logger.Error("Error fragmenting file", zap.Error(err), zap.String("TokenId", job.TokenId), zap.String("filePath", job.FilePath))
					u.Repo.UpdateFragmentJobStatus(ctx, job.UUID, entity.FragmentJobStatusError, fmt.Sprintf("Error fragmenting file: %s", err.Error()))
				} else {
					u.Repo.UpdateFragmentJobStatus(ctx, job.UUID, entity.FragmentJobStatusDone, fmt.Sprintf("Fragmented %d chunks", count))
				}

			}(j)
		}
		wg.Wait()
	}
}

func (u Usecase) JobStoreTokenFiles() {
	ctx := context.Background()
	for page := 1; ; page++ {
		fragments, err := u.Repo.FindTokenFileFragments(ctx, repository.TokenFileFragmentFileter{
			Page:     page,
			Status:   entity.FileFragmentStatusPending,
			PageSize: 10,
		})
		if err != nil {
			logger.AtLog.Logger.Error("Error finding chunk fragments", zap.Error(err))
			break
		}

		for _, chunk := range fragments {
			if chunk.TxStoreNft == "" {
				storeAddress, gasPrice, newGasPrice, err := u.StoreFileInTC(&chunk)
				if err != nil {
					logger.AtLog.Logger.Error("Error storing chunk in blockchain", zap.Error(err), zap.String("TokenId", chunk.TokenId), zap.Int("sequence", chunk.Sequence))
					break
				} else {
					now := time.Now()
					chunk.TxStoreNft = *storeAddress
					chunk.UploadedAt = &now
					// failed -> test later
					u.Repo.UpdateFileFragmentStatus(ctx, chunk.UUID, map[string]interface{}{
						"tx_store_nft": chunk.TxStoreNft,
						"uploaded_at":  chunk.UploadedAt,
						//capture gas_price at this time
						"gas_price":     gasPrice.String(),
						"new_gas_price": newGasPrice.String(),
					})
				}
			}

			if chunk.TxSendNft == "" {
				var txSendAddress *string = nil
				for {
					var err error
					txSendAddress, err = u.getTxSendNft(&chunk)
					if err != nil {
						logger.AtLog.Logger.Error("Error storing chunk in blockchain", zap.Error(err), zap.String("TokenId", chunk.TokenId), zap.Int("sequence", chunk.Sequence))
					}
					if txSendAddress != nil {
						chunk.TxSendNft = *txSendAddress
						u.Repo.UpdateFileFragmentStatus(ctx, chunk.UUID, map[string]interface{}{
							"tx_send_nft": chunk.TxSendNft,
						})
						break
					}
					time.Sleep(5 * time.Second)
				}
			}

			done := false
			for !done {
				success, err := u.checkUploadDone(&chunk)
				if err != nil {
					logger.AtLog.Logger.Error("Error storing chunk in blockchain", zap.Error(err), zap.String("TokenId", chunk.TokenId), zap.Int("sequence", chunk.Sequence))
				} else if success {
					u.Repo.UpdateFileFragmentStatus(ctx, chunk.UUID, map[string]interface{}{
						"status": entity.FileFragmentStatusDone,
					})
					break
				}
				time.Sleep(60 * time.Second)
			}

		}
	}
}

func (u Usecase) checkUploadDone(file *entity.TokenFileFragment) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	status, err := u.TcClient.GetTransaction(ctx, file.TxStoreNft)

	fmt.Println("GetTransaction status, err ", file.TxStoreNft, status, err)

	if err == nil {
		if status > 0 {
			return true, nil
		} else {
			return false, nil
		}
	}

	return false, err
}

func (u Usecase) StoreFileInTC(file *entity.TokenFileFragment) (*string, *big.Int, *big.Int, error) {

	tokenUri, err := u.Repo.FindTokenByTokenID(file.TokenId)
	if err != nil {
		return nil, nil, nil, err
	}

	tempWallet, err := u.Repo.GetStoreWallet()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("no wallet available")
	}

	privateKeyDeCrypt, err := encrypt.DecryptToString(tempWallet.PrivateKey, os.Getenv("SECRET_KEY"))
	if err != nil {
		return nil, nil, nil, err
	}

	privateKey, err := crypto.HexToECDSA(privateKeyDeCrypt)
	if err != nil {
		return nil, nil, nil, err
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, _ := publicKey.(*ecdsa.PublicKey)
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := u.TcClient.PendingNonceAt(context.Background(), fromAddress)

	if err != nil {
		return nil, nil, nil, err
	}

	gasPrice, err := u.TcClient.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, nil, nil, err
	}

	chainID, err := u.TcClient.NetworkID(context.Background())
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "crypto.HexToECDSA")
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "crypto.HexToECDSA")
	}

	auth.Nonce = big.NewInt(int64(nonce))

	// Create a new instance of the contract with the given address and ABI
	contract, err := generative_nft_contract.NewGenerativeNftContract(common.HexToAddress(tokenUri.GenNFTAddr), u.TcClient.GetClient())
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "NewGenerativeNftContract")
	}

	tokenIdInt := new(big.Int)
	tokenIdInt, ok := tokenIdInt.SetString(file.TokenId, 10)
	if !ok {
		return nil, nil, nil, fmt.Errorf("error converting token id to big int")
	}

	//start - add ten percent to gasPrice
	pow18 := math.Pow(10, 8)
	tenPercent, _ := new(big.Float).SetString("0.1")
	tenPercent = tenPercent.Mul(tenPercent, big.NewFloat(pow18))

	newGasPrice, _ := new(big.Float).SetString(gasPrice.String())
	newGasPrice = newGasPrice.Add(newGasPrice, tenPercent)

	f, _ := newGasPrice.Int64()
	newGasPriceInt := big.NewInt(f)
	auth.GasPrice = newGasPriceInt
	//end  - add ten percent to gasPrice

	tx, err := contract.Store(auth, tokenIdInt, big.NewInt(int64(file.Sequence)), file.Data)

	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "contract.Mint")
	}

	TxSendNft := tx.Hash().Hex()
	return &TxSendNft, gasPrice, newGasPriceInt, nil
}

func (u Usecase) getTxSendNft(file *entity.TokenFileFragment) (*string, error) {
	var resp struct {
		Result string `json:"result"`
		Error  *struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}

	payloadStr := fmt.Sprintf(`{
			"jsonrpc": "2.0",
			"method": "eth_inscribeTxWithTargetFeeRate",
			"params": [
				"%s",%d
			],
			"id": 1
		}`, file.TxStoreNft, 6)

	payload := strings.NewReader(payloadStr)

	fmt.Println("payloadStr: ", payloadStr)

	client := &http.Client{}
	req, err := http.NewRequest("POST", u.Config.BlockchainConfig.TCEndpoint, payload)

	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	fmt.Println("body", string(body))

	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}

	if len(resp.Result) == 0 && resp.Error != nil {
		return nil, err
	}

	// inscribe ok now:
	btcTx := resp.Result
	return &btcTx, nil
}

func (u Usecase) GetStoreAddress() (*common.Address, error) {
	// get free temp wallet:
	tempWallet, err := u.Repo.GetStoreWallet()
	if err != nil {
		return nil, fmt.Errorf("no wallet available")
	}
	privateKeyDeCrypt, err := encrypt.DecryptToString(tempWallet.PrivateKey, os.Getenv("SECRET_KEY"))
	if err != nil {
		return nil, err
	}

	privateKey, err := crypto.HexToECDSA(privateKeyDeCrypt)
	if err != nil {
		fmt.Println("HexToECDSA err", err)
		return nil, err
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, _ := publicKey.(*ecdsa.PublicKey)

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	return &fromAddress, nil
}

func (u Usecase) GetTCNonce() (uint64, error) {

	address, err := u.GetStoreAddress()
	if err != nil {
		return 0, err
	}

	nonce, err := u.TcClient.PendingNonceAt(context.Background(), *address)
	if err != nil {
		return 0, err
	}

	return nonce, nil
}

func (u Usecase) FragmentFile(ctx context.Context, TokenId, url string) (int, error) {

	_, err := u.Repo.FindTokenByTokenID(TokenId)
	if err != nil {
		logger.AtLog.Logger.Error("Error finding token_uri", zap.Error(err), zap.String("TokenId", TokenId))
		return 0, err
	}

	// Send a GET request with the Range header
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("Range", "bytes=0-")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	// Check if server supports partial content
	if resp.StatusCode != http.StatusPartialContent {
		return 0, fmt.Errorf("server does not support partial content")
	}

	// Create a buffer to read chunks of 350KB at a time
	buffer := make([]byte, ChunkSize)
	sequence := 0
	for {
		// Read from the file into the buffer
		// Read chunk from response body
		index, err := resp.Body.Read(buffer)
		if err != nil && err != io.EOF {
			return 0, err
		}
		if index == 0 {
			break
		}

		// Process the chunk of data as needed
		// The buffer[:n] contains the actual data read from the file
		if err = u.processFragmentData(ctx, url, TokenId, sequence, buffer[:index]); err != nil {
			logger.AtLog.Logger.Error("Error processing data", zap.Error(err), zap.String("filePath", url), zap.String("TokenId", TokenId))
			return 0, err
		}

		// increment sequence
		sequence++
	}
	logger.AtLog.Logger.Info("File fragmented", zap.String("filePath", url), zap.String("TokenId", TokenId), zap.Int("total fragments", sequence))
	return sequence, nil
}

func (u Usecase) processFragmentData(ctx context.Context, filePath, tokenID string, sequence int, data []byte) error {

	file, err := u.Repo.FindTokenFileFragment(ctx, tokenID, sequence)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			file = &entity.TokenFileFragment{
				TokenId:  tokenID,
				Sequence: sequence,
				Status:   entity.FileFragmentStatusPending,
				FilePath: filePath,
				Data:     data,
			}
			return u.Repo.InsertFileFragment(ctx, file)
		}
		logger.AtLog.Logger.Error("Error finding token_file_fragment", zap.Error(err), zap.String("TokenId", tokenID), zap.Int("sequence", sequence))
		return err
	}

	return nil
}

func (u Usecase) CreateFileFromFragments(ctx context.Context, tokenID string, filePath string) error {

	defer func() {
		if err := recover(); err != nil {
			logger.AtLog.Logger.Error("Error creating file", zap.Error(err.(error)), zap.String("filePath", filePath), zap.String("TokenId", tokenID))
			// remove file path if exist
			_, err := os.Stat(filePath)
			if err != nil && os.IsNotExist(err) {
				return
			}
			os.Remove(filePath)
		}
	}()

	file, err := os.Create(filePath)
	if err != nil {
		logger.AtLog.Logger.Error("Error creating file", zap.Error(err), zap.String("filePath", filePath), zap.String("TokenId", tokenID))
		return err
	}
	defer file.Close()

	for page := 1; ; page++ {
		fragments, err := u.Repo.FindTokenFileFragments(ctx, repository.TokenFileFragmentFileter{
			TokenID:  tokenID,
			Page:     page,
			PageSize: 10,
		})
		if err != nil {
			os.Remove(filePath)
			logger.AtLog.Logger.Error("Error finding token_file_fragment", zap.Error(err), zap.String("TokenId", tokenID))
			return err
		}
		if len(fragments) == 0 {
			break
		}
		for _, fragment := range fragments {
			_, err = file.Write(fragment.Data)
			if err != nil {
				os.Remove(filePath)
				logger.AtLog.Logger.Error("Error writing file", zap.Error(err), zap.String("TokenId", tokenID))
				return err
			}
		}

	}
	logger.AtLog.Info("File created", zap.String("filePath", filePath), zap.String("TokenId", tokenID))
	return nil
}

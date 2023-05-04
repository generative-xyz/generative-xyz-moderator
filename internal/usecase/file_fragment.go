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
			Status:   entity.FileFragmentStatusCreated,
			PageSize: 10,
		})
		if err != nil {
			logger.AtLog.Logger.Error("Error finding file fragments", zap.Error(err))
			break
		}

		for _, file := range fragments {
			storeAddress, err := u.StoreFileInTC(&file)
			if err != nil {
				logger.AtLog.Logger.Error("Error storing file in blockchain", zap.Error(err), zap.String("TokenId", file.TokenId), zap.Int("sequence", file.Sequence))
			} else {
				now := time.Now()
				file.TxStoreNft = *storeAddress
				file.Status = entity.FileFragmentStatusProcessing
				file.UploadedAt = &now
				u.Repo.UpdateFileFragmentStatus(ctx, file.TokenId, map[string]interface{}{
					"status":       entity.FileFragmentStatusProcessing,
					"tx_store_nft": *storeAddress,
					"uploaded_at":  time.Now(),
				})
				var txSendAddress *string = nil
				for {
					var err error
					txSendAddress, err = u.getTxSendNft(&file)
					if err != nil {
						logger.AtLog.Logger.Error("Error storing file in blockchain", zap.Error(err), zap.String("TokenId", file.TokenId), zap.Int("sequence", file.Sequence))
					}
					if txSendAddress != nil {
						file.TxSendNft = *txSendAddress
						u.Repo.UpdateFileFragmentStatus(ctx, file.TokenId, map[string]interface{}{
							"tx_send_nft": file.TxSendNft,
						})
						break
					}
					time.Sleep(5 * time.Second)
				}

			}
		}
	}

	for page := 1; ; page++ {
		fragments, err := u.Repo.FindTokenFileFragments(ctx, repository.TokenFileFragmentFileter{
			Page:     page,
			Status:   entity.FileFragmentStatusProcessing,
			PageSize: 10,
		})
		if err != nil {
			logger.AtLog.Logger.Error("Error finding token file fragments", zap.Error(err))
			break
		}
		var wg sync.WaitGroup
		for _, file := range fragments {
			wg.Add(1)
			go func(file *entity.TokenFileFragment) {
				success, err := u.checkUploadDone(file)
				if err != nil {
					logger.AtLog.Logger.Error("Error storing file in blockchain", zap.Error(err), zap.String("TokenId", file.TokenId), zap.Int("sequence", file.Sequence))
				} else if success {
					u.Repo.UpdateFileFragmentStatus(ctx, file.TokenId, map[string]interface{}{
						"status": entity.FileFragmentStatusDone,
					})
				}
			}(&file)
		}
	}
}

func (u Usecase) checkUploadDone(file *entity.TokenFileFragment) (bool, error) {
	return true, nil
}

func (u Usecase) StoreFileInTC(file *entity.TokenFileFragment) (*string, error) {

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
		return nil, err
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, _ := publicKey.(*ecdsa.PublicKey)
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := u.TcClient.PendingNonceAt(context.Background(), fromAddress)

	if err != nil {
		return nil, err
	}

	gasPrice, err := u.TcClient.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}

	fmt.Println("gasPrice: ", gasPrice)

	chainID, err := u.TcClient.NetworkID(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, "crypto.HexToECDSA")
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return nil, errors.Wrap(err, "crypto.HexToECDSA")
	}

	auth.Nonce = big.NewInt(int64(nonce))

	// Create a new instance of the contract with the given address and ABI
	contract, err := generative_nft_contract.NewGenerativeNftContract(fromAddress, u.TcClient.GetClient())
	if err != nil {
		return nil, errors.Wrap(err, "NewGenerativeNftContract")
	}

	tokenIdInt := new(big.Int)
	tokenIdInt, ok := tokenIdInt.SetString(file.TokenId, 10)
	if !ok {
		return nil, fmt.Errorf("error converting token id to big int")
	}

	tx, err := contract.Store(auth, tokenIdInt, big.NewInt(int64(file.Sequence)), file.Data)

	if err != nil {
		return nil, errors.Wrap(err, "contract.Mint")
	}

	TxSendNft := tx.Hash().Hex()
	return &TxSendNft, nil
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

func (u Usecase) FragmentFile(ctx context.Context, TokenId, filePath string) (int, error) {

	_, err := u.Repo.FindTokenByTokenID(TokenId)
	if err != nil {
		logger.AtLog.Logger.Error("Error finding token_uri", zap.Error(err), zap.String("TokenId", TokenId))
		return 0, err
	}

	file, err := os.Open(filePath)
	if err != nil {
		logger.AtLog.Logger.Error("Error opening file", zap.Error(err), zap.String("filePath", filePath), zap.String("TokenId", TokenId))
		return 0, err
	}
	defer file.Close()

	// Create a buffer to read chunks of 350KB at a time
	buffer := make([]byte, ChunkSize)
	sequence := 0
	for {
		// Read from the file into the buffer
		index, err := file.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			logger.AtLog.Logger.Error("Error reading file", zap.Error(err), zap.String("filePath", filePath), zap.String("TokenId", TokenId))
			return 0, err
		}

		// Process the chunk of data as needed
		// The buffer[:n] contains the actual data read from the file
		if err = u.processFragmentData(ctx, filePath, TokenId, sequence, buffer[:index]); err != nil {
			logger.AtLog.Logger.Error("Error processing data", zap.Error(err), zap.String("filePath", filePath), zap.String("TokenId", TokenId))
			return 0, err
		}

		// increment sequence
		sequence++
	}
	logger.AtLog.Logger.Info("File fragmented", zap.String("filePath", filePath), zap.String("TokenId", TokenId), zap.Int("total fragments", sequence))
	return sequence, nil
}

func (u Usecase) processFragmentData(ctx context.Context, filePath, tokenID string, sequence int, data []byte) error {

	file, err := u.Repo.FindTokenFileFragment(ctx, tokenID, sequence)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			file = &entity.TokenFileFragment{
				TokenId:  tokenID,
				Sequence: sequence,
				Status:   entity.FileFragmentStatusCreated,
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

package usecase

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"io"
	"os"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/repository"
	"rederinghub.io/utils/encrypt"
	"rederinghub.io/utils/logger"
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

		noun, err := u.GetTCNonce()
		if err != nil {
			logger.AtLog.Logger.Error("Error getting nonce", zap.Error(err))
			break
		}

		for index, file := range fragments {
			address, err := u.StoreFileInTC(&file, uint64(index+1)+noun)
			if err != nil {
				logger.AtLog.Logger.Error("Error storing file in blockchain", zap.Error(err), zap.String("TokenId", file.TokenId), zap.Int("sequence", file.Sequence))
			} else {
				u.Repo.UpdateFileFragmentStatus(ctx, file.TokenId, map[string]interface{}{
					"status":        entity.FileFragmentStatusProcessing,
					"store_address": address,
					"uploaded_at":   time.Now(),
				})
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
				success, err := u.CheckStoreStatus(file)
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

func (u Usecase) StoreFileInTC(file *entity.TokenFileFragment, nonce uint64) (string, error) {
	panic("Not implemented")
}

func (u Usecase) CheckStoreStatus(file *entity.TokenFileFragment) (bool, error) {
	panic("Not implemented")
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

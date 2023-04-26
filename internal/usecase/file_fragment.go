package usecase

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"io"
	"os"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/repository"
	"rederinghub.io/utils/logger"
)

const ChunkSize = 350 * 1024

func (u Usecase) JobFragmentBigFile() {
	panic("implement me")
}

func (u Usecase) FragmentFile(ctx context.Context, TokenId, filePath string) error {

	_, err := u.Repo.FindTokenByTokenID(TokenId)
	if err != nil {
		logger.AtLog.Logger.Error("Error finding token_uri", zap.Error(err), zap.String("TokenId", TokenId))
		return err
	}

	file, err := os.Open(filePath)
	if err != nil {
		logger.AtLog.Logger.Error("Error opening file", zap.Error(err), zap.String("filePath", filePath), zap.String("TokenId", TokenId))
		return err
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
			return err
		}

		// Process the chunk of data as needed
		// The buffer[:n] contains the actual data read from the file
		if err = u.processFragmentData(ctx, filePath, TokenId, sequence, buffer[:index]); err != nil {
			logger.AtLog.Logger.Error("Error processing data", zap.Error(err), zap.String("filePath", filePath), zap.String("TokenId", TokenId))
			return err
		}

		// increment sequence
		sequence++
	}
	logger.AtLog.Logger.Info("File fragmented", zap.String("filePath", filePath), zap.String("TokenId", TokenId), zap.Int("total fragments", sequence))
	return nil
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

	return fmt.Errorf("file fragment already exists")
}

func (u Usecase) CreateFileFromFragment(ctx context.Context, tokenID string, filePath string) error {

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

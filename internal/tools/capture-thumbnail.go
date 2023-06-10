package tools

import (
	"encoding/json"
	"fmt"
	"os"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/repository"
	"rederinghub.io/internal/usecase"
	"rederinghub.io/utils/helpers"
	"strings"
	"sync"
)

func CaptureThumbnail(repo *repository.Repository, uc *usecase.Usecase, projectID string) {
	tokens, err := repo.GetAllTokensByProjectID(projectID)
	if err != nil {
		return
	}

	tokenCap := []entity.TokenUri{}
	for _, tok := range tokens {

		if tok.OrderInscriptionIndex == 226 ||
			tok.OrderInscriptionIndex == 231 ||
			tok.OrderInscriptionIndex == 196 ||
			tok.OrderInscriptionIndex == 197 ||
			tok.OrderInscriptionIndex == 199 ||
			tok.OrderInscriptionIndex == 200 ||
			tok.OrderInscriptionIndex == 201 ||
			tok.OrderInscriptionIndex == 202 ||
			tok.OrderInscriptionIndex == 203 ||
			tok.OrderInscriptionIndex == 204 ||
			tok.OrderInscriptionIndex == 205 ||
			tok.OrderInscriptionIndex == 216 ||
			tok.OrderInscriptionIndex == 218 ||
			tok.OrderInscriptionIndex == 222 ||
			tok.OrderInscriptionIndex == 224 ||
			tok.OrderInscriptionIndex == 225 ||
			tok.OrderInscriptionIndex == 227 ||
			tok.OrderInscriptionIndex == 198 {
			tokenCap = append(tokenCap, tok)
		}
	}

	type dataChan struct {
		TokenID   string
		Thumbnail string
		ProjectID string
		Err       error
	}

	respArr := make(chan dataChan, len(tokenCap))
	inChan := make(chan entity.TokenUri, len(tokenCap))
	var wg sync.WaitGroup

	for _, _ = range tokenCap {
		go func(wg *sync.WaitGroup, inChan chan entity.TokenUri, outChan chan dataChan) {
			defer wg.Done()

			tok := <-inChan
			resp, err := uc.RunAndCap(&tok)
			outChan <- dataChan{
				TokenID:   tok.TokenID,
				Thumbnail: resp.Thumbnail,
				ProjectID: tok.ProjectID,
				Err:       err,
			}

		}(&wg, inChan, respArr)
	}

	for i, nft := range tokenCap {
		wg.Add(1)
		inChan <- nft
		if i%10 == 0 && i > 0 {
			wg.Wait()
		}
	}

	captured := make(map[string]string)
	for _, _ = range tokenCap {
		outFromChan := <-respArr
		key := fmt.Sprintf("%s-%s", outFromChan.ProjectID, outFromChan.TokenID)
		captured[key] = outFromChan.Thumbnail
	}

	helpers.CreateFile("token-thumbnail.json", captured)
}

func UpdateThumbnail(repo *repository.Repository, uc *usecase.Usecase, fileName string) {

	bytes, err := os.ReadFile(fileName)
	if err != nil {
		return
	}

	respArr := make(map[string]string)
	err = json.Unmarshal(bytes, &respArr)
	if err != nil {
		return
	}

	for key, value := range respArr {
		arr := strings.Split(key, "-")
		//contract := arr[0]
		tokenID := arr[1]
		projectID := arr[0]
		thumbnail := value

		repo.UpdateTokenThumbnailByTokenId(projectID, tokenID, thumbnail)
	}

	return
}

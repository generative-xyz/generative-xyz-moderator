package tools

import (
	"encoding/json"
	"fmt"
	"os"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/repository"
	"rederinghub.io/internal/usecase"
	"rederinghub.io/utils/helpers"
	"sync"
)

func CaptureThumbnail(repo *repository.Repository, uc *usecase.Usecase, projectID string) {
	tokens, err := repo.GetAllTokensByProjectID(projectID)
	if err != nil {
		return
	}

	tokenCap := []entity.TokenUri{}
	for _, tok := range tokens {
		tokenCap = append(tokenCap, tok)
	}

	type dataChan struct {
		TokenID   string
		Thumbnail string
		ProjectID string
		Err       error
		Attrs     []entity.TokenUriAttr
		AttrStrs  []entity.TokenUriAttrStr
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
				Attrs:     resp.Traits,
				AttrStrs:  resp.TraitsStr,
				Err:       err,
			}

		}(&wg, inChan, respArr)
	}

	maxProcess := 2
	for i, nft := range tokenCap {
		wg.Add(1)
		inChan <- nft
		if i%maxProcess == 0 && i > 0 || i == len(tokenCap)-1 {
			wg.Wait()
		}
	}

	captured := make(map[string]string)
	for _, _ = range tokenCap {
		outFromChan := <-respArr
		key := fmt.Sprintf("%s-%s", outFromChan.ProjectID, outFromChan.TokenID)
		captured[key] = outFromChan.Thumbnail
		repo.UpdateTokenThumbnailByTokenId(outFromChan.ProjectID, outFromChan.TokenID, outFromChan.Thumbnail, outFromChan.Attrs, outFromChan.AttrStrs)
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

	//for key, value := range respArr {
	//	arr := strings.Split(key, "-")
	//	//contract := arr[0]
	//	//tokenID := arr[1]
	//	//projectID := arr[0]
	//	//thumbnail := value
	//
	//	//repo.UpdateTokenThumbnailByTokenId(projectID, tokenID, thumbnail)
	//}

	return
}

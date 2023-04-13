package usecase

import (
	"encoding/json"
	"os"
	"time"

	"github.com/davecgh/go-spew/spew"
	"rederinghub.io/internal/entity"
	"rederinghub.io/utils/helpers"
)

type ChangeName struct {
	TokenID      string
	AnimationURL string
	Thumbnail    string
}

type ChangeNameChan struct {
	Err  error
	Data *ChangeName
}

func (u Usecase) GetTokenArtworkName() {
	bkFileName := "bk-1002573.json"
	tokens := []entity.TokenUri{}
	fc, err := os.ReadFile(bkFileName)
	if err != nil {
		tokens, err = u.Repo.GetAllTokensByProjectID("1002573")
		if err != nil {
			return
		}

		helpers.CreateFile(bkFileName, tokens)
	}

	err = json.Unmarshal(fc, &tokens)
	if err != nil {
		return
	}

	resp := []ChangeName{}
	i := 0

	chanChangeName := make(chan ChangeNameChan, len(tokens))
	for _, token := range tokens {

		go func(token entity.TokenUri, chanChangeName chan ChangeNameChan) {
			tmp := &ChangeName{}
			imageURL := token.AnimationURL

			defer func() {
				chanChangeName <- ChangeNameChan{
					Err:  err,
					Data: tmp,
				}
			}()

			capResp, err := u.RunAndCap(&token)
			if err != nil {
				return
			}

			tmp.TokenID = token.TokenID
			tmp.AnimationURL = imageURL
			tmp.Thumbnail = capResp.Thumbnail
		}(token, chanChangeName)

		if i > 0 && i%10 == 0 {
			time.Sleep(500 * time.Millisecond)
		}

		i++

	}

	for _, _ = range tokens {

		dataFromChan := <-chanChangeName
		if dataFromChan.Err != nil {
			continue
		}

		resp = append(resp, *dataFromChan.Data)
	}

	spew.Dump(resp)
	helpers.CreateFile("new-inscriptionID.json", resp)
}

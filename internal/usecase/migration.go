package usecase

import (
	"encoding/json"
	"os"

	"github.com/davecgh/go-spew/spew"
	"rederinghub.io/internal/entity"
	"rederinghub.io/utils/helpers"
)
 
func (u Usecase) GetTokenArtworkName(){
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

	
	spew.Dump(tokens)
}
package usecase

import (
	"encoding/json"
	"os"
	"strconv"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"rederinghub.io/internal/entity"
	"rederinghub.io/utils/helpers"
)

type ChangeName struct {
	TokenID              string
	OrderInsciptionID    int
	NewOrderInsciptionID int
	ArtworkName string
	AnimationURL string
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
	for _, token := range tokens {
		tmp := ChangeName{}
		artworkName := ""

		
		imageURL := token.AnimationURL
		an := strings.ReplaceAll(imageURL, "https://cdn.generative.xyz/btc-projects/aiseries:perceptrons-52561678/Perceptrons/","")
		an = strings.ReplaceAll(an, ".html","")
		aID, err := strconv.Atoi(an)
		if err != nil {
			return
		}
		aID ++

		tmp.TokenID = token.TokenID
		tmp.OrderInsciptionID = token.OrderInscriptionIndex
		tmp.NewOrderInsciptionID = aID
		tmp.ArtworkName = artworkName
		tmp.AnimationURL = imageURL


		resp = append(resp, tmp)
		
	}

	spew.Dump(resp)
	helpers.CreateFile("new-inscriptionID.json", resp)
}

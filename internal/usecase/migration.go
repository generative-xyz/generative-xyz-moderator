package usecase

import (
	"context"
	"encoding/json"
	"os"
	"strconv"
	"strings"

	"time"

	"github.com/chromedp/chromedp"
	"github.com/davecgh/go-spew/spew"
	"rederinghub.io/internal/entity"
	"rederinghub.io/utils/helpers"
)

type ChangeName struct {
	TokenID              string
	OrderInsciptionID    int
	NewOrderInsciptionID int
	ArtworkName string
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

		opts := append(chromedp.DefaultExecAllocatorOptions[:],
			//chromedp.ExecPath("google-chrome"),
			chromedp.Flag("headless", false),
			chromedp.Flag("disable-gpu", false),
			chromedp.Flag("no-first-run", true),
		)
		allocCtx, _ := chromedp.NewExecAllocator(context.Background(), opts...)
		cctx, cancel := chromedp.NewContext(allocCtx)
		 

		imageURL := token.AnimationURL

		err = chromedp.Run(cctx,
			chromedp.EmulateViewport(960, 960),
			chromedp.Navigate(imageURL),
			chromedp.Sleep(time.Second*time.Duration(25)),
			//chromedp.EvaluateAsDevTools("window.$generativeTraits", &traits),
			chromedp.EvaluateAsDevTools("$artworkName", &artworkName),
		)

		if err != nil {
			return
		}

		an := strings.ReplaceAll(artworkName, "Perceptron #","")
		aID, err := strconv.Atoi(an)
		if err != nil {
			return
		}

		tmp.TokenID = token.TokenID
		tmp.OrderInsciptionID = token.OrderInscriptionIndex
		tmp.NewOrderInsciptionID = aID
		tmp.ArtworkName = artworkName


		resp = append(resp, tmp)
		cancel()
	}

	spew.Dump(resp)
	helpers.CreateFile("new-inscriptionID.json", resp)
}

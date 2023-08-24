package usecase

import (
	"encoding/json"
	"fmt"
	"github.com/antchfx/htmlquery"
	"github.com/davecgh/go-spew/spew"
	"net/url"
	"os"
	"rederinghub.io/internal/entity"
	"rederinghub.io/utils/helpers"
	"strings"
	"sync"
)

type gmHolder struct {
	Items    []string `json:"items"`
	NextPage string   `json:"next_page_path"`
}

type parsedGmHolder struct {
	WalletAddress string
	GM            string
	Percent       string
}

func (u *Usecase) crawData(fullPath string) (*gmHolder, error) {
	var err error
	res := &gmHolder{}
	res, err = u.ReportGMHolder(fullPath)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (u *Usecase) ReportGMHolders() {
	tmpFile := "report-tmp.json"
	b, err := os.ReadFile(tmpFile)
	fullData := []string{}

	if err != nil {
		domain := "https://explorer.trustless.computer"
		queries := url.Values{}
		queries["type"] = []string{"JSON"}
		done := make(chan bool)

		fullPath := fmt.Sprintf("%s/token/0x2fe8d5a64affc1d703aeca8a566f5e9faee0c003/token-holders?%s", domain, queries.Encode())

		dataChan := make(chan []string)
		go func() {
			for {
				fmt.Println("craw data from: " + fullPath)
				res, err := u.crawData(fullPath)
				if err != nil {
					done <- true
					break
				}

				dataChan <- res.Items

				if res.NextPage == "" {
					done <- true
					break
				}

				fullPath = fmt.Sprintf("%s%s&type=JSON", domain, res.NextPage)
			}
		}()
		mux := sync.Mutex{}

	L:
		for {

			select {
			case dataFChan := <-dataChan:
				mux.Lock()
				fullData = append(fullData, dataFChan...)
				mux.Unlock()

			case <-done:
				spew.Dump(done)
				break L
			}

		}

		helpers.CreateFile(tmpFile, fullData)
	} else {
		err = json.Unmarshal(b, &fullData)
		if err != nil {
			return
		}
	}

	inChan := make(chan string, len(fullData))
	outChan := make(chan parsedGmHolder, len(fullData))

	for _, _ = range fullData {
		go u.ParseData(inChan, outChan)
	}

	for _, i := range fullData {
		inChan <- i
	}

	report := []parsedGmHolder{}
	for _, _ = range fullData {
		out := <-outChan
		report = append(report, out)
	}

	helpers.CreateFile("report-gm.json", report)

}

func (u *Usecase) ParseData(input chan string, outPut chan parsedGmHolder) {
	str := <-input
	doc, err := htmlquery.Parse(strings.NewReader(str))
	if err != nil {
		outPut <- parsedGmHolder{}
		return
	}

	w := make(chan string)
	g := make(chan string)
	p := make(chan string)

	//wallet address
	go func(w chan string) {
		walletAddress := ""
		a := htmlquery.Find(doc, "//span[@data-address-hash]")
		if len(a) > 0 {
			walletAddress = htmlquery.SelectAttr(a[0], "data-address-hash")
			walletAddress = strings.ReplaceAll(walletAddress, " ", "")
			walletAddress = strings.ToLower(walletAddress)
		}

		w <- walletAddress
	}(w)

	//GM
	go func(g chan string) {
		gm := ""

		a := htmlquery.Find(doc, "//span[contains(@class,'text-dark')]")
		if len(a) > 0 {
			f := a[0].FirstChild
			if f != nil {
				f1 := *f
				gm = strings.ReplaceAll(f1.Data, " ", "")
				gm = strings.ReplaceAll(gm, "\n", "")
				gm = strings.ReplaceAll(gm, "GM", "")
			}

		}

		g <- gm
	}(g)

	//percent
	go func(p chan string) {
		percent := ""

		a := htmlquery.Find(doc, "//div[contains(@class,'flex-column')]/span/text()")
		if len(a) == 4 {
			f := a[3].Data
			percent = strings.ReplaceAll(f, "\n", "")
			percent = strings.ReplaceAll(percent, " ", "")
			percent = strings.ReplaceAll(percent, "%", "")

		}

		p <- percent
	}(p)

	outPut <- parsedGmHolder{
		WalletAddress: <-w,
		GM:            <-g,
		Percent:       <-p,
	}
}

func (u *Usecase) ReportGMHolder(path string) (*gmHolder, error) {
	method := "GET"

	bytes, _, _, err := helpers.JsonRequest(path, method, map[string]string{}, nil)
	if err != nil {
		return nil, err
	}

	resp := &gmHolder{}
	err = json.Unmarshal(bytes, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

type inscriptionInfo struct {
	InscriptionID string
	Address       string
}

type inscriptionInfoChan struct {
	Data *inscriptionInfo
	Err  error
}

func (u *Usecase) ReportPerceptronOwners() {
	//get list percentron IDs
	genAddress := "1002573"
	tokens, err := u.Repo.AnalyticsTokenUriOwner(entity.FilterTokenUris{
		GenNFTAddr: &genAddress,
	})
	if err != nil {
		return
	}

	inchan := make(chan string, len(tokens))
	outChan := make(chan inscriptionInfoChan, len(tokens))

	var wg sync.WaitGroup
	for _, _ = range tokens {
		go u.ReportPerceptronOwner(&wg, inchan, outChan)
	}

	go func(wg *sync.WaitGroup) {
		for i, tok := range tokens {
			wg.Add(1)

			inchan <- tok.TokenID
			if i > 0 && i%10 == 0 || i == len(tokens)-1 {
				wg.Wait()
			}
		}
	}(&wg)

	reportUsers := []inscriptionInfoChan{}
	walletAddress := []string{}
	for _, _ = range tokens {
		out := <-outChan
		reportUsers = append(reportUsers, out)
		walletAddress = append(walletAddress, out.Data.Address)
		l := len(walletAddress) / len(tokens)
		fmt.Println(fmt.Sprintf("=====> loading: %d %", l))
	}

	type report struct {
		Inscription []string
		Total       int
		BtcAddress  string
	}

	users, err := u.Repo.FindUserByAddresses(walletAddress)
	ethWalletAddress := make(map[string]entity.Users)
	for _, u := range users {
		ethWalletAddress[u.WalletAddressBTC] = u
	}

	data := make(map[string]report)
	for _, user := range reportUsers {
		w, ok := ethWalletAddress[user.Data.Address]
		if ok {
			key := w.WalletAddress
			item, ok1 := data[key]
			if !ok1 {
				newArray := []string{
					user.Data.InscriptionID,
				}
				data[w.WalletAddress] = report{
					Inscription: newArray,
					BtcAddress:  user.Data.Address,
					Total:       len(newArray),
				}
			} else {
				a := append(item.Inscription, user.Data.InscriptionID)
				data[key] = report{
					Inscription: a,
					BtcAddress:  user.Data.Address,
					Total:       len(a),
				}
			}
			fmt.Println("Calculating for: " + key + fmt.Sprintf("%d", data[key].Total))
		}

	}

	helpers.CreateFile("report-perceptrons.json", data)

}

func (u *Usecase) ReportPerceptronOwner(wg *sync.WaitGroup, input chan string, output chan inscriptionInfoChan) {

	var err error
	resp := &inscriptionInfo{}
	inscription := <-input
	defer wg.Done()

	defer func() {
		output <- inscriptionInfoChan{
			Err:  err,
			Data: resp,
		}
	}()

	cacheKey := fmt.Sprintf("perceptron.%s", inscription)
	cached, err := u.Cache.GetData(cacheKey)
	if err != nil || cached == nil {
		fullUrl := fmt.Sprintf("https://dev.generativeexplorer.com/api/inscription/%s", inscription)
		fmt.Println("Craw: " + fullUrl)

		bytes, _, _, err := helpers.JsonRequest(fullUrl, "GET", map[string]string{}, nil)
		if err != nil {
			return
		}

		err = json.Unmarshal(bytes, resp)
		if err != nil {
			return
		}

		resp.InscriptionID = inscription

		u.Cache.SetDataWithExpireTime(cacheKey, resp, 86400)
	} else {
		bytes := []byte(*cached)
		err = json.Unmarshal(bytes, resp)
	}
}

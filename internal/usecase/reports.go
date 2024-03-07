package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/antchfx/htmlquery"
	"github.com/davecgh/go-spew/spew"
	"math"
	"math/big"
	"net/url"
	"os"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
	"rederinghub.io/utils/helpers"
	"sort"
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

func (u *Usecase) ReportGMHolders(wg *sync.WaitGroup, domain string, gmAddress string, level string) {
	defer wg.Done()

	tmpFile := fmt.Sprintf("report-%s-tmp.json", level)
	b, err := os.ReadFile(tmpFile)
	fullData := []string{}

	//defer os.Remove(tmpFile)

	if err != nil {
		queries := url.Values{}
		queries["type"] = []string{"JSON"}
		done := make(chan bool)

		fullPath := fmt.Sprintf("%s/token/%s/token-holders?%s", domain, gmAddress, queries.Encode())

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

	helpers.CreateFile(fmt.Sprintf("report-%s-gm.json", level), report)
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
	}

	type report struct {
		Inscriptions []string
		Total        int
		BtcAddress   string
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
					Inscriptions: newArray,
					BtcAddress:   user.Data.Address,
					Total:        len(newArray),
				}
			} else {
				a := append(item.Inscriptions, user.Data.InscriptionID)
				data[key] = report{
					Inscriptions: a,
					BtcAddress:   user.Data.Address,
					Total:        len(a),
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

type revenueReport struct {
	WalletAddress string
	BTCAmount     float64
	BTCRate       float64
	BTCToUSD      float64
	ETHAmount     float64
	ETHRate       float64
	EthToUSD      float64
	USDAmount     float64
}

func (u *Usecase) ReportUserRevenue(btcRate, ethRate float64) {
	maxProcess := 5

	artistTmpFile := "tmp-artist.json"
	s2ndTmpFile := "tmp-artist-2nd.json"
	minterTmpFile := "tmp-minter.json"
	minterBtcOldTmpFile := "tmp-minter-btc-old.json"
	minterEthOldTmpFile := "tmp-minter-eth-old.json"
	b2ndTmpFile := "tmp-buyer-2ndsale.json"

	report := []revenueReport{}
	reportMinter := []revenueReport{}
	reportMinterOldData := []revenueReport{}
	reportMinterETHOldData := []revenueReport{}

	data := map[string][]entity.ReportArtist{}
	minterData := map[string][]entity.ReportArtist{}
	minterBtcOLDData := map[string][]entity.ReportArtist{}
	minterEthOLDData := map[string][]entity.ReportArtist{}
	dataBuyer2ndSale := []*entity.Report2ndSale{}
	dataSeller2ndSale := []*entity.Report2ndSale{}

	b, err := os.ReadFile(artistTmpFile)
	s1, err := os.ReadFile(s2ndTmpFile)

	b2, err := os.ReadFile(b2ndTmpFile)
	mt, err := os.ReadFile(minterTmpFile)
	mtoldBTC, err := os.ReadFile(minterBtcOldTmpFile)
	mtoldEth, err := os.ReadFile(minterEthOldTmpFile)

	if err != nil {

		usersFChan := make(chan []*entity.Users)
		done := make(chan bool)
		wg := sync.WaitGroup{}

		userRevenueChan := make(chan *map[string][]entity.ReportArtist)
		seller2ndSaleChan := make(chan *[]*entity.Report2ndSale)

		minterChan := make(chan *map[string][]entity.ReportArtist)
		minterBTCOLDChan := make(chan *map[string][]entity.ReportArtist)
		minterEthOLDChan := make(chan *map[string][]entity.ReportArtist)
		buyer2ndSaleChan := make(chan *[]*entity.Report2ndSale)

		go u.ReportedUsers(usersFChan, done)
		count := 1

	LOOP:
		for {
			select {
			case users := <-usersFChan:
				{
					fmt.Println("Select =======>>>>", len(users), count)
					wg.Add(6)

					go u.CalculateUserRevenue(&wg, users, userRevenueChan)
					go u.CalculateMinter(&wg, users, minterChan)
					go u.CalculateMinterBTCOld(&wg, users, minterBTCOLDChan)
					go u.CalculateMinterETHOld(&wg, users, minterEthOLDChan)
					go u.CalculateBuyer2ndSale(&wg, "buyer", users, buyer2ndSaleChan)
					go u.CalculateBuyer2ndSale(&wg, "seller_address", users, seller2ndSaleChan)

					if count > 0 && count%maxProcess == 0 {
						//fmt.Println("[Wait] Select ===========> ", len(users), count)
						wg.Wait()
					}
					count++
				}
			case isDone := <-done:
				{
					fmt.Println("isDone =======>>>>", isDone)
					break LOOP
				}
			case urv := <-userRevenueChan:
				{
					if urv != nil {
						items := *urv
						fmt.Println(fmt.Sprintf("==> append  %d items", len(items)))
						//data = append(data, i)
						for k, i := range items {
							data[k] = i
						}
					}
				}
			case buyer2ndSale := <-buyer2ndSaleChan:
				{
					if buyer2ndSale != nil {
						items := *buyer2ndSale
						fmt.Println(fmt.Sprintf("[Buyer 2nd sale]==> append  %d items", len(dataBuyer2ndSale)))
						//data = append(data, i)

						for _, i := range items {
							z := u.ConvertAmount(i.Amount, 8)
							zz, _ := z.Float64()
							i.AmountUSD = zz * btcRate
							i.Amount = zz
							dataBuyer2ndSale = append(dataBuyer2ndSale, i)
						}
					}
				}
			case seller2ndSale := <-seller2ndSaleChan:
				{
					if seller2ndSale != nil {
						items := *seller2ndSale
						fmt.Println(fmt.Sprintf("[Seller 2nd sale]==> append  %d items", len(dataSeller2ndSale)))
						//data = append(data, i)

						for _, i := range items {
							z := u.ConvertAmount(i.Amount, 8)
							zz, _ := z.Float64()
							i.AmountUSD = zz * btcRate
							i.Amount = zz
							dataSeller2ndSale = append(dataSeller2ndSale, i)
						}
					}
				}
			case minter := <-minterChan:
				{
					if minter != nil {
						items := *minter
						fmt.Println(fmt.Sprintf("==> append  %d items", len(items)))
						//data = append(data, i)
						for k, i := range items {
							minterData[k] = i
						}
					}
				}
			case minterBtcOld := <-minterBTCOLDChan:
				{
					if minterBtcOld != nil {
						items := *minterBtcOld
						fmt.Println(fmt.Sprintf("==> append  %d items", len(items)))
						//data = append(data, i)
						for k, i := range items {
							minterBtcOLDData[k] = i
						}
					}
				}
			case minterEthOLD := <-minterEthOLDChan:
				{
					if minterEthOLD != nil {
						items := *minterEthOLD
						fmt.Println(fmt.Sprintf("==> append  %d items", len(items)))
						//data = append(data, i)
						for k, i := range items {
							minterEthOLDData[k] = i
						}
					}
				}
			}
		}

		//cache the data for the next
		helpers.CreateFile(artistTmpFile, data)
		helpers.CreateFile(minterTmpFile, minterData)
		helpers.CreateFile(minterBtcOldTmpFile, minterBtcOLDData)
		helpers.CreateFile(minterEthOldTmpFile, minterEthOLDData)
		helpers.CreateFile(b2ndTmpFile, dataBuyer2ndSale)
		helpers.CreateFile(s2ndTmpFile, dataSeller2ndSale)

	} else {
		err := json.Unmarshal(b, &data)
		if err != nil {
			return
		}

		err = json.Unmarshal(b2, &dataBuyer2ndSale)
		if err != nil {
			return
		}

		err = json.Unmarshal(s1, &dataSeller2ndSale)
		if err != nil {
			return
		}

		err = json.Unmarshal(mt, &minterData)
		if err != nil {
			return
		}

		err = json.Unmarshal(mtoldBTC, &minterBtcOLDData)
		if err != nil {
			return
		}

		err = json.Unmarshal(mtoldEth, &minterEthOLDData)
		if err != nil {
			return
		}
	}

	for w, info := range data {
		usd := float64(0)
		re := revenueReport{}
		for _, ptype := range info {
			z := u.ConvertAmount(ptype.Amount, 8)
			v := float64(0)
			tmp, _ := z.Float64()
			if ptype.PayType == "eth" {
				v = tmp * ethRate
				re.ETHAmount = tmp
				re.ETHRate = ethRate
				re.EthToUSD = v
			} else {
				v = tmp * btcRate
				re.BTCAmount = tmp
				re.BTCRate = btcRate
				re.BTCToUSD = v
			}

			usd += v

		}

		re.WalletAddress = w
		re.USDAmount = usd
		report = append(report, re)
	}

	for w, info := range minterData {
		usd := float64(0)
		re := revenueReport{}
		for _, ptype := range info {
			decimal := 8
			if ptype.PayType == "eth" {
				decimal = 18
			}

			z := u.ConvertAmount(ptype.Amount, decimal)
			v := float64(0)
			tmp, _ := z.Float64()
			if ptype.PayType == "eth" {
				v = tmp * ethRate
				re.ETHAmount = tmp
				re.ETHRate = ethRate
				re.EthToUSD = v
			} else {
				v = tmp * btcRate
				re.BTCAmount = tmp
				re.BTCRate = btcRate
				re.BTCToUSD = v
			}

			usd += v

		}

		re.WalletAddress = w
		re.USDAmount = usd
		reportMinter = append(reportMinter, re)
	}

	for w, info := range minterBtcOLDData {
		usd := float64(0)
		re := revenueReport{}
		for _, ptype := range info {
			z := u.ConvertAmount(ptype.Amount, 8)
			tmp, _ := z.Float64()
			v := tmp * btcRate
			re.BTCAmount = tmp
			re.BTCRate = btcRate
			re.BTCToUSD = v
			usd += v
		}

		re.WalletAddress = w
		re.USDAmount = usd
		reportMinterOldData = append(reportMinterOldData, re)
	}

	for w, info := range minterEthOLDData {
		usd := float64(0)
		re := revenueReport{}
		for _, ptype := range info {
			z := u.ConvertAmount(ptype.Amount, 18)
			tmp, _ := z.Float64()
			v := tmp * btcRate
			re.BTCAmount = tmp
			re.BTCRate = btcRate
			re.BTCToUSD = v
			usd += v
		}

		re.WalletAddress = w
		re.USDAmount = usd
		reportMinterETHOldData = append(reportMinterETHOldData, re)
	}

	reportMinter = u.MergeMinterArray(reportMinter, reportMinterOldData)
	reportMinter = u.MergeMinterArray(reportMinter, reportMinterETHOldData)

	sort.SliceStable(report, func(i, j int) bool {
		return report[i].USDAmount > report[j].USDAmount
	})
	sort.SliceStable(reportMinter, func(i, j int) bool {
		return reportMinter[i].USDAmount > reportMinter[j].USDAmount
	})
	sort.SliceStable(dataBuyer2ndSale, func(i, j int) bool {
		return dataBuyer2ndSale[i].Amount > dataBuyer2ndSale[j].Amount
	})
	sort.SliceStable(dataSeller2ndSale, func(i, j int) bool {
		return dataSeller2ndSale[i].Amount > dataSeller2ndSale[j].Amount
	})

	fmt.Println(" =======> [Exit] <=======")

	sellerR := u.MergeReportArray(report, dataSeller2ndSale)
	buyerR := u.MergeReportArray(reportMinter, dataBuyer2ndSale)

	helpers.CreateFile("report-seller.json", sellerR)
	helpers.CreateFile("report-buyer.json", buyerR)

}

type reportSale struct {
	WalletAddress string
	FirstSale     float64
	SecondSale    float64
	Total         float64
}

func (u *Usecase) MergeMinterArray(arr1 []revenueReport, arr2 []revenueReport) []revenueReport {
	resp := []revenueReport{}
	for _, i := range arr1 {
		resp = append(resp, i)
	}

	for _, i := range arr2 {
		isExisted, index := u.checkExistedInReport(i.WalletAddress, resp)
		if isExisted {
			resp[index].USDAmount = i.USDAmount + resp[index].USDAmount
		} else {
			resp = append(resp, i)
		}

	}

	return resp
}

func (u *Usecase) MergeReportArray(arr1 []revenueReport, arr2 []*entity.Report2ndSale) []*reportSale {
	resp := []*reportSale{}

	a1 := make(map[string]float64)
	a2 := make(map[string]float64)

	for _, i := range arr1 {
		a1[strings.ToLower(i.WalletAddress)] = i.USDAmount
	}

	for _, i := range arr2 {
		a2[strings.ToLower(i.WalletAddress)] = i.AmountUSD
	}

	for k, i := range a1 {
		item := &reportSale{
			WalletAddress: k,
			FirstSale:     i,
			SecondSale:    0,
			Total:         i,
		}
		resp = append(resp, item)
	}

	for k, i := range a2 {
		secondSale := u.checkWalletExistedInReport(k, resp)
		if secondSale != nil {
			//update
			secondSale.SecondSale = i
			secondSale.Total = i + secondSale.FirstSale

		} else {
			secondSale = &reportSale{
				WalletAddress: k,
				FirstSale:     0,
				SecondSale:    i,
				Total:         i,
			}
			resp = append(resp, secondSale)
		}
	}

	sort.SliceStable(resp, func(i, j int) bool {
		return resp[i].Total > resp[j].Total
	})

	return resp
}

func (u *Usecase) checkWalletExistedInReport(wallet string, arr []*reportSale) *reportSale {
	for _, i := range arr {
		if i.WalletAddress == wallet {
			return i
		}
	}

	return nil
}

func (u *Usecase) checkExistedInReport(wallet string, arr []revenueReport) (bool, int) {
	for key, i := range arr {
		if i.WalletAddress == wallet {
			return true, key
		}
	}

	return false, -1
}

func (u *Usecase) ConvertAmount(amount float64, decimal int) *big.Float {
	f := big.NewFloat(amount)
	pow := math.Pow10(decimal)
	powBig := big.NewFloat(0).SetFloat64(pow)
	z := f.Quo(f, powBig)
	return z
}

func (u *Usecase) ReportedUsers(out chan []*entity.Users, done chan bool) {

	limit := int64(100)
	page := int64(1)

	defer func() {
		done <- true
	}()

	filter := structure.FilterUsers{}
	for {

		users := []*entity.Users{}
		cacheKey := fmt.Sprintf("users.limit.%d.page.%d", limit, page)
		cached, err := u.Repo.Cache.GetData(cacheKey)
		if err != nil {
			filter.BaseFilters = structure.BaseFilters{
				Limit: limit,
				Page:  page,
			}
			pagination, err := u.Repo.ListUsersWithPagination(filter)
			if err != nil {
				break
			}

			users = pagination.Result.([]*entity.Users)
			u.Cache.SetDataWithExpireTime(cacheKey, users, 86400)
		} else {
			err = json.Unmarshal([]byte(*cached), &users)
		}

		if len(users) == 0 {
			break
		}

		out <- users
		page++
	}
}

func (u *Usecase) CalculateUserRevenue(wg *sync.WaitGroup, users []*entity.Users, out chan *map[string][]entity.ReportArtist) {
	walletAddress := []string{}
	respP := new(map[string][]entity.ReportArtist)
	resp := make(map[string][]entity.ReportArtist)
	respP = &resp

	v := []entity.ReportArtist{}

	var err error

	defer func() {
		wg.Done()
		out <- respP
	}()

	for _, u := range users {
		walletAddress = append(walletAddress, u.WalletAddress)
	}

	v, err = u.Repo.AggregateUsersVolumn(walletAddress)
	if err != nil {
		return
	}

	for _, i := range v {
		key := strings.ToLower(i.WalletAddress)
		info, ok := resp[key]
		if !ok {
			a := []entity.ReportArtist{}
			a = append(a, i)
			resp[key] = a
		} else {
			info = append(info, i)
			resp[key] = info
		}
	}

}

func (u *Usecase) CalculateMinter(wg *sync.WaitGroup, users []*entity.Users, out chan *map[string][]entity.ReportArtist) {
	walletAddress := []string{}
	respP := new(map[string][]entity.ReportArtist)
	resp := make(map[string][]entity.ReportArtist)
	respP = &resp

	v := []entity.ReportArtist{}

	var err error

	defer func() {
		wg.Done()
		out <- respP
	}()

	mapBTCwl := make(map[string]string)
	for _, u := range users {
		walletAddress = append(walletAddress, u.WalletAddressBTC)
		mapBTCwl[u.WalletAddressBTC] = u.WalletAddress
	}

	v, err = u.Repo.AggregateMinterVolumn(walletAddress)
	if err != nil {
		return
	}

	for _, i := range v {
		parsedEthWl, ok := mapBTCwl[i.WalletAddress]
		if !ok {
			continue
		}

		key := strings.ToLower(parsedEthWl)
		info, ok := resp[key]
		if !ok {
			a := []entity.ReportArtist{}
			a = append(a, i)
			resp[key] = a
		} else {
			info = append(info, i)
			resp[key] = info
		}
	}

}

func (u *Usecase) CalculateMinterBTCOld(wg *sync.WaitGroup, users []*entity.Users, out chan *map[string][]entity.ReportArtist) {
	walletAddress := []string{}
	respP := new(map[string][]entity.ReportArtist)
	resp := make(map[string][]entity.ReportArtist)
	respP = &resp

	v := []entity.ReportArtist{}

	var err error

	defer func() {
		wg.Done()
		out <- respP
	}()

	mapBTCwl := make(map[string]string)
	for _, u := range users {
		if u.WalletAddressBTC == "" {
			continue
		}

		walletAddress = append(walletAddress, u.WalletAddressBTC)
		mapBTCwl[u.WalletAddressBTC] = u.WalletAddress
	}

	v, err = u.Repo.AggregateMinterBTCVolumnOld(walletAddress)
	if err != nil {
		return
	}

	for _, i := range v {
		parsedEthWl, ok := mapBTCwl[i.WalletAddress]
		if !ok {
			continue
		}

		key := strings.ToLower(parsedEthWl)
		info, ok := resp[key]
		if !ok {
			a := []entity.ReportArtist{}
			a = append(a, i)
			resp[key] = a
		} else {
			info = append(info, i)
			resp[key] = info
		}
	}

}

func (u *Usecase) CalculateMinterETHOld(wg *sync.WaitGroup, users []*entity.Users, out chan *map[string][]entity.ReportArtist) {
	walletAddress := []string{}
	respP := new(map[string][]entity.ReportArtist)
	resp := make(map[string][]entity.ReportArtist)
	respP = &resp

	v := []entity.ReportArtist{}

	var err error

	defer func() {
		wg.Done()
		out <- respP
	}()

	mapBTCwl := make(map[string]string)
	for _, u := range users {
		if u.WalletAddressBTC == "" {
			continue
		}

		walletAddress = append(walletAddress, u.WalletAddressBTC)
		mapBTCwl[u.WalletAddressBTC] = u.WalletAddress
	}

	v, err = u.Repo.AggregateMinterEthVolumnOld(walletAddress)
	if err != nil {
		return
	}

	for _, i := range v {
		parsedEthWl, ok := mapBTCwl[i.WalletAddress]
		if !ok {
			continue
		}

		key := strings.ToLower(parsedEthWl)
		info, ok := resp[key]
		if !ok {
			a := []entity.ReportArtist{}
			a = append(a, i)
			resp[key] = a
		} else {
			info = append(info, i)
			resp[key] = info
		}
	}

}

func (u *Usecase) CalculateBuyer2ndSale(wg *sync.WaitGroup, userType string, users []*entity.Users, out chan *[]*entity.Report2ndSale) {
	vp := new([]*entity.Report2ndSale)
	v := []*entity.Report2ndSale{}
	vp = &v
	walletAddress := []string{}
	mapBtcVsEth := make(map[string]string)
	var err error

	defer func() {
		wg.Done()
		out <- vp
	}()

	for _, u := range users {
		walletAddress = append(walletAddress, u.WalletAddressBTC)
		mapBtcVsEth[u.WalletAddressBTC] = u.WalletAddress
	}

	v, err = u.Repo.AggregateBuyer2ndSaleVolumn(walletAddress, userType)
	if err != nil {
		return
	}

	for _, i := range v {
		w := mapBtcVsEth[i.WalletAddressBTC]
		i.WalletAddress = w
	}
}

func (u *Usecase) ExportMagicEdend(collection string) {
	project, err := u.Repo.FindProjectByTokenID(collection)
	if err != nil {
		return
	}

	f := structure.FilterTokens{}
	genNFTAddr := collection
	cached := fmt.Sprintf("_exp.%s", genNFTAddr)
	data := []entity.ModularTokenUri{}

	u.Cache.Delete(cached)
	err = u.Cache.GetObjectData(cached, &data)
	if err != nil {
		f.GenNFTAddr = &genNFTAddr
		inscriptions, err := u.Repo.AllModularInscriptions(context.Background(), f)
		if err != nil {
			return
		}
		u.Cache.SetDataWithExpireTime(cached, inscriptions, 86400)
		data = inscriptions
	}

	type magicedenMetaAttributes struct {
		TraitType string `json:"trait_type"`
		Value     string `json:"value"`
	}

	type magicedenMeta struct {
		Name          string                    `json:"name"`
		HighResImgURL string                    `json:"high_res_img_url"`
		Attributes    []magicedenMetaAttributes `json:"attributes,omitempty"`
	}

	type magiceden struct {
		ID   string        `json:"id"`
		Meta magicedenMeta `json:"meta"`
	}

	jsonData := []magiceden{}
	for _, i := range data {

		attrs := []magicedenMetaAttributes{}

		for _, _attr := range i.ParsedAttributesStr {
			attrs = append(attrs, magicedenMetaAttributes{
				TraitType: _attr.TraitType,
				Value:     _attr.Value,
			})
		}

		jsonDataItem := magiceden{
			ID: i.TokenID,
			Meta: magicedenMeta{
				Name:          fmt.Sprintf("%s #%d", project.Name, i.OrderInscriptionIndex),
				HighResImgURL: i.Thumbnail,
				Attributes:    attrs,
			},
		}
		jsonData = append(jsonData, jsonDataItem)
	}

	helpers.CreateFile("exported.json", jsonData)
}

func (u *Usecase) CaptureThumbnails(collectionID string) {
	u.Capture(&structure.TokenImagePayload{
		TokenID:         "d7b0c6e9e8b143288973bc77dfb9ddac44534e26dd01ffce1a60134b9564efcai0",
		ContractAddress: "0x0000000000000000000000000000000000000000",
	})
}

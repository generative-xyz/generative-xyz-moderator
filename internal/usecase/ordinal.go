package usecase

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/davecgh/go-spew/spew"
	"go.uber.org/zap"
	"rederinghub.io/utils/helpers"
	"rederinghub.io/utils/logger"
)



const (
	API_URL string = `https://dev-v5.generativeexplorer.com/api/`
)

func (u Usecase) FindInscriptions() {
	totalItem := 100
	max := 10000
	inscriptions := []string{}
	for {
		if totalItem > max {
			break
		}

		data, err := u.Dev5service.Inscriptions(fmt.Sprintf("%d",totalItem))
		if err != nil {
			logger.AtLog.Logger.Error("u.Dev5service.Inscriptions", zap.Error(err))
		}
		logger.AtLog.Logger.Info("u.Dev5service.Inscriptions", zap.Any("next", data.Next))
		totalItem = data.Next

		for i := len(data.Inscriptions) - 1; i >= 0 ; i -- {
			inscriptions = append(inscriptions, data.Inscriptions[i])
		}

		//inscriptions = append(inscriptions, data.Inscriptions)

		
	}

	bytes, err := json.Marshal(inscriptions)
	if err != nil {
		return
	}

	err = helpers.CreateFile("inscription.json", bytes)
	if err != nil {
		return
	}

	//spew.Dump(inscriptions)
}

type Inscription struct {
	ID string `json:"id"`
	Meta map[string]string `json:"meta"`
}

type MetaJson struct {
	Name string `json:"name"`
	InscriptionIcon string `json:"inscription_icon"`
	Supply string `json:"supply"`
	Slug string `json:"slug"`
	Description string `json:"description"`
	TwitterLink string `json:"twitter_link"`
	DiscordLink string `json:"discord_link"`
	WebsiteLink string `json:"website_link"`
	WalletAddress string `json:"wallet_address"`
	Royalty string `json:"royalty"`
}

func (u Usecase) CreateInscriptionFiles() {
	inscriptions, err := helpers.ReadFile("inscription.json") 
	if err != nil {
		return
	}

	data := []string{}
	err = json.Unmarshal(inscriptions, &data)
	if err != nil {
		return
	}

	var wg sync.WaitGroup
	wg.Add(3)

	go func(wg *sync.WaitGroup){
		defer wg.Done()

		data100 := data[0:100]
		u.CreateData("100 items", "100", data100)
	}(&wg)
	
	go func(wg *sync.WaitGroup){
		defer wg.Done()
		data100 := data[101:1000]
		u.CreateData("1000 items", "1000", data100)
	}(&wg)
	
	go func(wg *sync.WaitGroup){
		defer wg.Done()
		data100 := data[1001:10000]
		u.CreateData("9000 items", "9000", data100)
	}(&wg)
	
	wg.Wait()
	
}

func (u Usecase) CreateData(pjName string, folderName string, data[]string) {
	data100Data := []Inscription{}
	for i, item := range data {
		data100Item :=  Inscription{}
		data100Item.ID = item
		data100Item.Meta = make(map[string]string)
		data100Item.Meta["name"] = fmt.Sprintf("#%d", i+1)
		data100Data = append(data100Data, data100Item)
		logger.AtLog.Infof("Creating data %v", data100Item)
	}

	bytes, err := json.Marshal(data100Data)
	if err != nil {
		return
	}
	err = helpers.CreateFile(folderName+"/inscriptions.json", bytes)
	if err != nil {
		return
	}
	logger.AtLog.Infof("Created file inscriptions.json")

	mtdata := MetaJson{
		Name: fmt.Sprintf("Project %s",pjName),
		InscriptionIcon: data[0],
		Supply: fmt.Sprintf("%d", len(data)),
		Slug: helpers.GenerateSlug(fmt.Sprintf("Project %s",pjName)),
	}

	byteMt, err := json.Marshal(mtdata)
	if err != nil {
		return
	}
	err = helpers.CreateFile(folderName+"/meta.json", byteMt)
	if err != nil {
		return
	}

	logger.AtLog.Infof("Created file meta.json")
	spew.Dump(len(data100Data))
}
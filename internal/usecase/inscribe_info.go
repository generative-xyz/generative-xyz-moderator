package usecase

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gocolly/colly"
	"go.mongodb.org/mongo-driver/mongo"
	"rederinghub.io/internal/entity"
)


func (u Usecase) crawlInscribeWebsite( id string) (inscriptionInfo *entity.InscribeInfo, err error) {
	defer func() {
			if r := recover(); r != nil {
					err = fmt.Errorf("was panic, id=%s recovered value: %v", id, r)
			}
	}()

	url := fmt.Sprintf("https://generativeexplorer.com/inscription/%s", id)
	dts := []string{}
	dds := []string{}
	hrefs := []string{}
	var inscribeIndex string

	c := colly.NewCollector()

	c.OnHTML("dl", func(e *colly.HTMLElement) {
		if e == nil {
			err = fmt.Errorf("something went wrong went crawl inscribe id")
			return
		}
		e.ForEach("dt", func(id int, e *colly.HTMLElement) {
			if e == nil {
				err = fmt.Errorf("something went wrong went crawl inscribe id")
				return
			}
			dts = append(dts, e.Text)
		})
		e.ForEach("dd", func(id int, e *colly.HTMLElement) {
			if e == nil {
				err = fmt.Errorf("something went wrong went crawl inscribe id")
				return
			}
			dds = append(dds, e.Text)
			hrefs = append(hrefs, e.ChildAttr("a", "href"))
		})
	})
	c.OnHTML("h1", func(e *colly.HTMLElement) {
		if e == nil {
			err = fmt.Errorf("something went wrong went crawl inscribe id")
			return
		}
		text := e.Text
		inscribeIndex = strings.Replace(text, "Inscription ", "", -1)
	})

	c.OnError(func(r *colly.Response, e error) {
		// u.Logger.Error(fmt.Sprintf("request to url %s failed", url), err.Error(), err)
		// err = e
		err = fmt.Errorf("something went wrong went crawl inscribe id")
	})


	c.Visit(url)

	if err != nil {
		return nil, err
	}

	if len(dts) != len(dds) {
		return nil, fmt.Errorf("something went wrong went crawl inscribe id %s", id)
	}

	inscribeInfo := map[string]string{}
	inscribeInfoToHref := map[string]string{}

	for i := 0; i < len(dts); i++ {
		inscribeInfo[dts[i]] = dds[i]
		inscribeInfoToHref[dts[i]] = hrefs[i]
	}

	return &entity.InscribeInfo{
		ID: inscribeInfo["id"],
		Index: inscribeIndex,
		Address: inscribeInfo["address"],
		OutputValue: inscribeInfo["output value"],
		Preview: inscribeInfoToHref["preview"],
		Content: inscribeInfoToHref["content"],
		ContentLength: inscribeInfo["content length"],
		ContentType: inscribeInfo["content type"],
		Timestamp: inscribeInfo["timestamp"],
		GenesisHeight: inscribeInfo["genesis height"],
		GenesisTransaction: inscribeInfo["genesis transaction"],
		Location: inscribeInfo["location"],
		Output: inscribeInfo["output"],
		Offset: inscribeInfo["offset"],
	}, nil
}

func (u Usecase) GetInscribeInfo( id string) (*entity.InscribeInfo, error) {
	inscribeInfo, err := u.Repo.GetInscribeInfo(id);
	if err != nil {
		// Failed to find inscribe info in database, try to crawl it from website
		inscribeInfo, err = u.crawlInscribeWebsite(id)
		if err != nil {
			return nil, err
		} else {
			// If crawl successfully, create the inscribe info
			err := u.Repo.CreateInscribeInfo(inscribeInfo)
			if err != nil {
				return nil, err
			}
		}
	}

	return inscribeInfo, nil
}

func (u Usecase) SyncInscribeInfo( id string) (*entity.InscribeInfo, bool, error) {

	updated := false
// try to find old record in mongo
	oldInscribeInfo, err := u.Repo.GetInscribeInfo(id);
	var newInscribeInfo *entity.InscribeInfo
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			// make sure oldInscribeInfo == nil
			oldInscribeInfo = nil
			newInscribeInfo, err = u.crawlInscribeWebsite(id) 
			if err != nil {
				return nil, updated, err
			}
		} else {
			return nil, updated, err		}
	}
// need an update
	if oldInscribeInfo == nil || oldInscribeInfo.Address != newInscribeInfo.Address {
		// update inscribe info document
		updated = true
		u.Repo.UpsertTokenUri(id, newInscribeInfo)
	} 

	return newInscribeInfo, updated, nil
}

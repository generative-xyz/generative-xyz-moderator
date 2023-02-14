package usecase

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly"
	"github.com/opentracing/opentracing-go"
	"rederinghub.io/internal/entity"
)


func (u Usecase) crawlInscribeWebsite(rootSpan opentracing.Span, id string) (*entity.InscribeInfo, error) {
	span, log := u.StartSpan("usecase.crawlInscribeWebsite", rootSpan)
	defer u.Tracer.FinishSpan(span, log)
	
	dts := []string{}
	dds := []string{}
	hrefs := []string{}
	var inscribeIndex string

	c := colly.NewCollector()

	c.OnHTML("dl", func(e *colly.HTMLElement) {
		e.ForEach("dt", func(id int, e *colly.HTMLElement) {
			dts = append(dts, e.Text)
		})
		e.ForEach("dd", func(id int, e *colly.HTMLElement) {
			dds = append(dds, e.Text)
			hrefs = append(hrefs, e.ChildAttr("a", "href"))
		})
	})
	
	c.OnHTML("h1", func(e *colly.HTMLElement) {
		text := e.Text
		inscribeIndex = strings.Replace(text, "Inscription ", "", -1)
	})

	c.Visit(fmt.Sprintf("https://ordinals-explorer.generative.xyz/inscription/%s", id))

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
		Timestamp: inscribeInfo["timestamp"],
		GenesisHeight: inscribeInfo["genesis height"],
		GenesisTransaction: inscribeInfo["genesis transaction"],
		Location: inscribeInfo["location"],
		Output: inscribeInfo["output"],
		Offset: inscribeInfo["offset"],
	}, nil
}

func (u Usecase) GetInscribeInfo(rootSpan opentracing.Span, id string) (*entity.InscribeInfo, error) {
	span, log := u.StartSpan("usecase.GetInscribeInfo", rootSpan)
	defer u.Tracer.FinishSpan(span, log)
	
	inscribeInfo, err := u.Repo.GetInscribeInfo(id);
	if err != nil {
		// Failed to find inscribe info in database, try to crawl it from website
		inscribeInfo, err = u.crawlInscribeWebsite(span, id)
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

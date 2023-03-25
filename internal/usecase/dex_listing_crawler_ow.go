package usecase

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/itchyny/timefmt-go"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
)

func (u Usecase) JobDexListingCrawlerOW() {
	var wg sync.WaitGroup

	wg.Add(1)

	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		err := u.crawlOWListing()
		if err != nil {
			log.Println("JobWatchPendingBTCTxSubmit watchPendingBTCTxSubmit err", err)
		}
	}(&wg)

	wg.Wait()
}

func (u Usecase) crawlOWListing() error {
	offset := 0
	retryT := 0
	for {
		if retryT >= 5 {
			return errors.New("crawlOrdinalsWalletCollection getOrdinalsWalletCollections max retry")
		}
		collections, err := getOWCollections(offset)
		if err != nil {
			log.Println("crawlOrdinalsWalletCollection getOrdinalsWalletCollections err", err)
			time.Sleep(5 * time.Second)
			retryT++
			continue
		}
		if collections == nil {
			return errors.New("crawlOrdinalsWalletCollection getOrdinalsWalletCollections collections is nil")
		}
		retryT = 0
		collectionLen := len(collections.Collections)
		for _, collection := range collections.Collections {
			if !collection.Active {
				continue
			}
			collectionItems, err := getOWCollectionItems(collection.Slug)
			if err != nil {
				log.Println("crawlOrdinalsWalletCollection getOWCollectionItems err", err)
				time.Sleep(5 * time.Second)
				continue
			}
			log.Printf("collection %v has %v items\n", collection.Slug, len(collectionItems))
			totalCollectionValue := 0
			itemsToSave := []entity.DexBTCOWInscription{}
			err = u.Repo.ClearDexBTCOWCollectionListing(collection.Slug)
			if err != nil {
				log.Println("crawlOrdinalsWalletCollection ClearDexBTCOWCollectionListing err", err)
				continue
			}
			for _, v := range collectionItems {
				if v.Escrow.BoughtAt != "" {
					log.Println("collection", collection.Slug, "item", v.ID, "is already bought")
				}
				totalCollectionValue += v.Escrow.SatoshiPrice
				itemsToSave = append(itemsToSave, entity.DexBTCOWInscription{
					CollectionSlug: collection.Slug,
					InscriptionID:  v.ID,
					Price:          v.Escrow.SatoshiPrice,
					SellerAddress:  strings.ToLower(v.Escrow.SellerAddress),
				})
			}
			err = u.Repo.CreateDexBTCOWInscriptions(itemsToSave)
			if err != nil {
				log.Println("crawlOrdinalsWalletCollection CreateDexBTCOWInscriptions err", err)
				continue
			}
			log.Println("collection", collection.Slug, "total value", totalCollectionValue)
		}
		log.Printf("got %v collections in total %v\n", offset, collections.Total)

		offset += collectionLen
		if offset >= collections.Total {
			break
		}
	}
	return nil
}

func getOWCollections(offset int) (*structure.DexListingOWCollections, error) {
	url := fmt.Sprintf("https://turbo.ordinalswallet.com/v2/collections?order=VolumeWeekDesc&offset=%v", offset)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var result structure.DexListingOWCollections
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func getOWCollectionItems(collectionSlug string) ([]structure.DexListingOWCollectionItem, error) {
	skip := isExcludeCollection(collectionSlug)
	if skip {
		return nil, nil
	}
	var result []structure.DexListingOWCollectionItem
	offset := 0
	for {
		url := fmt.Sprintf("https://turbo.ordinalswallet.com/collection/%v/inscriptions?offset=%v&order=PriceAsc&listed=true", collectionSlug, offset)

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		var respond []structure.DexListingOWCollectionItem
		err = json.Unmarshal(body, &respond)
		if err != nil {
			return nil, err
		}
		if len(respond) == 0 {
			break
		}
		result = append(result, respond...)
		offset = len(result)
	}

	return result, nil
}

func GetOWListingRaw(inscriptionID string, address string, pubkey string) (*structure.DexListingOWPurchaseRespond, error) {
	url := "https://turbo.ordinalswallet.com/wallet/purchase"

	payload := strings.NewReader("{\n\t\"inscription\": \"" + inscriptionID + "\",\n\t\"from\": \"" + address + "\",\n\t\"public_key\": \"" + pubkey + "\"\n}")

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var result structure.DexListingOWPurchaseRespond
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func submitOWPurchaseRawTx(setupTx string, purchaseTx string) (string, string, error) {
	url := "https://turbo.ordinalswallet.com/market/purchase"

	payload := strings.NewReader("{\n\t\"setup_rawtx\": \"" + setupTx + "\",\n\t\"purchase_rawtx\": \"" + purchaseTx + "\"\n}")

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return "", "", err
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", "", err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", "", err
	}
	var result structure.DexListingOWPurchaseRespond
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", "", err
	}

	return result.Purchase, result.Setup, nil
}

func isExcludeCollection(collectionSlug string) bool {
	excludeList := make(map[string]struct{})
	excludeList["sub-100k"] = struct{}{}
	excludeList["sub-10k"] = struct{}{}
	excludeList["sub-1k"] = struct{}{}
	if _, exist := excludeList[collectionSlug]; exist {
		return true
	}
	return false
}

func (u *Usecase) DEXSubmitOWPurchaseRawTx(address string, inscriptionID string, setupTx string, purchaseTx string) error {
	purchaseHash, setupHash, err := submitOWPurchaseRawTx(setupTx, purchaseTx)
	if err != nil {
		return err
	}

	if setupHash != "" {
		err = u.TrackWalletTx(address, structure.WalletTrackTx{Txhash: setupHash, Type: "buy-split-inscription", Amount: 0, InscriptionID: "", InscriptionNumber: 0, Receiver: address})
		if err != nil {
			log.Println("httpDelivery.trackTx.TrackWalletTx", err.Error(), err)
			return err
		}
	}

	err = u.TrackWalletTx(address, structure.WalletTrackTx{Txhash: purchaseHash, Type: "buy-inscription", Amount: 0, InscriptionID: inscriptionID, InscriptionNumber: 0, Receiver: address})
	if err != nil {
		log.Println("httpDelivery.trackTx.TrackWalletTx", err.Error(), err)
		return err
	}

	return nil
}

func crawlOWActivities(before time.Time) error {
	offset := 0
	stop := false
	newActivity := []structure.DexListingOWCollectionItem{}
	for {
		activities, err := getOWActivities(offset)
		if err != nil {
			return err
		}
		if len(activities) == 0 {
			break
		}
		for _, activity := range activities {
			if activity.Escrow.BoughtAt != "" {
				boughtAt, err := timefmt.Parse(activity.Escrow.BoughtAt, "%Y-%m-%d %H:%M:%S.%f%z")
				if err != nil {
					log.Fatal(err)
				}
				if boughtAt.Before(before) {
					stop = true
					break
				}
			}
			newActivity = append(newActivity, activity)
		}
		if stop {
			break
		}
		offset += len(activities)
	}

	log.Println("newActivity", len(newActivity))
	return nil
}

func getOWActivities(offset int) ([]structure.DexListingOWCollectionItem, error) {
	url := fmt.Sprintf("https://turbo.ordinalswallet.com/inscriptions/activity?offset=%v&order=BoughtAtDesc", offset)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var result []structure.DexListingOWCollectionItem
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

package crontab_ordinal_collections

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/robfig/cron.v2"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase"
	"rederinghub.io/utils/global"
	"rederinghub.io/utils/logger"
)

type ScronOrdinalCollectionHandler struct {
	Logger  logger.Ilogger
	Usecase usecase.Usecase
}

func NewScronOrdinalCollectionHandler(global *global.Global, uc usecase.Usecase) *ScronOrdinalCollectionHandler {
	return &ScronOrdinalCollectionHandler{
		Logger:  global.Logger,
		Usecase: uc,
	}
}

func (h ScronOrdinalCollectionHandler) syncCollection(collectionFoldersPath string, source string) error {
	collectionMetaFilePath := fmt.Sprintf("%s/meta.json", collectionFoldersPath)
	collectionInscriptionFilePath := fmt.Sprintf("%s/inscriptions.json", collectionFoldersPath)
	metaJsonFile, err := os.Open(collectionMetaFilePath)
	if err != nil {
		return err
	}
	defer metaJsonFile.Close()

	inscrJsonFile, err := os.Open(collectionInscriptionFilePath)
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer inscrJsonFile.Close()

	var meta entity.CollectionMeta
	var inscriptions []entity.CollectionInscription

	byteValue, _ := ioutil.ReadAll(metaJsonFile)
	json.Unmarshal(byteValue, &meta)
	byteValue, _ = ioutil.ReadAll(inscrJsonFile)
	json.Unmarshal(byteValue, &inscriptions)

	meta.Source = source

	_, err = h.Usecase.Repo.FindCollectionMetaByInscriptionIcon(meta.InscriptionIcon)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			h.Usecase.Repo.InsertCollectionMeta(&meta)
		} else {
			return err
		}
	} 

	insertedInscriptions, err := h.Usecase.Repo.FindCollectionInscriptionByInscriptionIcon(meta.InscriptionIcon)
	if err != nil {
		return err
	}
	insertedIds := map[string]bool{}
	for _, inscription := range insertedInscriptions {
		insertedIds[inscription.ID] = true
	}
	
	processed := 0
	for _, inscription := range inscriptions {
		if insertedIds[inscription.ID] {
			continue
		}
		processed++
		inscription.CollectionInscriptionIcon = meta.InscriptionIcon
		inscription.Source = source
		h.Usecase.Repo.InsertCollectionInscription(&inscription)

		if processed % 10 == 0 {
			time.Sleep(1 * time.Second)
		}
	}

	fmt.Printf("Done for collection %s \n", meta.Name)

	return nil
}

func (h ScronOrdinalCollectionHandler) crawlOrdinalCollection(source string) error {
	uuid := uuid.New().String()
	folder_path := fmt.Sprintf("/tmp/ordinals-collection-%s", uuid)

	_, err := git.PlainClone(folder_path, false, &git.CloneOptions{
		URL:      source,
		Progress: os.Stdout,
	})

	if err != nil {
		return err
	}

	collectionFoldersPath := fmt.Sprintf("%s/collections", folder_path)
	collectionFolders, err := ioutil.ReadDir(collectionFoldersPath)
	if err != nil {
		return err
	}
	for _, f := range collectionFolders {
		h.syncCollection(fmt.Sprintf("%s/%s", collectionFoldersPath, f.Name()), source)
	}
	return nil
}

func (h ScronOrdinalCollectionHandler) StartServer() {
	c := cron.New()
	// cronjob to sync ordinals collection
	c.AddFunc("0 */2 * * *", func() {
		source := "https://github.com/ordinals-wallet/ordinals-collections.git"
		err := h.crawlOrdinalCollection(source)
		if err != nil {
			h.Logger.Error("DispatchCron.EveryTwoHour.SyncOrdinalWalletCollections", err.Error(), err)
		}
	})
	c.AddFunc("*/15 * * * *", func() {
		source := "https://github.com/generative-xyz/ordinals-collections.git"
		err := h.crawlOrdinalCollection(source)
		if err != nil {
			h.Logger.Error("DispatchCron.EveryFifteenMinutes.SyncGenerativeCollections", err.Error(), err)
		}
	})
	c.Start()

	go func() {
		for {
			err := h.Usecase.CreateProjectsFromMetas()
			if err != nil {
				h.Logger.Error("error at cronjob create projects from metas", err.Error(), err)
				return
			}
			err = h.Usecase.CreateTokensFromCollectionInscriptions()
			if err != nil {
				h.Logger.Error("error at cronjob create tokens from collection inscription", err.Error(), err)
				return
			}
			// Sleep 5 minutes after recreate again
			time.Sleep(1 * time.Minute)
		}
	}()

}

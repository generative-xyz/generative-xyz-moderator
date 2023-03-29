package usecase

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/repository"
)

func (u Usecase) FindProjectByInscriptionIcon(inscriptionIcon string) (*entity.Projects, error) {
	project, err := u.Repo.FindProjectByInscriptionIcon(inscriptionIcon)
	if err != nil {
		if !errors.Is(err, mongo.ErrNoDocuments) {
			return nil, err
		}
	} else {
		return project, nil
	}

	// find project from inscription icon
	inscriptions, err := u.Repo.FindCollectionInscriptionByInscriptionIcon(inscriptionIcon)
	if err != nil {
		return nil, err
	}
	tokenIds := []string{}
	for _, inscription := range inscriptions {
		tokenIds = append(tokenIds, inscription.ID)
	}

	token, err := u.Repo.FindOneTokenByListOfTokenIds(tokenIds)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	project, err = u.Repo.FindProjectByTokenID(token.ProjectID)
	if err != nil {
		return nil, err
	}

	if project.InscriptionIcon == "" {
		err := u.Repo.SetProjectInscriptionIcon(project.TokenID, inscriptionIcon)
		if err != nil {
			return nil, err
		}
	}

	return project, nil
}

func (u Usecase) CreateProjectsFromMetas() error {
	uncreatedMetas, err := u.Repo.FindUncreatedCollectionMeta()
	if err != nil {
		return err
	}
	u.Logger.Info("Start create projects from meta len(uncreatedMetas)=%v", len(uncreatedMetas))
	processed := 0
	for _, meta := range uncreatedMetas {
		u.Logger.Info(fmt.Sprintf("Trying to create project from collection meta %s %s", meta.Name, meta.InscriptionIcon))
		processed++
		project, err := u.FindProjectByInscriptionIcon(meta.InscriptionIcon)
		if err != nil {
			u.Logger.ErrorAny("u.FindProjectByInscriptionIcon", zap.Error(err))
			continue
		}

		existed := true

		if project == nil {
			project, err = u.CreateProjectFromCollectionMeta(meta)
			if err != nil {
				u.Logger.Error("u.CreateProjectFromCollectionMeta", err.Error(), err)
				return err
			}
			u.Logger.LogAny(fmt.Sprintf("Created project from collection meta %s %s", meta.Name, meta.InscriptionIcon))
			existed = false
		}

		err = u.Repo.SetProjectCreatedMeta(meta)
		if err != nil {
			u.Logger.Error("u.Repo.SetProjectCreatedMeta", err.Error(), err)
			continue
		}

		if project == nil {
			u.Logger.ErrorAny("CreateProjectsFromMetas.NilProject")
			continue
		}
		
		err = u.Repo.SetMetaMappedProjectID(meta, project.TokenID)
		if err != nil {
			u.Logger.Error("CreateProjectsFromMetas.SetMetaMappedProjectID", zap.Error(err))
			continue
		}
		u.Logger.LogAny("SetMetaMappedProjectID", zap.Any("projectID", project.TokenID))

		err = u.Repo.SetMetaProjectExisted(meta, existed)
		if err != nil {
			u.Logger.Error("CreateProjectsFromMetas.SetMetaProjectExisted", zap.Error(err))
			continue
		}
		u.Logger.LogAny("SetMetaProjectExisted", zap.Any("projectID", project.TokenID))
		
		
		if processed % 20 == 0 {
			time.Sleep(1 * time.Second)
		}
	}

	return nil
}


func (u Usecase) CreateTokensFromCollectionInscriptions() error {
	uncreatedInscription, err := u.Repo.FindUncreatedCollectionInscription()
	if err != nil {
		return err
	}
	u.Logger.Info("c len(uncreatedMetas)=%v", len(uncreatedInscription))
	processed := 0
	for _, inscription := range uncreatedInscription {
		u.Logger.Info(fmt.Sprintf("Trying to create token from collection inscription %s %s", inscription.Meta.Name, inscription.ID))
		processed++

		_, err = u.Repo.FindTokenByTokenID(inscription.ID)
		if err != nil {
			if !errors.Is(err, mongo.ErrNoDocuments) {
				u.Logger.ErrorAny("u.Repo.FindTokenByTokenID " + inscription.ID, zap.Error(err))
				continue
			} else {
				meta, err := u.Repo.FindCollectionMetaByInscriptionIcon(inscription.CollectionInscriptionIcon)
				if err != nil {
					u.Logger.ErrorAny("u.Repo.FindCollectionMetaByInscriptionIcon", zap.Error(err))
					continue
				}

				if meta.ProjectExisted {
					u.Logger.LogAny("MetaProjectIsAlreadyExisted", zap.Any("meta", meta))
					continue
				}

				_, err = u.CreateBTCTokenURIFromCollectionInscription(*meta, inscription)
				if err != nil {
					if !errors.Is(err, repository.ErrNoProjectsFound) {
						u.Logger.ErrorAny("u.CreateBTCTokenURIFromCollectionInscription", zap.Error(err))
						continue
					}
				} else {
					u.Logger.Info(fmt.Sprintf("Done create token %s", inscription.ID))
				}
			}
		}

		err = u.Repo.SetTokenCreatedInscription(inscription)
		u.Logger.Info(fmt.Sprintf("Done set token created %s", inscription.ID))

		if err != nil {
			u.Logger.ErrorAny("u.CreateBTCTokenURIFromCollectionInscription", zap.Error(err))
			continue
		}
		
		if processed % 20 == 0 {
			time.Sleep(1 * time.Second)
		}
	}

	return nil
}

func (h Usecase) syncCollection(collectionFoldersPath string, source string) error {
	h.Logger.LogAny("syncCollection.start", zap.String("collectionFoldersPath", collectionFoldersPath), zap.String("source", source))
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
	meta.WalletAddress = strings.ToLower(meta.WalletAddress)

	_, err = h.Repo.FindCollectionMetaByInscriptionIcon(meta.InscriptionIcon)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			h.Repo.InsertCollectionMeta(&meta)
		} else {
			return err
		}
	}

	insertedInscriptions, err := h.Repo.FindCollectionInscriptionByInscriptionIcon(meta.InscriptionIcon)
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
		h.Repo.InsertCollectionInscription(&inscription)

		if processed%10 == 0 {
			time.Sleep(1 * time.Second)
		}
	}

	fmt.Printf("Done for collection %s \n", meta.Name)

	return nil
}

func (h Usecase) crawlOrdinalCollection(source string) error {
	h.Logger.LogAny("crawlOrdinalCollection.start", zap.String("source", source))
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

func (u Usecase) JobCrawlGenerativeOrdinalCollection() error {
	source := "https://github.com/generative-xyz/ordinals-collections.git"
	err := u.crawlOrdinalCollection(source)
	if err != nil {
		u.Logger.ErrorAny(
			"JobCrawlGenerativeOrdinalCollection.crawlOrdinalCollection",
			zap.Error(err),
		)
	}
	return nil
}

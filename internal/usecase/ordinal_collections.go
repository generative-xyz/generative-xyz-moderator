package usecase

import (
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"rederinghub.io/internal/entity"
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
			u.Logger.Error("u.FindProjectByInscriptionIcon", err.Error(), err)
			continue
		}

		if project == nil {
			_, err = u.CreateProjectFromCollectionMeta(meta)
			if err != nil {
				u.Logger.Error("u.CreateProjectFromCollectionMeta", err.Error(), err)
				return err
			}
			u.Logger.Info(fmt.Sprintf("Created project from collection meta %s %s", meta.Name, meta.InscriptionIcon))
		}

		err = u.Repo.SetProjectCreatedMeta(meta)
		if err != nil {
			u.Logger.Error("u.Repo.SetProjectCreatedMeta", err.Error(), err)
			continue
		}
		
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
				_, err = u.CreateBTCTokenURIFromCollectionInscription(*meta, inscription)
				if err != nil {
					u.Logger.ErrorAny("u.CreateBTCTokenURIFromCollectionInscription", zap.Error(err))
					continue
				}
				u.Logger.Info(fmt.Sprintf("Done create token %s", inscription.ID))
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

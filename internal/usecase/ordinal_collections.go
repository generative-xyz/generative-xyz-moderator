package usecase

import (
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/repository"
	"rederinghub.io/utils/logger"
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
	logger.AtLog.Logger.Info("Start create projects from meta len(uncreatedMetas)=%v", zap.Any("len(uncreatedMetas)", len(uncreatedMetas)))
	processed := 0
	for _, meta := range uncreatedMetas {
		logger.AtLog.Logger.Info(fmt.Sprintf("Trying to create project from collection meta %s %s", meta.Name, meta.InscriptionIcon))
		processed++
		project, err := u.FindProjectByInscriptionIcon(meta.InscriptionIcon)
		if err != nil {
			logger.AtLog.Logger.Error("u.FindProjectByInscriptionIcon", zap.Error(err))
			continue
		}

		existed := true

		if project == nil {
			project, err = u.CreateProjectFromCollectionMeta(meta)
			if err != nil {
				logger.AtLog.Logger.Error("u.CreateProjectFromCollectionMeta", zap.Error(err))
				return err
			}
			u.Logger.LogAny(fmt.Sprintf("Created project from collection meta %s %s", meta.Name, meta.InscriptionIcon))
			existed = false
		}

		err = u.Repo.SetProjectCreatedMeta(meta)
		if err != nil {
			logger.AtLog.Logger.Error("u.Repo.SetProjectCreatedMeta", zap.Error(err))
			continue
		}

		if project == nil {
			logger.AtLog.Logger.Error("CreateProjectsFromMetas.NilProject")
			continue
		}
		
		err = u.Repo.SetMetaMappedProjectID(meta, project.TokenID)
		if err != nil {
			logger.AtLog.Logger.Error("CreateProjectsFromMetas.SetMetaMappedProjectID", zap.Error(err))
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
	logger.AtLog.Logger.Info("c len(uncreatedMetas)=%v", zap.Any("len(uncreatedInscription)", len(uncreatedInscription)))
	processed := 0
	for _, inscription := range uncreatedInscription {
		logger.AtLog.Logger.Info(fmt.Sprintf("Trying to create token from collection inscription %s %s", inscription.Meta.Name, inscription.ID))
		processed++

		_, err = u.Repo.FindTokenByTokenID(inscription.ID)
		if err != nil {
			if !errors.Is(err, mongo.ErrNoDocuments) {
				logger.AtLog.Logger.Error("u.Repo.FindTokenByTokenID " + inscription.ID, zap.Error(err))
				continue
			} else {
				meta, err := u.Repo.FindCollectionMetaByInscriptionIcon(inscription.CollectionInscriptionIcon)
				if err != nil {
					logger.AtLog.Logger.Error("u.Repo.FindCollectionMetaByInscriptionIcon", zap.Error(err))
					continue
				}

				if meta.ProjectExisted {
					u.Logger.LogAny("MetaProjectIsAlreadyExisted", zap.Any("meta", meta))
					continue
				}

				_, err = u.CreateBTCTokenURIFromCollectionInscription(*meta, inscription)
				if err != nil {
					if !errors.Is(err, repository.ErrNoProjectsFound) {
						logger.AtLog.Logger.Error("u.CreateBTCTokenURIFromCollectionInscription", zap.Error(err))
						continue
					}
				} else {
					logger.AtLog.Logger.Info(fmt.Sprintf("Done create token %s", inscription.ID))
				}
			}
		}

		err = u.Repo.SetTokenCreatedInscription(inscription)
		logger.AtLog.Logger.Info(fmt.Sprintf("Done set token created %s", inscription.ID))

		if err != nil {
			logger.AtLog.Logger.Error("u.CreateBTCTokenURIFromCollectionInscription", zap.Error(err))
			continue
		}
		
		if processed % 20 == 0 {
			time.Sleep(1 * time.Second)
		}
	}

	return nil
}

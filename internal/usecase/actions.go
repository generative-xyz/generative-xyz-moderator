package usecase

import (
	"rederinghub.io/internal/entity"
	"strings"
)

func (u *Usecase) LikeProject(projectID string, walletAddress string) (*entity.Actions, error) {
	obj := &entity.Actions{
		CreatedBy:  strings.ToLower(walletAddress),
		ObjectID:   strings.ToLower(projectID),
		ObjectType: entity.Project,
		Action:     entity.LIKE,
	}

	err := u.Repo.InsertAction(obj)
	if err != nil {
		return nil, err
	}

	return obj, nil
}

func (u *Usecase) DisLikeProject(projectID string, walletAddress string) (*entity.Actions, error) {
	obj := &entity.Actions{
		CreatedBy:  strings.ToLower(walletAddress),
		ObjectID:   strings.ToLower(projectID),
		ObjectType: entity.Project,
		Action:     entity.DISLIKE,
	}

	err := u.Repo.InsertAction(obj)
	if err != nil {
		return nil, err
	}

	return obj, nil
}

func (u *Usecase) LikeToken(tokenID string, walletAddress string) (*entity.Actions, error) {
	obj := &entity.Actions{
		CreatedBy:  strings.ToLower(walletAddress),
		ObjectID:   strings.ToLower(tokenID),
		ObjectType: entity.TokenURI,
		Action:     entity.LIKE,
	}

	err := u.Repo.InsertAction(obj)
	if err != nil {
		return nil, err
	}

	return obj, nil
}

func (u *Usecase) DisLikeToken(tokenID string, walletAddress string) (*entity.Actions, error) {
	obj := &entity.Actions{
		CreatedBy:  strings.ToLower(walletAddress),
		ObjectID:   strings.ToLower(tokenID),
		ObjectType: entity.TokenURI,
		Action:     entity.DISLIKE,
	}

	err := u.Repo.InsertAction(obj)
	if err != nil {
		return nil, err
	}

	return obj, nil
}

package usecase

import (
	"rederinghub.io/internal/entity"
	"strings"
)

func (u *Usecase) LikeProject(contractAddress string, projectID string) (*entity.Actions, error) {
	obj := &entity.Actions{
		Parent:     strings.ToLower(contractAddress),
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

func (u *Usecase) DisLikeProject(contractAddress string, projectID string) (*entity.Actions, error) {
	obj := &entity.Actions{
		Parent:     strings.ToLower(contractAddress),
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

func (u *Usecase) LikeToken(projectID string, tokenID string) (*entity.Actions, error) {
	obj := &entity.Actions{
		Parent:     strings.ToLower(projectID),
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

func (u *Usecase) DisLikeToken(projectID string, tokenID string) (*entity.Actions, error) {
	obj := &entity.Actions{
		Parent:     strings.ToLower(projectID),
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

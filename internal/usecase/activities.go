package usecase

import (
	"fmt"
	"strconv"

	"rederinghub.io/internal/entity"
)

func (u Usecase) CreateViewProjectActivity(projectID string) {
	err := u.Repo.InsertActitvy(&entity.Activity{
		Type: entity.View,
		ProjectID: projectID,
	})

	if err != nil {
		fmt.Printf("CreateViewProjectActivity.%s.Error:%s", projectID, err.Error())
	}
}

func (u Usecase) CreateMintActivity(inscriptionID string, value string) {
	iValue, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
			iValue = 0
	}

	tokenUri, err := u.Repo.FindTokenByTokenID(inscriptionID)
	if err != nil {
		fmt.Printf("CreateMintActivity.FindTokenByTokenID.%s.Error:%s", inscriptionID, err.Error())
		return
	}

	err = u.Repo.InsertActitvy(&entity.Activity{
		Type: entity.Mint,
		Reference: inscriptionID,
		Value: iValue,
		ProjectID: tokenUri.ProjectID,
	})

	if err != nil {
		fmt.Printf("CreateMintActivity.%s.Error:%s", inscriptionID, err.Error())
		return
	}
}

func (u Usecase) CreateBuyActivity(inscriptionID string, value string) {
	iValue, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
			iValue = 0
	}

	tokenUri, err := u.Repo.FindTokenByTokenID(inscriptionID)
	if err != nil {
		fmt.Printf("CreateBuyActivity.FindTokenByTokenID.%s.Error:%s", inscriptionID, err.Error())
		return
	}

	err = u.Repo.InsertActitvy(&entity.Activity{
		Type: entity.Buy,
		Reference: inscriptionID,
		Value: iValue,
		ProjectID: tokenUri.ProjectID,
	})

	if err != nil {
		fmt.Printf("CreateBuyActivity.%s.Error:%s", inscriptionID, err.Error())
		return
	}
}

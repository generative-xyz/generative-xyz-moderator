package structure

import (
	"errors"
	"fmt"
	"strconv"

	"rederinghub.io/utils/helpers"
)
type OrdinalCollectionMeta struct {
	Name string	`json:"name"`
	InscriptionIcon string `json:"inscription_icon"`
	Supply string `json:"supply"`
	Slug string `json:"slug"`
	Description string `json:"description"`
	TwitterLink string `json:"twitter_link"`
	DiscordLink string `json:"discord_link"`
	WebsiteLink string `json:"website_link"`
	WalletAddress string `json:"wallet_address"`
	Royalty string `json:"royalty"`
}

func (v *OrdinalCollectionMeta) Verify() error {
	if v.Name == "" {
		err := errors.New("Collection is required")
		return err
	}
	
	if v.InscriptionIcon == "" {
		err := errors.New("InscriptionIcon is required")
		return err
	}
	
	if v.Supply == "" {
		err := errors.New("Supply is required")
		return err
	}
	
	supplyInt, err := strconv.Atoi(v.Supply) 
	if err != nil {
		err = fmt.Errorf("Supply must be a number %v", err)
		return err
	}

	if supplyInt <= 0 {
		err := errors.New("Supply must be greater than 0")
		return err
	}

	if v.Royalty == "" {
		err := errors.New("Royalty is required")
		return err
	}
	
	royaltyInt, err := strconv.Atoi(v.Royalty) 
	if err != nil {
		err = fmt.Errorf("Royalty must be a number %v", err)
		return err
	}

	if royaltyInt < 0 {
		err := errors.New("Royalty must be greater than 0")
		return err
	}


	if v.Slug == "" {
		v.Slug = helpers.GenerateSlug(v.Name)
		return err
	}

	return nil
}

type OrdinalInscriptionMeta struct {
	ID string	`json:"id"`
	Meta map[string]string `json:"meta"`
}


func (v OrdinalInscriptionMeta) Verify() error {
	if v.ID == "" {
		err := errors.New("Inscription ID is required")
		return err
	}
	
	if len(v.Meta) == 0 {
		err := errors.New("Meta is required")
		return err
	}

	name, ok := v.Meta["name"]
	if !ok {
		err := errors.New("Meta.Name is required")
		return err
	}
	
	if name=="" {
		err := errors.New("Meta.Name is not empty")
		return err
	}

	return nil
}

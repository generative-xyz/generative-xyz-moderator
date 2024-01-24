package usecase

import (
	"context"
	"fmt"
	"os"
	"rederinghub.io/internal/usecase/structure"
	"strings"
)

type InsOwner struct {
	InscriptionID string
	OwnerAddress  string
	IsUpdated     bool
	Err           error
}

type Ins struct {
	InscriptionID string
	OwnerAddress  string
}

func (u Usecase) CrontabUpdateModularInscOwners(ctx context.Context) error {

	page := 1
	limit := 100
	genNFTAddr := os.Getenv("MODULAR_PROJECT_ID")

	for {
		offset := (page - 1) * limit
		inscriptions, err := u.Repo.AggregateModularInscriptions(ctx, genNFTAddr, offset, limit)
		if err != nil {
			return err
		}

		if len(inscriptions) == 0 {
			break
		}

		inChan := make(chan Ins, len(inscriptions))
		outChan := make(chan InsOwner, len(inscriptions))

		for range inscriptions {
			go u.FindModularInscOwner(inChan, outChan)
		}

		for _, i := range inscriptions {
			inChan <- Ins{
				InscriptionID: i.TokenID,
				OwnerAddress:  i.OwnerAddr,
			}
		}

		for range inscriptions {
			outFChan := <-outChan
			if outFChan.Err != nil {
				continue
			}

			if outFChan.IsUpdated {
				//TODO - update owner
				fmt.Println(fmt.Sprintf("[ins] %s-%s-%v", outFChan.InscriptionID, outFChan.OwnerAddress, outFChan.IsUpdated))
			}
		}

		page++
	}

	return nil
}

func (u Usecase) FindModularInscOwner(in chan Ins, out chan InsOwner) {
	var err error
	addr := ""
	inscID := <-in
	info := &structure.InscriptionOrdInfoByID{}
	isUpdate := false

	defer func() {
		isUpdate = !strings.EqualFold(inscID.OwnerAddress, addr)
		out <- InsOwner{
			Err:           err,
			InscriptionID: inscID.InscriptionID,
			OwnerAddress:  addr,
			IsUpdated:     isUpdate,
		}
	}()

	info, err = u.GetInscriptionByIDFromOrd(inscID.InscriptionID)
	if err != nil {
		return
	}

	addr = info.Address
}

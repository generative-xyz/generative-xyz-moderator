package usecase

import "rederinghub.io/internal/usecase/structure"

func (u Usecase) SyncTraitStats() error {
	allTokens, err := u.Repo.GetAllTokensHasTraitSeletedFields()
	if err != nil {
		return err
	}
	projectSet := map[string]bool{}
	for _, token := range allTokens {
		projectSet[token.ProjectID] = true
	}
	for projectID := range projectSet {
		_, traitStats, err := u.GetUpdatedProjectStats(structure.GetProjectReq{TokenID: projectID, ContractAddr: "0x0000000000000000000000000000000000000000"})
		if err != nil {
			return err
		}
		u.Repo.UpdateProjectTraitStats(projectID, traitStats)
	}
	return nil
}

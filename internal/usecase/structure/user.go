package structure

import "rederinghub.io/internal/entity"

type FilterUsers struct {
	BaseFilters
	Search           *string
	Email            *string
	WalletAddress    *string
	WalletAddressBTC *string
	UserType         *entity.UserType
	IsUpdatedAvatar  *bool
}

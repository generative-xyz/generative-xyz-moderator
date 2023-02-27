package structure

import "rederinghub.io/internal/entity"

type FilterReferrals struct {
	BaseFilters
	ReferrerID *string
	ReferreeID *string
	PayType *string
}

type ReferalResp struct {
	ReferrerID string  
	ReferreeID string 
	Referrer   *entity.Users 
	Referree   *entity.Users  
	Percent    int32 	
	ReferreeVolume    ReferralVolumnResp	
}

type ReferralVolumnResp struct {
	Amount string 
	AmountType string 
	ProjectID string 
	Percent int 
	Earn string 
	GenEarn string 
}

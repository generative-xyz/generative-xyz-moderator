package structure

import "time"

type DexVolumeInscritionFilter struct {
}

type AggerateChartForProject struct {
	ProjectID *string `json:"projectID"`
	FromDate *time.Time `json:"fromDate"`
	ToDate *time.Time `json:"toDate"`
}

type AggerateChartForToken struct {
	TokenID *string `json:"tokenID"`
	FromDate *time.Time `json:"fromDate"`
	ToDate *time.Time `json:"toDate"`
}

type AggragetedCollection struct {
	ProjectID string `json:"projectID"` 
	ProjectName string `json:"projectName"` 
	Timestamp string `json:"timestamp"` 
	Amount int64 `json:"amount"`
}

type AggragetedTokenURI struct {
	TokenID string `json:"tokenID"` 
	Timestamp string `json:"timestamp"` 
	Amount int64 `json:"amount"`
}

type AggragetedCollectionVolumnResp struct {
	Volumns []AggragetedCollection `json:"volumns"` 
}

type AggragetedTokenVolumnResp struct {
	Volumns []AggragetedTokenURI `json:"volumns"` 
}
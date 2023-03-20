package structure

import "time"

type DexVolumeInscritionFilter struct {
}

type AggerateChartForProject struct {
	ProjectID *string `json:"projectID"`
	FromDate *time.Time `json:"fromDate"`
	ToDate *time.Time `json:"toDate"`
}

type AggragetedInscription struct {
	ProjectID string `json:"projectID"` 
	ProjectName string `json:"projectName"` 
	Timestamp string `json:"timestamp"` 
	Amount int64 `json:"amount"`
}

type AggragetedInscriptionVolumnResp struct {
	Volumns []AggragetedInscription `json:"volumns"` 
}
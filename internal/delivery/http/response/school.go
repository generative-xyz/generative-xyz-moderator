package response

type AISchoolJobProgress struct {
	Progress    int    `json:"progress"`
	Status      string `json:"status"`
	JobID       string `json:"job_id"`
	Output      string `json:"output"`
	CompletedAt int64  `json:"completed_at"`
	CreatedAt   int64  `json:"created_at"`
	ModelName   string `json:"model_name"`
}

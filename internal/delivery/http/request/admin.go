package request

type UpsertRedisRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

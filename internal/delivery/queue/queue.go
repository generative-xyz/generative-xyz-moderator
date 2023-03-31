package queue

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/hibiken/asynq"
	"go.uber.org/zap"
	"rederinghub.io/internal/usecase"
	"rederinghub.io/utils"
	"rederinghub.io/utils/config"
	"rederinghub.io/utils/logger"
)

type QueueHandler struct {
	usecase usecase.Usecase
	redisConf     config.RedisConfig
	log logger.Ilogger
}

func NewQueueHandler(usecase usecase.Usecase, redisConf config.RedisConfig, log logger.Ilogger) *QueueHandler {
	return &QueueHandler{
		usecase: usecase,
		redisConf:     redisConf,
		log: log,
	}
}

func (h QueueHandler) StartServer() {

	redisDBInt, err := strconv.Atoi(h.redisConf.DB)
	if err != nil {
		redisDBInt = 1
	}

	srv := asynq.NewServer(
        asynq.RedisClientOpt{
			Addr: h.redisConf.Address,
			Password: h.redisConf.Password,
			DB: redisDBInt,
		},
        asynq.Config{
            // Specify how many concurrent workers to use
            Concurrency: Concurrency,
            // Optionally specify multiple queues with different priority.
            Queues: map[string]int{
                "critical": Critical,
                "default":  Default,
                "low":      Low,
            },
            // See the godoc for other configuration options
        },
    )

	mux := asynq.HandlerFunc(h.handler)
	err = srv.Run(mux)
    if err != nil {
        log.Fatalf("could not run server: %v", err)
    }
}

// ImageProcessor implements asynq.Handler interface.
func (h *QueueHandler) handler(ctx context.Context, t *asynq.Task) error {
	logger.AtLog.Logger.Info("recived task", zap.Any("payload", t.Payload()), zap.Any("taskType", t.Type))
	unizipType := fmt.Sprintf("%s:%s", h.redisConf.ENV, utils.PUBSUB_PROJECT_UNZIP)
	tokenThumbnail := fmt.Sprintf("%s:%s", h.redisConf.ENV, utils.PUBSUB_TOKEN_THUMBNAIL)

	//spew.Dump(t.Type())
	switch t.Type() {
		case unizipType:
			return h.usecase.ProccessUnzip(ctx, t)
		case tokenThumbnail:
			return h.usecase.ProccessCreateTokenThumbnail(ctx, t)
	}

	return nil
}

package worker

import (
	db "github.com/kitanoyoru/kita-ci/internal/db/postgres"
	"github.com/kitanoyoru/kita-ci/internal/queue"
	"github.com/kitanoyoru/kita-ci/pkg/config"
	"github.com/kitanoyoru/kita-ci/pkg/log"
)

type CIWorker struct {
	config    *config.WorkerConfig
	jobsQueue queue.Queue
	dbClient  *db.PostgresClient // TODO: Change to interface
	logger    log.ILogger
}

func NewCIWorker(cfg *config.WorkerConfig, logger log.ILogger) *CIWorker {
	dbConfig := config.PostgresConfig{
		DB:         cfg.DB,
		DBUser:     cfg.DBUser,
		DBPassword: cfg.DBPassword,
		DBAddr:     cfg.DBAddr,
	}
	dbClient := db.NewPostgresClient(dbConfig)

	ciJobsQueue := queue.NewRMQQueue(cfg.RabbitMQAddr)

	worker := &CIWorker{
		config:    cfg,
		jobsQueue: ciJobsQueue,
		dbClient:  dbClient,
		logger:    logger,
	}

	return worker
}

func (w *CIWorker) Run() {
	w.startConsuming()
}

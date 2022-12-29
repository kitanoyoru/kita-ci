package worker

import (
	"github.com/kitanoyoru/kita-ci/internal/db"
	"github.com/kitanoyoru/kita-ci/internal/docker"
	"github.com/kitanoyoru/kita-ci/internal/queue"
	"github.com/kitanoyoru/kita-ci/pkg/config"
	"github.com/kitanoyoru/kita-ci/pkg/log"
)

type CIWorker struct {
	config       *config.WorkerConfig
	jobsQueue    queue.Queue
	dbClient     *db.PostgresClient // TODO: Change to interface
	dockerClient *docker.DockerClient
	logger       log.ILogger
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

	dockerConfig := config.DockerConfig{
		Image: cfg.ImageBuilder,
	}
	dockerClient := docker.NewDockerClient(dockerConfig)

	worker := &CIWorker{
		config:       cfg,
		jobsQueue:    ciJobsQueue,
		dbClient:     dbClient,
		dockerClient: dockerClient,
		logger:       logger,
	}

	return worker
}

func (w *CIWorker) Run() {
	w.StartConsuming()
}

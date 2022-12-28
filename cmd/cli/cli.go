package cli

import (
	"github.com/spf13/cobra"

	"github.com/kitanoyoru/kita-ci/internal/worker"
	"github.com/kitanoyoru/kita-ci/pkg/config"
	"github.com/kitanoyoru/kita-ci/pkg/log"
)

var (
	port         int
	rabbitMQAddr string
	db           string
	dbUser       string
	dbPassword   string
	dbAddr       string
	imageBuilder string
	logLevel     string
)

var rootCmd = &cobra.Command{
	Use:   "kita-ci",
	Short: "Kita CI worker microservice",
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start kita CI worker",
	Run: func(cmd *cobra.Command, args []string) {
		config := &config.WorkerConfig{
			Port:         port,
			RabbitMQAddr: rabbitMQAddr,
			DB:           db,
			DBUser:       dbUser,
			DBPassword:   dbPassword,
			DBAddr:       dbAddr,
			ImageBuilder: imageBuilder,
		}

		logger := log.NewLogger(logLevel)

		worker := worker.NewCIWorker(config, logger)

		worker.Run()
	},
}

func Run() {
	if err := rootCmd.Execute(); err != nil {
		panic(err) // TODO: Make it more cleaner
	}
}

func init() {
	rootCmd.AddCommand(startCmd)

	startCmd.Flags().IntVarP(&port, "port", "p", 2022, "CI worker port")
	startCmd.Flags().StringVarP(&rabbitMQAddr, "rabbitmq", "rmq", "amqp://guest:guest@localhost:5672", "RabbitMQ address")
	startCmd.Flags().StringVar(&db, "db", "kita", "PostgreSQL databse")
	startCmd.Flags().StringVar(&dbUser, "db-user", "kita", "PostgreSQL user")
	startCmd.Flags().StringVar(&dbPassword, "db-pass", "kita", "PostgreSQL password")
	startCmd.Flags().StringVar(&dbAddr, "db-addr", "postgres:5432", "PostgreSQL address")
	startCmd.Flags().StringVar(&imageBuilder, "builder", "image-builder", "Docker image builder name")
	startCmd.Flags().StringVar(&logLevel, "log-level", "INFO", "Log level")
}

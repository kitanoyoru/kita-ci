package cli

import (
	"log"

	"github.com/spf13/cobra"
)

var (
	port int

	rabbitMQAddr string

	db         string
	dbUser     string
	dbPassword string
	dbAddr     string

	imageBuilder string

	logLevel string
)

var rootCmd = &cobra.Command{
	User:  "kita-ci",
	Short: "Kita CI worker microservice",
}

var startCmd = &cobra.Command{
	User:  "start",
	Short: "Start kita CI worker",
	Run: func(cmd *cobra.Command, args []string) {
		config := config.NewWorkerConfig{
			Port:         port,
			RabbitMQAddr: rabbitMQAddr,
			DB:           db,
			DBUser:       dbUser,
			DBPassword:   dbPassword,
			DBAddr:       dbAddr,
			ImageBuilder: imageBuilder,
		}

		logger = log.NewLogger(logLevel)

		worker := worker.NewWorker(config, logLevel)

		worker.Run()
	},
}

func Run() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
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

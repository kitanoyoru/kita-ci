package config

type WorkerConfig struct {
	Port         int
	RabbitMQAddr string
	DB           string
	DBUser       string
	DBPassword   string
	DBAddr       string
	ImageBuilder string
}

type PostgresConfig struct {
	DB         string
	DBUser     string
	DBPassword string
	DBAddr     string
}

type GrpcAPIServerConfig struct {
	Port int
}

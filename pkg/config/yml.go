package config

import "github.com/spf13/viper"

const (
	ConfigDir = "./config"
)

func ReadConfigFile() (*WorkerConfig, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yml")

	viper.AddConfigPath(ConfigDir)

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	workerConfig := &WorkerConfig{
		Port:         viper.GetInt("kita-ci.port"),
		RabbitMQAddr: viper.GetString("kita-ci.rabbitmq"),
		DB:           viper.GetString("kita-ci.db"),
		DBUser:       viper.GetString("kita-ci.dbUser"),
		DBPassword:   viper.GetString("kita-ci.dbPassword"),
		ImageBuilder: viper.GetString("kita-ci.imageBuilder"),
	}

	return workerConfig, nil
}

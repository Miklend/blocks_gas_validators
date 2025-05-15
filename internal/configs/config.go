package configs

import (
	"blocks_gas_validators/pkg/logging"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	IsDebug *bool         `yaml:"is_debug"`
	Listen  ListenConfig  `yaml:"listen"`
	Storage StorageConfig `yaml:"storage"`
	Alchemy AlchemyConfig `yaml:"alchemy"`
}

type ListenConfig struct {
	Type   string `yaml:"type"`
	BindIP string `yaml:"bind_ip"`
	Port   string `yaml:"port"`
}
type StorageConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Database string `yaml:"database"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type AlchemyConfig struct {
	Mode        string `yaml:"mode"`
	NetworkName string `yaml:"network_name"`
	NameApiKey  string `yaml:"name_api_key"`
	Limiter     int    `yaml:"limiter"`
	MaxRetries  int    `yaml:"max_retries"`
	BatchSize   int    `yaml:"batch_size"`
	Start       uint64 `yaml:"start"`
	End         uint64 `yaml:"end"`
	Workers     int    `yaml:"workers"`
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		logger := logging.GetLogger()
		logger.Info("read application configuration")
		instance = &Config{}
		if err := cleanenv.ReadConfig("config.yaml", instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			logger.Info(help)
			logger.Fatal(err)
		}
	})
	return instance
}

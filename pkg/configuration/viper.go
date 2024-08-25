package configuration

import (
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type (
	// Config
	Config struct {
	}
)

func getConfigName() string {
	mode := "dev"
	switch os.Getenv("MODE") {
	case "uat":
		mode = "uat"
	case "prod":
		mode = "prod"
	}
	return mode
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}
	vn := viper.New()

	// dường dẫn file config
	path := "config"
	configName := getConfigName()
	vn.AddConfigPath(path)
	vn.SetConfigName(configName)
	// ví du: config/dev.yaml

	vn.SetConfigType("yaml")
	vn.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	vn.AutomaticEnv()

	err := vn.ReadInConfig()
	if err != nil {
		return cfg, err

	}

	for _, key := range vn.AllKeys() {
		str := strings.ToUpper(strings.ReplaceAll(key, ".", "_"))
		if configName != "prod" {
			// đùng để print value ra trên môi trường non-prod
			log.Default().Println(key, str, vn.Get(key))
		}
		vn.BindEnv(key, str)
	}

	err = vn.Unmarshal(&cfg)
	if err != nil {
		return cfg, err
	}

	return cfg, err
}

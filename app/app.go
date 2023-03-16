package app

import (
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"os"
	"path"
	"strings"
)

func InitApp(rootDir ...string) error {

	basePath := "."
	if len(rootDir) != 0 {
		basePath = rootDir[0]
	}

	envPath := path.Join(basePath, ".env")
	if _, err := os.Stat(envPath); err == nil {
		err := godotenv.Load(envPath)
		if err != nil {
			return errors.Wrap(err, "godotenv cannot load config")
		}
	}

	yamlPath := path.Join(basePath, "config.yaml")
	if _, err := os.Stat(yamlPath); errors.Is(err, os.ErrNotExist) {
		return errors.Errorf("create config.yaml first in root directory")
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(basePath)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	err := viper.ReadInConfig()
	if err != nil {
		return errors.Wrap(err, "viper cannot read config")
	}

	return nil
}

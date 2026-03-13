package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-viper/mapstructure/v2"
	"github.com/spf13/viper"
)

type Environments struct {
	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     int    `mapstructure:"DB_PORT"`
	DBUsername string `mapstructure:"DB_USERNAME"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBName     string `mapstructure:"DB_NAME"`
}

func LoadEnvs() *Environments {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	
	env := &Environments{}
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	if err := loadMappedEnvVars(env); err != nil {
		panic(err)
	}

	if err := viper.Unmarshal(&env); err != nil {
		panic(err)
	}

	setEnvVars(viper.AllSettings())

	return env
}

func setEnvVars(s map[string]any) {
	for k, v := range s {
		if os.Getenv(strings.ToUpper(k)) == "" {
			os.Setenv(strings.ToUpper(k), fmt.Sprintf("%v", v))
		}
	}
}

func loadMappedEnvVars(env *Environments) error {
	envKeysMap := &map[string]any{}
	if err := mapstructure.Decode(env, &envKeysMap); err != nil {
		return err
	}

	for k := range *envKeysMap {
		if err := viper.BindEnv(k); err != nil {
			return err
		}
	}

	return nil
}

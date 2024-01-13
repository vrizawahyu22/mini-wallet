package utils

import "github.com/spf13/viper"

type BaseConfig struct {
	DBConnString string `mapstructure:"DB_CONN_STRING"`
	MigrationURL string `mapstructure:"MIGRATION_URL"`
}

func LoadBaseConfig(path string, configName string) (config *BaseConfig) {
	viper.AddConfigPath(path)
	viper.SetConfigName(configName)
	viper.SetTypeByDefaultValue(true)
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	LogAndPanicIfError(err, "failed when reading config")

	err = viper.Unmarshal(&config)
	LogAndPanicIfError(err, "failed when unmarshal config")

	return
}

func CheckAndSetConfig(path string, configName string, saFile ...string) *BaseConfig {
	config := LoadBaseConfig(path, configName)

	return config
}

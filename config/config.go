package config

import (
	"log"
	// "os"  <-- DELETE THIS LINE
	"github.com/spf13/viper"
)

type Config struct {
	DBHost         string `mapstructure:"DB_HOST"`
	DBUser         string `mapstructure:"DB_USER"`
	DBPassword     string `mapstructure:"DB_PASSWORD"`
	DBName         string `mapstructure:"DB_NAME"`
	DBPort         string `mapstructure:"DB_PORT"`
	HospitalAPIURL string `mapstructure:"HOSPITAL_A_API_URL"`
	JWTSecret      string `mapstructure:"JWT_SECRET"`
}

func LoadConfig() (Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	viper.SetDefault("DB_HOST", "db")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("HOSPITAL_A_API_URL", "https://hospital-a.api.co.th")
	viper.SetDefault("JWT_SECRET", "change_this_secret_in_prod")

	if err := viper.ReadInConfig(); err != nil {
		log.Println("Config file not found, using environment variables")
	}

	var config Config
	err := viper.Unmarshal(&config)
	return config, err
}
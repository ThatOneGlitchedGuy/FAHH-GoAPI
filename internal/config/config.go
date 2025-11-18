package config

import (
	"strings"
	"sync"

	"golang-app/internal/security"
	"github.com/spf13/viper"
)

type Settings struct {
	AppName                  string   `mapstructure:"APP_NAME"`
	Env                      string   `mapstructure:"ENV"`
	Debug                    bool     `mapstructure:"DEBUG"`
	APIV1Prefix              string   `mapstructure:"API_V1_PREFIX"`
	CORSOrigins              []string `mapstructure:"CORS_ORIGINS"`
	SecretKey                string   `mapstructure:"SECRET_KEY"`
	AccessTokenExpireMinutes int      `mapstructure:"ACCESS_TOKEN_EXPIRE_MINUTES"`
	JWTAlgorithm             string   `mapstructure:"JWT_ALGORITHM"`
	FernetKey                string   `mapstructure:"FERNET_KEY"`
	DatabaseURL              string   `mapstructure:"DATABASE_URL"`
	PageSizeDefault          int      `mapstructure:"PAGE_SIZE_DEFAULT"`
	PageSizeMax              int      `mapstructure:"PAGE_SIZE_MAX"`
}

var (
	settings    Settings
	once        sync.Once
	settingsErr error
)

func GetSettings() (*Settings, error) {
	once.Do(func() {
		viper.SetDefault("APP_NAME", "Mediated Marketplace API")
		viper.SetDefault("ENV", "dev")
		viper.SetDefault("DEBUG", true)
		viper.SetDefault("API_V1_PREFIX", "/api/v1")
		viper.SetDefault("CORS_ORIGINS", "http://localhost:3000,http://127.0.0.1:3000")
		viper.SetDefault("ACCESS_TOKEN_EXPIRE_MINUTES", 60*24)
		viper.SetDefault("JWT_ALGORITHM", "HS256")
		viper.SetDefault("DATABASE_URL", "sqlite://app.db")
		viper.SetDefault("PAGE_SIZE_DEFAULT", 20)
		viper.SetDefault("PAGE_SIZE_MAX", 100)

		secretKey, err := security.GenerateSecureKey(32)
		if err != nil {
			settingsErr = err
			return
		}
		viper.SetDefault("SECRET_KEY", secretKey)

		fernetKey, err := security.GenerateSecureKey(32)
		if err != nil {
			settingsErr = err
			return
		}
		viper.SetDefault("FERNET_KEY", fernetKey)

		viper.AutomaticEnv()
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

		viper.SetConfigFile(".env")
		viper.SetConfigType("env")
		if err := viper.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				// .env file not found, which is fine.
			} else {
				settingsErr = err
				return
			}
		}

		if err := viper.Unmarshal(&settings); err != nil {
			settingsErr = err
			return
		}

		if originsStr := viper.GetString("CORS_ORIGINS"); originsStr != "" && len(settings.CORSOrigins) == 0 {
			settings.CORSOrigins = strings.Split(originsStr, ",")
			for i, origin := range settings.CORSOrigins {
				settings.CORSOrigins[i] = strings.TrimSpace(origin)
			}
		}
	})
	return &settings, settingsErr
}

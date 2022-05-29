package router

import (
	"github.com/spf13/viper"
)

type Config struct {
	ApiUrl         string     `mapstructure:"ACCOUNT_API_URL" default:"/api/account"`
	Port           string     `mapstructure:"PORT" default:"8080"`
	MaxBodyBytes   int64      `mapstructure:"MAX_BODY_BYTES" default:"4194304"` // 4MB in Bytes ~ 4 * 1024 * 1024
	HandlerTimeout int64      `mapstructure:"HANDLER_TIMEOUT" default:"5"`
	DataSource     DataSource `mapstructure:"DATA_SOURCE,omitempty"`
	Token          Token      `mapstructure:"TOKEN,omitempty"`
}

type DataSource struct {
	PostGreSQL PostGreSQL `mapstructure:"POST_GRESQL,omitempty"`
	Redis      Redis      `mapstructure:"REDIS,omitempty"`
	GCP        GCP        `mapstructure:"GCP,omitempty"`
}

type PostGreSQL struct {
	PostGresHost              string `mapstructure:"POSTGRES_HOST" default:"postgres-account"`
	PostGresPort              string `mapstructure:"POSTGRES_PORT" default:"5432"`
	PostGresUser              string `mapstructure:"POSTGRES_USER" default:"postgres"`
	PostGresPassword          string `mapstructure:"POSTGRES_PASSWORD" required:"true"`
	PostGresDB                string `mapstructure:"POSTGRES_DB" default:"postgres"`
	PostGresSSL               string `mapstructure:"POSTGRES_SSL" default:"disable"`
	PostGresConnectionTimeOut int64  `mapstructure:"POSTGRES_CONNECTION_TIMEOUT" default:"10"`
}

type Redis struct {
	RedisHost string `mapstructure:"REDIS_HOST" default:"redis-account"`
	RedisPort string `mapstructure:"REDIS_PORT" default:"6379"`
}

type GCP struct {
	GCPImageBucket               string `mapstructure:"GCP_IMAGE_BUCKET" required:"true"`
	GoogleApplicationCredentials string `mapstructure:"GOOGLE_APPLICATION_CREDENTIALS" required:"true"`
	CloudConnectionTimeout       int64  `mapstructure:"CLOUD_CONNECTION_TIMEOUT" default:"5"`
}

type Token struct {
	AccessToken  AccessToken  `mapstructure:"ACCESS_TOKEN,omitempty"`
	RefreshToken RefreshToken `mapstructure:"REFRESH_TOKEN,omitempty"`
}

type AccessToken struct {
	AccessTokenExpire int64  `mapstructure:"ACCESS_TOKEN_EXPIRE" default:"900"` // 15 min in secs
	PublicKeyFile     string `mapstructure:"PUBLIC_KEY_FILE" required:"true"`
	PrivateKeyFile    string `mapstructure:"PRIVATE_KEY_FILE" required:"true"`
}

type RefreshToken struct {
	RefreshTokenExpire int64  `mapstructure:"REFRESH_TOKEN_EXPIRE" default:"259200"` // 3 days
	RefreshTokenSecret string `mapstructure:"REFRESH_TOKEN_SECRET" required:"true"`
}

func GetConfig(path string, name string, fileType string) (*Config, error) {
	var config *Config
	viper.AddConfigPath(path)
	viper.SetConfigName(name)
	viper.SetConfigType(fileType)

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

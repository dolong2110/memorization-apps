package router

import "github.com/spf13/viper"

type Config struct {
	APIUrl         string     `mapstructure:"WORD_API_URL" default:"/api/word"`
	MaxBodyBytes   int        `mapstructure:"MAX_BODY_BYTES" default:"4194304"`
	HandlerTimeOut int        `mapstructure:"HANDLER_TIME_OUT" default:"5"`
	DataSource     DataSource `mapstructure:"DATA_SOURCE"`
}

type DataSource struct {
	PostGreSQL PostGreSQL `mapstructure:"POST_GRESQL"`
}

type PostGreSQL struct {
	PostGreSHost     string `mapstructure:"POSTGRES_HOST" required:"true"`
	PostGreSPort     int    `mapstructure:"POSTGRES_PORT" required:"true"`
	PostGreSUser     string `mapstructure:"POSTGRES_USER" default:"postgres"`
	PostGreSPassword string `mapstructure:"POSTGRES_PASSWORD" required:"true"`
	PostGreSDB       string `mapstructure:"POSTGRES_DB" default:"postgres"`
	PostGreSSSL      string `mapstructure:"POSTGRES_SSL" default:"disable"`
}

func GetConfig(path string, name string, fileType string) (*Config, error) {
	var config *Config
	viper.AddConfigPath(path)
	viper.SetConfigName(name)     // Register config file name (no extension)
	viper.SetConfigType(fileType) // Look for specific type
	viper.ReadInConfig()

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

package router

type config struct {
	DataSource DataSource `mapstructure:"DATA_SOURCE"`
}

type DataSource struct {
	DataBase PostGreSQL `mapstructure:"POST_GRESQL"`
}

type PostGreSQL struct {
	PostGresHost     string `mapstructure:"POSTGRES_HOST" default:"postgres-account"`
	PostGresPort     int    `mapstructure:"POSTGRES_PORT" default:"5432"`
	PostGresUser     string `mapstructure:"POSTGRES_USER" default:"postgres"`
	PostGresPassword string `mapstructure:"POSTGRES_PASSWORD" required:"true"`
	PostGresDB       string `mapstructure:"POSTGRES_DB" default:"postgres"`
	PostGresSSL      string `mapstructure:"POSTGRES_SSL" default:"disable"`
}

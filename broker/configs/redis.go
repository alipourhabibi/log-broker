package configs

type RMQ struct {
	Host     string `mapstructure:"Host"`
	Port     string `mapstructure:"Port"`
	Password string `mapstructure:"Password"`
	DBName   int    `mapstructure:"DBName"`
}

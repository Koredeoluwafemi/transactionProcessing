package config

var App struct {
	Name   string `env:"Name" envDefault:"TransactionProcessing"`
	JWTKey string `env:"JWTKey" envDefault:"SGa37bXXtT1ZfkB1maTha3h9jLJQpEpd-dZ7aYqEvkB5M"`
	Mode   string `env:"Mode"  envDefault:"test"`
	Port   string `env:"Port" envDefault:"3200"`
	ENV    string `env:"ENV" envDefault:"local"`
}

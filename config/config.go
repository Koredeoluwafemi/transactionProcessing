package config

import "github.com/caarlos0/env/v6"

var Status struct {
	Active    string `env:"Active" envDefault:"active"`
	Initiated string `env:"Initiated" envDefault:"initiated"`
	Success   string `env:"Success" envDefault:"success"`
	Failed    string `env:"Failed" envDefault:"failed"`
}

func init() {
	_ = env.Parse(&App)
	_ = env.Parse(&Status)
}

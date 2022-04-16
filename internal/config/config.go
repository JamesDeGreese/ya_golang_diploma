package config

type Config struct {
	RunAddress           string `env:"RUN_ADDRESS" envDefault:"127.0.0.1:8080"`
	DatabaseURI          string `env:"DATABASE_URI" envDefault:"postgres://user:password@localhost:5432/ya_golang_diploma_db"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS" envDefault:"127.0.0.1:8081"`
}

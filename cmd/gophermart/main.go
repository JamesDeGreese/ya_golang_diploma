package main

import (
	"flag"

	"github.com/JamesDeGreese/ya_golang_diploma/internal/config"
	"github.com/JamesDeGreese/ya_golang_diploma/internal/database"
	router "github.com/JamesDeGreese/ya_golang_diploma/internal/http"
	"github.com/caarlos0/env/v6"
)

func main() {
	c := config.Config{}
	err := env.Parse(&c)
	if err != nil {
		panic(err)
	}

	flag.StringVar(&c.RunAddress, "a", c.RunAddress, "a 127.0.0.1:8080")
	flag.StringVar(&c.DatabaseURI, "d", c.DatabaseURI, "d postgres://username:password@host:port/database_name")
	flag.StringVar(&c.AccrualSystemAddress, "r", c.AccrualSystemAddress, "r 127.0.0.1:8081")
	flag.Parse()

	s := database.InitStorage(c)
	r := router.SetupRouter(c, s)

	err = r.Run(c.RunAddress)
	if err != nil {
		panic(err)
	}
}

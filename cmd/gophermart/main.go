package main

import (
	"flag"
	"go.uber.org/zap"

	"github.com/JamesDeGreese/ya_golang_diploma/internal/config"
	"github.com/JamesDeGreese/ya_golang_diploma/internal/database"
	"github.com/JamesDeGreese/ya_golang_diploma/internal/entities"
	router "github.com/JamesDeGreese/ya_golang_diploma/internal/http"
	"github.com/JamesDeGreese/ya_golang_diploma/internal/integrations"
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
	flag.StringVar(&c.AccrualSystemAddress, "r", c.AccrualSystemAddress, "r http://127.0.0.1:8081")
	flag.Parse()

	s := database.InitStorage(c)
	ur := entities.UserStorage{Storage: *s}
	or := entities.OrderStorage{Storage: *s}
	wr := entities.WithdrawnStorage{Storage: *s}
	as := integrations.AccrualService{Address: c.AccrualSystemAddress, OrderRepository: or}
	l, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	h := router.Handler{Config: c, UserRepository: ur, OrderRepository: or, WithdrawnRepository: wr}
	r := router.SetupRouter(as, h, ur, or, l)

	err = r.Run(c.RunAddress)
	if err != nil {
		panic(err)
	}
}

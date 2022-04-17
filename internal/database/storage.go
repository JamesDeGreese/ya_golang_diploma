package database

import (
	"context"
	"database/sql"

	"github.com/JamesDeGreese/ya_golang_diploma/internal/config"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source"
	"github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4"
)

type Storage struct {
	DBConn *pgx.Conn
}

func InitStorage(c config.Config) *Storage {
	makeMigration(c.DatabaseURI)
	conn, err := pgx.Connect(context.Background(), c.DatabaseURI)
	if err != nil {
		panic(err)
	}

	return &Storage{
		DBConn: conn,
	}
}

func makeMigration(uri string) {
	uri += "?sslmode=disable"
	source.Register("myfile", &file.File{})
	db, err := sql.Open("postgres", uri)
	_, err = postgres.WithInstance(db, &postgres.Config{})
	m, err := migrate.New(
		"myfile://internal/database/migrations",
		uri)
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		panic(err)
	}
}

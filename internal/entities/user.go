package entities

import (
	"context"
	"fmt"

	"github.com/JamesDeGreese/ya_golang_diploma/internal/auth"
	"github.com/JamesDeGreese/ya_golang_diploma/internal/database"
)

var tableName string = "users"

type User struct {
	ID        int
	Login     string
	Password  string
	AuthToken string
}

type UserRepository struct {
	Storage database.Storage
}

func (ur UserRepository) Add(login string, password string) (bool, error) {
	query := fmt.Sprintf("INSERT INTO %s (login, password) VALUES ('%s', '%s');", tableName, login, auth.MakeMD5(password))
	_, err := ur.Storage.DBConn.Exec(context.Background(), query)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (ur UserRepository) GetByLogin(login string) (User, error) {
	var res User
	query := fmt.Sprintf("SELECT id, login, password, auth_token FROM %s WHERE login = '%s';", tableName, login)
	err := ur.Storage.DBConn.QueryRow(context.Background(), query).Scan(&res.ID, &res.Login, &res.Password, &res.AuthToken)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (ur UserRepository) SetAuthToken(login string, token string) error {
	query := fmt.Sprintf("UPDATE %s set auth_token = %s WHERE login = '%s';", tableName, token, login)
	_, err := ur.Storage.DBConn.Exec(context.Background(), query)
	if err != nil {
		return err
	}

	return nil
}

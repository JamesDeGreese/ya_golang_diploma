package entities

import (
	"context"
	"fmt"

	"github.com/JamesDeGreese/ya_golang_diploma/internal/database"
)

type User struct {
	ID        int
	Login     string
	Password  string
	AuthToken string
	Balance   int
}

type UserRepository struct {
	Storage database.Storage
}

func (ur UserRepository) getTableName() string {
	return "users"
}

func (ur UserRepository) Add(login string, password string) (bool, error) {
	query := fmt.Sprintf("INSERT INTO %s (login, password) VALUES ('%s', '%s');", ur.getTableName(), login, password)
	_, err := ur.Storage.DBConn.Exec(context.Background(), query)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (ur UserRepository) GetByLogin(login string) (User, error) {
	var res User
	query := fmt.Sprintf("SELECT id, login, password, auth_token, balance FROM %s WHERE login = '%s';", ur.getTableName(), login)
	err := ur.Storage.DBConn.QueryRow(context.Background(), query).Scan(&res.ID, &res.Login, &res.Password, &res.AuthToken, &res.Balance)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (ur UserRepository) SetAuthToken(login string, token string) error {
	query := fmt.Sprintf("UPDATE %s set auth_token = '%s' WHERE login = '%s';", ur.getTableName(), token, login)
	_, err := ur.Storage.DBConn.Exec(context.Background(), query)
	if err != nil {
		return err
	}

	return nil
}

func (ur UserRepository) RecalculateBalance(login string) error {
	var or OrderRepository
	var balance int
	query := fmt.Sprintf("SELECT SUM(o.accrual) FROM %s u JOIN %s o WHERE o.status = 'REGISTERED' AND u.login = '%s';", ur.getTableName(), or.getTableName(), login)
	err := ur.Storage.DBConn.QueryRow(context.Background(), query).Scan(&balance)
	if err != nil {
		return err
	}
	query = fmt.Sprintf("UPDATE %s set balance = %d WHERE login = '%s';", ur.getTableName(), balance, login)
	_, err = ur.Storage.DBConn.Exec(context.Background(), query)
	if err != nil {
		return err
	}

	return nil
}

func (ur UserRepository) GetByToken(authToken string) (interface{}, interface{}) {
	var res User
	query := fmt.Sprintf("SELECT id, login, password, auth_token, balance / 100 FROM %s WHERE auth_token = '%s';", ur.getTableName(), authToken)
	err := ur.Storage.DBConn.QueryRow(context.Background(), query).Scan(&res.ID, &res.Login, &res.Password, &res.AuthToken, &res.Balance)
	if err != nil {
		return res, err
	}

	return res, nil
}

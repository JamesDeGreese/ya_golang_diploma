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
}

type UserRepository interface {
	Add(login string, password string) (bool, error)
	GetByLogin(login string) (User, error)
	SetAuthToken(login string, token string) error
	GetBalance(userID int) (int, error)
	GetByToken(authToken string) (interface{}, interface{})
}

type UserStorage struct {
	Storage database.Storage
}

func (ur UserStorage) getTableName() string {
	return "users"
}

func (ur UserStorage) Add(login string, password string) (bool, error) {
	query := fmt.Sprintf("INSERT INTO %s (login, password) VALUES ('%s', '%s');", ur.getTableName(), login, password)
	_, err := ur.Storage.DBConn.Exec(context.Background(), query)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (ur UserStorage) GetByLogin(login string) (User, error) {
	var res User
	query := fmt.Sprintf("SELECT id, login, password, auth_token FROM %s WHERE login = '%s';", ur.getTableName(), login)
	err := ur.Storage.DBConn.QueryRow(context.Background(), query).Scan(&res.ID, &res.Login, &res.Password, &res.AuthToken)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (ur UserStorage) SetAuthToken(login string, token string) error {
	query := fmt.Sprintf("UPDATE %s set auth_token = '%s' WHERE login = '%s';", ur.getTableName(), token, login)
	_, err := ur.Storage.DBConn.Exec(context.Background(), query)
	if err != nil {
		return err
	}

	return nil
}

func (ur UserStorage) GetBalance(userID int) (int, error) {
	var or OrderStorage
	var wr WithdrawnStorage
	var balance int
	var withdraw int
	query := fmt.Sprintf("SELECT COALESCE(SUM(o.accrual), 0) FROM %s o WHERE o.user_id = %d AND o.status = 'PROCESSED';", or.getTableName(), userID)
	err := ur.Storage.DBConn.QueryRow(context.Background(), query).Scan(&balance)
	if err != nil {
		return 0, err
	}
	query = fmt.Sprintf("SELECT COALESCE(SUM(w.sum), 0) FROM %s w WHERE w.user_id = %d;", wr.getTableName(), userID)
	err = ur.Storage.DBConn.QueryRow(context.Background(), query).Scan(&withdraw)
	if err != nil {
		return 0, err
	}

	return balance - withdraw, nil
}

func (ur UserStorage) GetByToken(authToken string) (interface{}, interface{}) {
	var res User
	query := fmt.Sprintf("SELECT id, login, password, auth_token FROM %s WHERE auth_token = '%s';", ur.getTableName(), authToken)
	err := ur.Storage.DBConn.QueryRow(context.Background(), query).Scan(&res.ID, &res.Login, &res.Password, &res.AuthToken)
	if err != nil {
		return res, err
	}

	return res, nil
}

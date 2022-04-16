package entities

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"

	"github.com/JamesDeGreese/ya_golang_diploma/internal/database"
)

var tableName string = "users"

type User struct {
	ID       string
	Login    string
	Password string
}

type UserRepository struct {
	Storage database.Storage
}

func (ur UserRepository) Add(login string, password string) (bool, error) {
	query := fmt.Sprintf("INSERT INTO %s (login, password) VALUES ($1, $2);", tableName)
	_, err := ur.Storage.DBConn.Exec(context.Background(),
		query,
		login, makeMD5(password))
	if err != nil {
		return false, err
	}

	return true, nil
}

func makeMD5(in string) string {
	binHash := md5.Sum([]byte(in))
	return hex.EncodeToString(binHash[:])
}

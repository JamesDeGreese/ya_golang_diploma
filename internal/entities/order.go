package entities

import (
	"context"
	"fmt"
	"time"

	"github.com/JamesDeGreese/ya_golang_diploma/internal/database"
	"github.com/jackc/pgtype"
)

type Order struct {
	ID         int
	UserID     int
	Number     string
	Status     string
	Accrual    int
	UploadedAt pgtype.Timestamp
}

type OrderRepository struct {
	Storage database.Storage
}

func (or OrderRepository) getTableName() string {
	return "orders"
}

func (or OrderRepository) Add(userID int, number string, status string, accrual int) (bool, error) {
	query := fmt.Sprintf("INSERT INTO %s (user_id, number, status, accrual, uploaded_at) VALUES (%d, '%s', '%s', %d, '%s');", or.getTableName(), userID, number, status, accrual, time.Now().Format(time.RFC3339))
	_, err := or.Storage.DBConn.Exec(context.Background(), query)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (or OrderRepository) GetByNumber(number string) (Order, error) {
	var res Order
	query := fmt.Sprintf("SELECT id, user_id, number, status, accrual, uploaded_at FROM %s WHERE number = '%s';", or.getTableName(), number)
	err := or.Storage.DBConn.QueryRow(context.Background(), query).Scan(&res.ID, &res.UserID, &res.Number, &res.Status, &res.Accrual, &res.UploadedAt)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (or OrderRepository) GetByUserID(userID int) ([]Order, error) {
	res := make([]Order, 0)
	query := fmt.Sprintf("SELECT id, user_id, number, status, accrual, uploaded_at FROM %s WHERE user_id = %d ORDER BY uploaded_at ASC;", or.getTableName(), userID)
	rows, err := or.Storage.DBConn.Query(context.Background(), query)
	if err != nil {
		return res, err
	}
	defer rows.Close()

	for rows.Next() {
		var r Order
		err := rows.Scan(&r.ID, &r.UserID, &r.Number, &r.Status, &r.Accrual, &r.UploadedAt)
		if err != nil {
			return nil, nil
		}
		res = append(res, r)
	}

	return res, nil
}

func (or OrderRepository) Update(number string, status string, accrual int) (bool, error) {
	query := fmt.Sprintf("UPDATE %s SET status = %s, accrual = %d WHERE number = '%s;", or.getTableName(), status, accrual, number)
	_, err := or.Storage.DBConn.Exec(context.Background(), query)
	if err != nil {
		return false, err
	}

	return true, nil
}

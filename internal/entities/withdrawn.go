package entities

import (
	"context"
	"fmt"
	"time"

	"github.com/JamesDeGreese/ya_golang_diploma/internal/database"
	"github.com/jackc/pgtype"
)

type Withdraw struct {
	ID          int
	UserID      int
	Order       string
	Sum         int
	ProcessedAt pgtype.Timestamp
}

type WithdrawnRepository struct {
	Storage database.Storage
}

func (wr WithdrawnRepository) getTableName() string {
	return "withdrawals"
}

func (wr WithdrawnRepository) Add(userID int, order string, sum int) (bool, error) {
	query := fmt.Sprintf("INSERT INTO %s (user_id, order_number, sum, processed_at) VALUES (%d, %s, %d, '%s');", wr.getTableName(), userID, order, sum, time.Now().Format(time.RFC3339))
	_, err := wr.Storage.DBConn.Exec(context.Background(), query)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (wr WithdrawnRepository) GetByOrderID(orderID int) (Withdraw, error) {
	var res Withdraw
	query := fmt.Sprintf("SELECT id, user_id, order_number, sum, processed_at FROM %s WHERE order_id = %d;", wr.getTableName(), orderID)
	err := wr.Storage.DBConn.QueryRow(context.Background(), query).Scan(&res.ID, &res.UserID, &res.Order, &res.Sum, &res.ProcessedAt)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (wr WithdrawnRepository) GetByUserID(userID int) ([]Withdraw, error) {
	res := make([]Withdraw, 0)
	query := fmt.Sprintf("SELECT id, user_id, order_number, sum, processed_at FROM %s WHERE user_id = %d ORDER BY w.processed_at ASC;", wr.getTableName(), userID)
	rows, err := wr.Storage.DBConn.Query(context.Background(), query)
	if err != nil {
		return res, err
	}
	defer rows.Close()

	for rows.Next() {
		var r Withdraw
		err := rows.Scan(&r.ID, &r.UserID, &r.Order, &r.Sum, &r.ProcessedAt)
		if err != nil {
			return nil, nil
		}
		res = append(res, r)
	}

	return res, nil
}

func (wr WithdrawnRepository) GetUserWithdrawnSum(userID int) (int, error) {
	var res int
	query := fmt.Sprintf("SELECT COALESCE(SUM(sum), 0) FROM %s WHERE user_id = %d;", wr.getTableName(), userID)
	err := wr.Storage.DBConn.QueryRow(context.Background(), query).Scan(&res)
	if err != nil {
		return 0, err
	}

	return res, nil
}

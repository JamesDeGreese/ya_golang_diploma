package integrations

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/JamesDeGreese/ya_golang_diploma/internal/database"
	"github.com/JamesDeGreese/ya_golang_diploma/internal/entities"
)

type Order struct {
	Number  string `json:"number"`
	Status  string `json:"status"`
	Accrual int    `json:"accrual"`
}

type AccrualService struct {
	Address string
	Storage *database.Storage
}

func (as AccrualService) getOrderInfo(orderNumber string) (Order, error) {
	var o Order
	r, err := http.Get(fmt.Sprintf("http://%s/api/orders/%s", as.Address, orderNumber))
	if err != nil {
		return o, errors.New("accrual response error")
	}
	if r.StatusCode != http.StatusOK {
		return o, nil
	}
	defer r.Body.Close()
	err = json.NewDecoder(r.Body).Decode(&o)
	if err != nil {
		return o, err
	}

	return o, nil
}

func (as AccrualService) SyncOrder(orderNumber string) error {
	or := entities.OrderRepository{Storage: *as.Storage}
	orderInfo, err := as.getOrderInfo(orderNumber)
	if err != nil {
		return err
	}
	success, err := or.Update(orderNumber, orderInfo.Status, orderInfo.Accrual*100)
	if !success || err != nil {
		return err
	}

	return nil
}

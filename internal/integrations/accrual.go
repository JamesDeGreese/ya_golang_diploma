package integrations

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/JamesDeGreese/ya_golang_diploma/internal/entities"
)

type Order struct {
	Number  string  `json:"number"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual"`
}

type AccrualService struct {
	Address         string
	OrderRepository entities.OrderRepository
}

func (as AccrualService) getOrderInfo(orderNumber string) (Order, error) {
	var o Order
	r, err := http.Get(fmt.Sprintf("%s/api/orders/%s", as.Address, orderNumber))
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
	orderInfo, err := as.getOrderInfo(orderNumber)
	if err != nil {
		return err
	}
	success, err := as.OrderRepository.Update(orderNumber, orderInfo.Status, int(orderInfo.Accrual*100))
	if !success || err != nil {
		return err
	}

	return nil
}

package integrations

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type Order struct {
	Number  string `json:"number"`
	Status  string `json:"status"`
	Accrual int    `json:"accrual"`
}

type AccrualService struct {
	Address string
}

func (as AccrualService) GetOrderInfo(orderNumber string) (Order, error) {
	var o Order
	r, err := http.Get(fmt.Sprintf("http://%s/api/orders/%s", as.Address, orderNumber))
	if err != nil || r.StatusCode != http.StatusOK {
		return o, errors.New("accrual response error")
	}
	defer r.Body.Close()
	err = json.NewDecoder(r.Body).Decode(&o)
	if err != nil {
		return o, err
	}

	return o, nil
}

package http

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/JamesDeGreese/ya_golang_diploma/internal/auth"
	"github.com/JamesDeGreese/ya_golang_diploma/internal/config"
	"github.com/JamesDeGreese/ya_golang_diploma/internal/entities"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"github.com/theplant/luhn"
)

type Handler struct {
	Config              config.Config
	UserRepository      entities.UserRepository
	OrderRepository     entities.OrderRepository
	WithdrawnRepository entities.WithdrawnRepository
}

func (h Handler) Dummy(c *gin.Context) {
	c.String(http.StatusOK, "")
}

func (h Handler) UserRegister(c *gin.Context) {
	var req RegisterRequest
	err := json.NewDecoder(c.Request.Body).Decode(&req)
	if err != nil {
		c.String(http.StatusBadRequest, "")
		return
	}

	user, err := h.UserRepository.GetByLogin(req.Login)
	if user.ID != 0 {
		c.String(http.StatusConflict, "")
		return
	}
	if err != nil && err != pgx.ErrNoRows {
		c.String(http.StatusInternalServerError, "")
		return
	}

	success, err := h.UserRepository.Add(req.Login, auth.MakeMD5(req.Password))
	if !success || err != nil {
		c.String(http.StatusInternalServerError, "")
		return
	}
	token := auth.GenerateAuthToken()
	err = h.UserRepository.SetAuthToken(req.Login, token)
	if err != nil {
		c.String(http.StatusInternalServerError, "")
		return
	}
	auth.SetAuthCookie(c, token)
	c.String(http.StatusOK, "")
}

func (h Handler) UserLogin(c *gin.Context) {
	var req LoginRequest
	err := json.NewDecoder(c.Request.Body).Decode(&req)
	if err != nil {
		c.String(http.StatusBadRequest, "")
		return
	}

	user, err := h.UserRepository.GetByLogin(req.Login)
	if err == pgx.ErrNoRows {
		c.String(http.StatusUnauthorized, "")
		return
	}

	if !auth.ComparePasswords(req.Password, user.Password) {
		c.String(http.StatusUnauthorized, "")
		return
	}

	token := auth.GenerateAuthToken()
	err = h.UserRepository.SetAuthToken(req.Login, token)
	if err != nil {
		c.String(http.StatusInternalServerError, "")
		return
	}
	auth.SetAuthCookie(c, token)

	c.String(http.StatusOK, "")
}

func (h Handler) OrderRegister(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusInternalServerError, "")
		return
	}
	orderNumber, err := strconv.Atoi(string(body))
	if err != nil {
		c.String(http.StatusInternalServerError, "")
		return
	}

	if !luhn.Valid(orderNumber) {
		c.String(http.StatusUnprocessableEntity, "")
		return
	}

	order, err := h.OrderRepository.GetByNumber(string(body))
	if err != nil && err != pgx.ErrNoRows {
		c.String(http.StatusInternalServerError, "")
		return
	}

	user, err := getUser(c)
	if err != nil {
		c.String(http.StatusInternalServerError, "")
		return
	}

	if order.ID != 0 && order.UserID == user.ID {
		c.String(http.StatusOK, "")
		return
	}
	if order.ID != 0 && order.UserID != user.ID {
		c.String(http.StatusConflict, "")
		return
	}

	success, err := h.OrderRepository.Add(user.ID, string(body), "NEW", 0)
	if !success || err != nil {
		c.String(http.StatusInternalServerError, "")
		return
	}

	c.String(http.StatusAccepted, "")
}

func (h Handler) OrdersGet(c *gin.Context) {
	user, err := getUser(c)
	if err != nil {
		c.String(http.StatusInternalServerError, "")
		return
	}
	orders, err := h.OrderRepository.GetByUserID(user.ID)
	if len(orders) == 0 {
		c.JSON(http.StatusNoContent, orders)
		return
	}
	if err != nil {
		c.String(http.StatusInternalServerError, "")
		return
	}

	res := make([]Order, 0)
	for _, o := range orders {
		res = append(res, Order{
			o.Number,
			o.Status,
			float32(o.Accrual) / 100,
			o.UploadedAt.Time.Format(time.RFC3339),
		})
	}

	c.JSON(http.StatusOK, res)
}

func (h Handler) BalanceGet(c *gin.Context) {
	user, err := getUser(c)
	if err != nil {
		c.String(http.StatusInternalServerError, "")
		return
	}

	withdrawn, err := h.WithdrawnRepository.GetUserWithdrawnSum(user.ID)
	if err != nil {
		c.String(http.StatusInternalServerError, "")
		return
	}

	balance, err := h.UserRepository.GetBalance(user.ID)
	if err != nil {
		c.String(http.StatusInternalServerError, "")
		return
	}
	res := Balance{
		float32(balance) / 100,
		float32(withdrawn) / 100,
	}

	c.JSON(http.StatusOK, res)
}

func (h Handler) WithdrawRegister(c *gin.Context) {
	var req WithdrawRequest
	err := json.NewDecoder(c.Request.Body).Decode(&req)
	if err != nil {
		c.String(http.StatusBadRequest, "")
		return
	}

	user, err := getUser(c)
	if err != nil {
		c.String(http.StatusInternalServerError, "")
		return
	}
	balance, err := h.UserRepository.GetBalance(user.ID)
	if err != nil {
		c.String(http.StatusInternalServerError, "")
		return
	}
	if balance < int(req.Sum*100) {
		c.String(http.StatusPaymentRequired, "")
	}

	orderNumber, err := strconv.Atoi(req.Order)
	if err != nil {
		c.String(http.StatusInternalServerError, "")
		return
	}

	if !luhn.Valid(orderNumber) {
		c.String(http.StatusUnprocessableEntity, "")
		return
	}

	success, err := h.WithdrawnRepository.Add(user.ID, req.Order, int(req.Sum*100))
	if !success || err != nil {
		c.String(http.StatusInternalServerError, "")
		return
	}

	c.String(http.StatusOK, "")
}

func (h Handler) WithdrawalsGet(c *gin.Context) {
	user, err := getUser(c)
	if err != nil {
		c.String(http.StatusInternalServerError, "")
		return
	}
	withdrawals, err := h.WithdrawnRepository.GetByUserID(user.ID)
	if len(withdrawals) == 0 {
		c.JSON(http.StatusNoContent, withdrawals)
		return
	}
	if err != nil {
		c.String(http.StatusInternalServerError, "")
		return
	}

	res := make([]Withdraw, len(withdrawals))
	for _, w := range withdrawals {
		res = append(res, Withdraw{
			w.Order,
			float32(w.Sum / 100),
			w.ProcessedAt.Time.Format(time.RFC3339),
		})
	}

	c.JSON(http.StatusOK, res)
}

func getUser(c *gin.Context) (entities.User, error) {
	u, exists := c.Get("user")
	if !exists {
		c.String(http.StatusUnauthorized, "")
		c.Abort()
	}

	user, ok := u.(entities.User)
	if !ok {
		return user, errors.New("can't get user")
	}

	return user, nil
}

package http

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/JamesDeGreese/ya_golang_diploma/internal/auth"
	"github.com/JamesDeGreese/ya_golang_diploma/internal/config"
	"github.com/JamesDeGreese/ya_golang_diploma/internal/database"
	"github.com/JamesDeGreese/ya_golang_diploma/internal/entities"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"github.com/theplant/luhn"
)

type Handler struct {
	Config  config.Config
	Storage *database.Storage
}

func (h Handler) Dummy(c *gin.Context) {
	c.String(http.StatusOK, "")
}

func (h Handler) UserRegister(c *gin.Context) {
	var req RegisterRequest
	ur := entities.UserRepository{Storage: *h.Storage}
	err := json.NewDecoder(c.Request.Body).Decode(&req)
	if err != nil {
		c.String(http.StatusBadRequest, "")
		return
	}

	user, err := ur.GetByLogin(req.Login)
	if user.ID != 0 {
		c.String(http.StatusConflict, "")
		return
	}
	if err != nil && err != pgx.ErrNoRows {
		c.String(http.StatusInternalServerError, "")
		return
	}

	success, err := ur.Add(req.Login, auth.MakeMD5(req.Password))
	if !success || err != nil {
		c.String(http.StatusInternalServerError, "")
		return
	}
	token := auth.GenerateAuthToken()
	err = ur.SetAuthToken(req.Login, token)
	if err != nil {
		c.String(http.StatusInternalServerError, "")
		return
	}
	auth.SetAuthCookie(c, token)
	c.String(http.StatusOK, "")
}

func (h Handler) UserLogin(c *gin.Context) {
	var req LoginRequest
	ur := entities.UserRepository{Storage: *h.Storage}
	err := json.NewDecoder(c.Request.Body).Decode(&req)
	if err != nil {
		c.String(http.StatusBadRequest, "")
		return
	}

	user, err := ur.GetByLogin(req.Login)
	if err == pgx.ErrNoRows {
		c.String(http.StatusUnauthorized, "")
		return
	}

	if !auth.ComparePasswords(req.Password, user.Password) {
		c.String(http.StatusUnauthorized, "")
		return
	}

	token := auth.GenerateAuthToken()
	err = ur.SetAuthToken(req.Login, token)
	if err != nil {
		c.String(http.StatusInternalServerError, "")
		return
	}
	auth.SetAuthCookie(c, token)

	c.String(http.StatusOK, "")
}

func (h Handler) OrderStore(c *gin.Context) {
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

	or := entities.OrderRepository{Storage: *h.Storage}
	order, err := or.GetByNumber(string(body))
	if err != nil && err != pgx.ErrNoRows {
		c.String(http.StatusInternalServerError, "")
		return
	}

	user := getUser(c)

	if order.ID != 0 && order.UserID == user.ID {
		c.String(http.StatusOK, "")
		return
	}
	if order.ID != 0 && order.UserID != user.ID {
		c.String(http.StatusConflict, "")
		return
	}

	success, err := or.Add(user.ID, string(body))
	if !success || err != nil {
		c.String(http.StatusInternalServerError, "")
		return
	}

	c.String(http.StatusAccepted, "")
}

func (h Handler) OrdersGet(c *gin.Context) {
	user := getUser(c)

	or := entities.OrderRepository{Storage: *h.Storage}
	orders, err := or.GetByUserID(user.ID)
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
			o.Accrual,
			o.UploadedAt.Time.Format(time.RFC3339),
		})
	}

	c.JSON(http.StatusOK, res)
}

func (h Handler) BalanceGet(c *gin.Context) {
	user := getUser(c)

	wr := entities.WithdrawnRepository{Storage: *h.Storage}
	withdrawn, err := wr.GetUserWithdrawnSum(user.ID)
	if err != nil {
		c.String(http.StatusInternalServerError, "")
		return
	}
	res := Balance{
		float32(user.Balance / 100),
		float32(withdrawn / 100),
	}

	c.JSON(http.StatusOK, res)
}

func (h Handler) WithdrawRegister(c *gin.Context) {
	var req WithdrawRequest
	wr := entities.WithdrawnRepository{Storage: *h.Storage}
	or := entities.OrderRepository{Storage: *h.Storage}
	err := json.NewDecoder(c.Request.Body).Decode(&req)
	if err != nil {
		c.String(http.StatusBadRequest, "")
		return
	}

	user := getUser(c)
	if user.Balance < int(req.Sum*100) {
		c.String(http.StatusPaymentRequired, "")
	}
	order, err := or.GetByNumber(req.Order)
	if order.ID == 0 || err != nil {
		c.String(http.StatusUnprocessableEntity, "")
		return
	}

	success, err := wr.Add(user.ID, order.ID, req.Sum)
	if !success || err != nil {
		c.String(http.StatusInternalServerError, "")
		return
	}

	c.String(http.StatusOK, "")
}

func (h Handler) WithdrawalsGet(c *gin.Context) {
	wr := entities.WithdrawnRepository{Storage: *h.Storage}
	user := getUser(c)
	withdrawals, err := wr.GetByUserID(user.ID)
	if len(withdrawals) == 0 {
		c.JSON(http.StatusNoContent, withdrawals)
		return
	}
	if err != nil {
		c.String(http.StatusInternalServerError, "")
		return
	}

	res := make([]Withdraw, 0)
	for _, w := range withdrawals {
		res = append(res, Withdraw{
			w.Order,
			float32(w.Sum / 100),
			w.ProcessedAt.Time.Format(time.RFC3339),
		})
	}

	c.JSON(http.StatusOK, res)
}

func getUser(c *gin.Context) entities.User {
	u, exists := c.Get("user")
	if !exists {
		c.String(http.StatusUnauthorized, "")
		c.Abort()
	}
	return u.(entities.User)
}

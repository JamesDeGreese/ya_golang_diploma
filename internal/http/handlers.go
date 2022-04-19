package responses

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

	u, exists := c.Get("user")
	if !exists {
		c.String(http.StatusUnauthorized, "")
		return
	}
	user := u.(entities.User)

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
	u, exists := c.Get("user")
	if !exists {
		c.String(http.StatusInternalServerError, "")
		return
	}
	user := u.(entities.User)

	or := entities.OrderRepository{Storage: *h.Storage}
	orders, err := or.GetByUserID(user.ID)
	if len(orders) == 0 {
		c.JSON(http.StatusNoContent, orders)
		return
	}
	if err != nil && err != pgx.ErrNoRows {
		c.String(http.StatusInternalServerError, "")
		return
	}

	res := make([]Order, 0)
	for _, o := range orders {
		if err != nil {
			c.String(http.StatusInternalServerError, "")
			return
		}
		res = append(res, Order{
			o.Number,
			o.Status,
			o.Accrual,
			o.UploadedAt.Time.Format(time.RFC3339),
		})
	}

	c.JSON(http.StatusOK, res)
}

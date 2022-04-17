package responses

import (
	"encoding/json"
	"net/http"

	"github.com/JamesDeGreese/ya_golang_diploma/internal/auth"
	"github.com/JamesDeGreese/ya_golang_diploma/internal/config"
	"github.com/JamesDeGreese/ya_golang_diploma/internal/database"
	"github.com/JamesDeGreese/ya_golang_diploma/internal/entities"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
)

type Handler struct {
	Config  config.Config
	Storage *database.Storage
}

func (h Handler) Dummy(c *gin.Context) {
	c.String(http.StatusOK, "")
	return
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

	success, err := ur.Add(req.Login, req.Password)
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
	return
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
	return
}

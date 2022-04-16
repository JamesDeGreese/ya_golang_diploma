package responses

import (
	"encoding/json"
	"net/http"

	"github.com/JamesDeGreese/ya_golang_diploma/internal/config"
	"github.com/JamesDeGreese/ya_golang_diploma/internal/database"
	"github.com/JamesDeGreese/ya_golang_diploma/internal/entities"
	"github.com/gin-gonic/gin"
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
	}

	add, err := ur.Add(req.Login, req.Password)
	if add == false || err != nil {
		c.String(http.StatusInternalServerError, "")
	}
	c.String(http.StatusOK, "")
	return
}

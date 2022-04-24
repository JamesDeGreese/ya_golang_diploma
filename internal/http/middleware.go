package http

import (
	"net/http"

	"github.com/JamesDeGreese/ya_golang_diploma/internal/database"
	"github.com/JamesDeGreese/ya_golang_diploma/internal/entities"
	"github.com/JamesDeGreese/ya_golang_diploma/internal/integrations"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(s *database.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("auth")
		if cookie == "" || err != nil {
			if err != nil {
				c.String(http.StatusUnauthorized, "")
				c.Abort()
			}
		} else {
			ur := entities.UserRepository{Storage: *s}
			user, err := ur.GetByToken(cookie)
			if err != nil {
				c.String(http.StatusInternalServerError, "")
				c.Abort()
			}
			c.Set("user", user)
		}
		c.Next()
	}
}

func OrdersSyncMiddleware(s *database.Storage, as integrations.AccrualService) gin.HandlerFunc {
	return func(c *gin.Context) {

		u, exists := c.Get("user")
		if !exists {
			c.String(http.StatusUnauthorized, "")
			c.Abort()
		}
		user := u.(entities.User)
		ur := entities.UserRepository{Storage: *s}
		or := entities.OrderRepository{Storage: *s}
		orders, err := or.GetByUserID(user.ID)
		if err != nil {
			c.String(http.StatusInternalServerError, "")
			c.Abort()
		}
		for _, order := range orders {
			go func() {
				err := as.SyncOrder(order.Number)
				if err != nil {
					return
				}
				err = ur.RecalculateBalance(user.Login)
				if err != nil {
					return
				}
			}()
		}
		c.Next()
	}
}

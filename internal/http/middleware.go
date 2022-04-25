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
		or := entities.OrderRepository{Storage: *s}
		orders, err := or.GetUserNonFinalOrders(user.ID)
		if err != nil {
			c.String(http.StatusInternalServerError, "")
			c.Abort()
		}
		go func() {
			for _, order := range orders {
				_ = as.SyncOrder(order.Number)
			}
		}()

		c.Next()
	}
}

package auth

import (
	"net/http"

	"github.com/JamesDeGreese/ya_golang_diploma/internal/database"
	"github.com/JamesDeGreese/ya_golang_diploma/internal/entities"
	"github.com/gin-gonic/gin"
)

func Middleware(s *database.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("auth")
		if cookie == "" || err != nil {
			if err != nil {
				c.String(http.StatusUnauthorized, "")
				return
			}
		} else {
			ur := entities.UserRepository{Storage: *s}
			user, err := ur.GetByToken(cookie)
			if err != nil {
				c.String(http.StatusInternalServerError, "")
				return
			}
			c.Set("user", user)
		}
		c.Next()
	}
}

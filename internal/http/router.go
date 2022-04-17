package responses

import (
	"github.com/JamesDeGreese/ya_golang_diploma/internal/config"
	"github.com/JamesDeGreese/ya_golang_diploma/internal/database"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

func SetupRouter(c config.Config, s *database.Storage) *gin.Engine {
	r := gin.Default()
	r.Use(gzip.Gzip(gzip.BestSpeed, gzip.WithDecompressFn(gzip.DefaultDecompressHandle)))
	h := Handler{
		Config:  c,
		Storage: s,
	}
	r.POST("/api/user/register", h.UserRegister)
	r.POST("/api/user/login", h.UserLogin)
	r.POST("/api/user/orders", h.Dummy)
	r.GET("/api/user/orders", h.Dummy)
	r.GET("/api/user/balance", h.Dummy)
	r.POST("/api/user/balance/withdraw", h.Dummy)
	r.GET("/api/user/balance/withdrawals", h.Dummy)
	return r
}

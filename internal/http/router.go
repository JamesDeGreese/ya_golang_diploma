package responses

import (
	"github.com/JamesDeGreese/ya_golang_diploma/internal/auth"
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

	authorized := r.Group("/").Use(auth.Middleware(s))

	authorized.POST("/api/user/orders", h.OrderStore)
	authorized.GET("/api/user/orders", h.OrdersGet)
	authorized.GET("/api/user/balance", h.Dummy)
	authorized.POST("/api/user/balance/withdraw", h.Dummy)
	authorized.GET("/api/user/balance/withdrawals", h.Dummy)

	return r
}

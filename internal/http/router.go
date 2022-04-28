package http

import (
	"github.com/JamesDeGreese/ya_golang_diploma/internal/entities"
	"github.com/JamesDeGreese/ya_golang_diploma/internal/integrations"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

func SetupRouter(as integrations.AccrualService, h Handler, ur entities.UserRepository, or entities.OrderRepository) *gin.Engine {
	r := gin.Default()
	r.Use(gzip.Gzip(gzip.BestSpeed, gzip.WithDecompressFn(gzip.DefaultDecompressHandle)))
	r.POST("/api/user/register", h.UserRegister)
	r.POST("/api/user/login", h.UserLogin)

	authorized := r.Group("/").Use(AuthMiddleware(ur), OrdersSyncMiddleware(or, as))

	authorized.POST("/api/user/orders", h.OrderRegister)
	authorized.GET("/api/user/orders", h.OrdersGet)
	authorized.GET("/api/user/balance", h.BalanceGet)
	authorized.POST("/api/user/balance/withdraw", h.WithdrawRegister)
	authorized.GET("/api/user/balance/withdrawals", h.WithdrawalsGet)

	return r
}

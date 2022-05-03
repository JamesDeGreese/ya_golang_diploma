package http

import (
	"time"

	"github.com/JamesDeGreese/ya_golang_diploma/internal/entities"
	"github.com/JamesDeGreese/ya_golang_diploma/internal/integrations"
	"github.com/gin-contrib/gzip"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func SetupRouter(as integrations.AccrualService, h Handler, ur entities.UserStorage, or entities.OrderStorage) *gin.Engine {
	logger, _ := zap.NewProduction()

	r := gin.Default()

	r.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	r.Use(ginzap.RecoveryWithZap(logger, true))

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

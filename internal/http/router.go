package http

import (
	"go.uber.org/zap"
	"time"

	"github.com/JamesDeGreese/ya_golang_diploma/internal/entities"
	"github.com/JamesDeGreese/ya_golang_diploma/internal/integrations"
	"github.com/gin-contrib/gzip"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
)

func SetupRouter(as integrations.AccrualService, h Handler, ur entities.UserStorage, or entities.OrderStorage, l *zap.Logger) *gin.Engine {

	r := gin.Default()

	r.Use(ginzap.Ginzap(l, time.RFC3339, true))
	r.Use(ginzap.RecoveryWithZap(l, true))

	r.Use(gzip.Gzip(gzip.BestSpeed, gzip.WithDecompressFn(gzip.DefaultDecompressHandle)))
	r.POST("/api/user/register", h.UserRegister)
	r.POST("/api/user/login", h.UserLogin)

	authorized := r.Group("/").Use(AuthMiddleware(ur))

	authorized.POST("/api/user/orders", h.OrderRegister)
	authorized.GET("/api/user/orders", OrdersSyncMiddleware(or, as, l), h.OrdersGet)
	authorized.GET("/api/user/balance", OrdersSyncMiddleware(or, as, l), h.BalanceGet)
	authorized.POST("/api/user/balance/withdraw", OrdersSyncMiddleware(or, as, l), h.WithdrawRegister)
	authorized.GET("/api/user/balance/withdrawals", OrdersSyncMiddleware(or, as, l), h.WithdrawalsGet)

	return r
}

package router

import (
	"count/backend/internal/handler"
	"count/backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func New(h *handler.Handler) *gin.Engine {
	r := gin.Default()
	r.Use(middleware.CORS())

	api := r.Group("/api")
	api.GET("/health", handler.Health)
	api.POST("/auth/register", h.Register)
	api.POST("/auth/login", h.Login)

	authGroup := api.Group("")
	authGroup.Use(middleware.Auth(h.JWTSecret))
	authGroup.GET("/auth/me", h.Me)
	authGroup.GET("/processes", h.ListProcesses)
	authGroup.GET("/me/processes", h.ListMyProcesses)
	authGroup.GET("/scans/preview/:qrToken", h.PreviewScanOrder)
	authGroup.POST("/scans/record", h.RecordScan)
	authGroup.GET("/scan-records/mine", h.ListMyScanRecords)
	authGroup.GET("/orders/:id", h.GetOrder)
	authGroup.GET("/orders/by-no/:orderNo", h.GetOrderByNo)
	authGroup.GET("/orders", h.ListMyOrders)

	admin := api.Group("/admin")
	admin.Use(middleware.Auth(h.JWTSecret), middleware.RequireRole("admin"))
	admin.GET("/workers", h.ListWorkers)
	admin.PUT("/workers/:id/processes", h.UpdateWorkerProcesses)
	admin.POST("/processes", h.CreateProcess)
	admin.POST("/orders", h.CreateOrder)
	admin.GET("/orders", h.ListOrders)
	admin.GET("/scan-records", h.ListAllScanRecords)

	return r
}

package router

import (
	"github.com/gin-gonic/gin"
	"gme/controller"
	"net/http"
)

func SetupRouter(router *gin.Engine) {
	router.GET("/orders", func(c *gin.Context) {
		c.JSON(http.StatusOK, controller.Order.GetOrders)
	})
}

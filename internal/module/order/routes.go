package order

import "github.com/gin-gonic/gin"

func RegisterOrderRoutes(r gin.IRoutes, ctrl *OrderController) {
	r.GET("/orders/:id", ctrl.GetOrder)
	r.POST("/orders", ctrl.CreateOrder)
}

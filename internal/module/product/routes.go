package product

import "github.com/gin-gonic/gin"

func RegisterProductRoutes(r gin.IRoutes, ctrl *ProductController) {
	r.GET("/products", ctrl.ListProducts)
	r.GET("/products/:id", ctrl.GetProduct)
	r.POST("/products", ctrl.CreateProduct)
	r.PUT("/products/:id", ctrl.UpdateProduct)
}

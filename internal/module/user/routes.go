package user

import "github.com/gin-gonic/gin"

func RegisterUserRoutes(r gin.IRoutes, ctrl *UserController) {
	r.GET("/users/:id", ctrl.GetUser)
	r.POST("/users", ctrl.CreateUser)
}

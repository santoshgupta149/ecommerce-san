package admin

import (
	// "net/http"

	"github.com/gin-gonic/gin"
)

func RegisterAdminRoutes(r gin.IRoutes, ctrl *AdminController) {
	r.GET("/admin/get-all", ctrl.GetAllAdmin)
	r.POST("/admin/create", ctrl.CreateAdmin)
	// r.PUT("/admin/update/:id", ctrl.UpdateAdmin)
	// r.DELETE("/admin/delete/:id", ctrl.DeleteAdmin)
	// r.GET("/admin/get/:id", ctrl.GetAdminByID)
}

package admin

import (
	"net/http"

	"ecommerce-go/internal/module/admin/dto"
	"ecommerce-go/pkg/httputil"
	"ecommerce-go/pkg/validator"

	"github.com/gin-gonic/gin"
)

type AdminController struct {
	svc *AdminService // ← controller owns a reference to service
}

func NewController(svc *AdminService) *AdminController {
	return &AdminController{svc: svc}
}

// Controller's ONLY job: parse HTTP → call service → format response
func (ac *AdminController) CreateAdmin(c *gin.Context) {
	ctx := c.Request.Context()
	// Step 1: bind HTTP request into DTO
	var req dto.CreateAdminRequest
	// if err := c.ShouldBindJSON(&req); err != nil {
	//     httputil.ErrorMessage(c, http.StatusBadRequest, err.Error())
	//     return
	// }

	if err := c.ShouldBindJSON(&req); err != nil {
		// Format ugly error → clean array of field errors
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"errors":  validator.FormatErrors(err),
		})
		return
	}

	// Step 2: hand off to service — controller doesn't care HOW it's done
	result, err := ac.svc.CreateAdmin(ctx, req)
	if err != nil {
		httputil.ErrorMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Step 3: format and send HTTP response
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    result,
	})
}

func (ac *AdminController) GetAllAdmin(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"ok":   true,
		"data": []gin.H{},
		"meta": gin.H{
			"message":  "wire service/repository to return rows",
			"resource": "admins",
		},
	})
}

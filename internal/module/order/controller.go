package order

import (
	"errors"
	"net/http"
	"strconv"

	"ecommerce-go/pkg/httputil"

	"github.com/gin-gonic/gin"
)

type OrderController struct {
	svc *OrderService
}

func NewController(svc *OrderService) *OrderController {
	return &OrderController{svc: svc}
}

type createOrderRequest struct {
	UserID int64   `json:"user_id" binding:"required"`
	Total  float64 `json:"total" binding:"required,gt=0"`
}

func (oc *OrderController) CreateOrder(c *gin.Context) {
	ctx := c.Request.Context()

	var req createOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.ErrorMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	o, err := oc.svc.Create(ctx, CreateOrderInput{
		UserID: req.UserID,
		Total:  req.Total,
	})
	if err != nil {
		httputil.Error(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, o)
}

func (oc *OrderController) GetOrder(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		httputil.ErrorMessage(c, http.StatusBadRequest, "invalid id")
		return
	}

	o, err := oc.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			httputil.ErrorMessage(c, http.StatusNotFound, "order not found")
			return
		}
		httputil.Error(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, o)
}

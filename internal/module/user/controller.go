package user

import (
	"errors"
	"net/http"
	"strconv"

	"ecommerce-go/pkg/httputil"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	svc *UserService
}

func NewController(svc *UserService) *UserController {
	return &UserController{svc: svc}
}

type createUserRequest struct {
	Name   string `json:"name" binding:"required"`
	Email  string `json:"email" binding:"required,email"`
	Mobile string `json:"mobile_number" binding:"required,min=10,max=20"`
}

func (uc *UserController) CreateUser(c *gin.Context) {
	ctx := c.Request.Context()

	var req createUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.ErrorMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	u, err := uc.svc.Create(ctx, CreateUserInput{
		Name:   req.Name,
		Email:  req.Email,
		Mobile: req.Mobile,
	})
	if err != nil {
		if errors.Is(err, ErrDuplicateEmail) || errors.Is(err, ErrDuplicateMobile) {
			httputil.ErrorMessage(c, http.StatusConflict, err.Error())
			return
		}
		httputil.Error(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, u)
}

func (uc *UserController) GetUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		httputil.ErrorMessage(c, http.StatusBadRequest, "invalid id")
		return
	}

	u, err := uc.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			httputil.ErrorMessage(c, http.StatusNotFound, "user not found")
			return
		}
		httputil.Error(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, u)
}

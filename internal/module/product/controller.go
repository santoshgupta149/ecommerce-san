package product

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"ecommerce-go/pkg/httputil"
)

type ProductController struct {
	svc *ProductService
}

func NewController(svc *ProductService) *ProductController {
	return &ProductController{svc: svc}
}

func (pc *ProductController) CreateProduct(c *gin.Context) {
	var req createProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.ErrorMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	p, err := pc.svc.Create(c.Request.Context(), CreateProductInput{
		Name:        req.Name,
		SKU:         req.SKU,
		Description: req.Description,
		Category:    req.Category,
		Brand:       req.Brand,
		Price:       req.Price,
		Stock:       req.Stock,
		ImageURL:    req.ImageURL,
		IsActive:    req.IsActive,
	})
	if err != nil {
		var validationErr *ValidationError
		if errors.As(err, &validationErr) {
			httputil.ErrorMessage(c, http.StatusBadRequest, validationErr.Error())
			return
		}
		if errors.Is(err, ErrDuplicateSKU) {
			httputil.ErrorMessage(c, http.StatusConflict, err.Error())
			return
		}
		httputil.Error(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, p)
}

func (pc *ProductController) GetProduct(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		httputil.ErrorMessage(c, http.StatusBadRequest, "invalid id")
		return
	}

	p, err := pc.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			httputil.ErrorMessage(c, http.StatusNotFound, "product not found")
			return
		}
		httputil.Error(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, p)
}

func (pc *ProductController) ListProducts(c *gin.Context) {
	products, err := pc.svc.List(c.Request.Context())
	if err != nil {
		httputil.Error(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, products)
}

func (pc *ProductController) UpdateProduct(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		httputil.ErrorMessage(c, http.StatusBadRequest, "invalid id")
		return
	}

	var req updateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.ErrorMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	p, err := pc.svc.Update(c.Request.Context(), id, UpdateProductInput{
		Name:        req.Name,
		SKU:         req.SKU,
		Description: req.Description,
		Category:    req.Category,
		Brand:       req.Brand,
		Price:       req.Price,
		Stock:       req.Stock,
		ImageURL:    req.ImageURL,
		IsActive:    req.IsActive,
	})
	if err != nil {
		var validationErr *ValidationError
		if errors.As(err, &validationErr) {
			httputil.ErrorMessage(c, http.StatusBadRequest, validationErr.Error())
			return
		}
		if errors.Is(err, ErrNotFound) {
			httputil.ErrorMessage(c, http.StatusNotFound, "product not found")
			return
		}
		if errors.Is(err, ErrDuplicateSKU) {
			httputil.ErrorMessage(c, http.StatusConflict, err.Error())
			return
		}
		httputil.Error(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, p)
}

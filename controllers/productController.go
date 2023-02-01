package controllers

import (
	"net/http"
	"strconv"

	"github.com/aldisaputra17/shoping-cart/entities"
	"github.com/aldisaputra17/shoping-cart/helpers"
	"github.com/aldisaputra17/shoping-cart/request"
	"github.com/aldisaputra17/shoping-cart/services"
	"github.com/gin-gonic/gin"
)

type ProductController interface {
	AddCart(ctx *gin.Context)
	Deleted(ctx *gin.Context)
	GetCarts(ctx *gin.Context)
	Pagination(ctx *gin.Context)
}

type productController struct {
	productService services.ProductService
}

func NewProductController(prodService services.ProductService) ProductController {
	return &productController{
		productService: prodService,
	}
}

func (c *productController) AddCart(ctx *gin.Context) {
	var req request.ProductRequest

	err := ctx.ShouldBind(&req)
	if err != nil {
		res := helpers.BuildErrorResponse("Failed input data", err.Error(), helpers.EmptyObj{})
		ctx.JSON(http.StatusUnprocessableEntity, res)
		return
	} else {
		result, err := c.productService.AddCart(ctx, &req)
		if err != nil {
			res := helpers.BuildErrorResponse("Failed add cart", err.Error(), helpers.EmptyObj{})
			ctx.JSON(http.StatusBadRequest, res)
			return
		}
		response := helpers.BuildResponse(true, "Created", result)
		ctx.JSON(http.StatusCreated, response)
	}

}

func (c *productController) Deleted(ctx *gin.Context) {
	var prod entities.Products
	idpr := ctx.Param("id")
	id, _ := strconv.ParseInt(idpr, 0, 0)
	if prod.ID == int(id) {
		c.productService.Deleted(ctx, prod)
		res := helpers.BuildResponse(true, "Deleted", helpers.EmptyObj{})
		ctx.JSON(http.StatusOK, res)
	} else {
		response := helpers.BuildErrorResponse("Failed deleted cart", "deleted again", helpers.EmptyObj{})
		ctx.JSON(http.StatusForbidden, response)
	}

}

func (c *productController) GetCarts(ctx *gin.Context) {
	result, err := c.productService.GetCarts(ctx)
	if err != nil {
		res := helpers.BuildErrorResponse("Failed find carts", err.Error(), helpers.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
	}
	res := helpers.BuildResponse(true, "Ok", result)
	ctx.JSON(http.StatusOK, res)
}

func (c *productController) Pagination(ctx *gin.Context) {
	code := http.StatusOK
	pagination := helpers.GeneratePaginationRequest(ctx)

	response, err := c.productService.Paginantion(ctx, pagination)
	if err != nil {
		res := helpers.BuildErrorResponse("Failed pagination products", err.Error(), helpers.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	if !response.Success {
		code = http.StatusBadRequest
	}

	res := helpers.BuildResponse(true, "Ok", response)
	ctx.AbortWithStatusJSON(code, res)
}

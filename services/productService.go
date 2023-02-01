package services

import (
	"context"
	"fmt"
	"time"

	"github.com/aldisaputra17/shoping-cart/entities"
	"github.com/aldisaputra17/shoping-cart/helpers"
	"github.com/aldisaputra17/shoping-cart/repositories"
	"github.com/aldisaputra17/shoping-cart/request"
	"github.com/gin-gonic/gin"
)

type ProductService interface {
	AddCart(ctx context.Context, req *request.ProductRequest) (*entities.Products, error)
	Deleted(ctx context.Context, prod entities.Products) error
	GetCarts(ctx context.Context) (*[]entities.Products, error)
	Paginantion(ctx *gin.Context, paginat *entities.Pagination) (helpers.Response, error)
}

type productService struct {
	productRepository repositories.ProductRepository
	contextTimeOut    time.Duration
}

func NewProductService(repo repositories.ProductRepository, time time.Duration) ProductService {
	return &productService{
		productRepository: repo,
		contextTimeOut:    time,
	}
}

func (service *productService) AddCart(ctx context.Context, req *request.ProductRequest) (*entities.Products, error) {
	prodCreate := &entities.Products{
		ID:           req.ID,
		Name:         req.Name,
		Product_Code: req.ProductCode,
		Quantity:     req.Quantity,
	}
	ctx, cancel := context.WithTimeout(ctx, service.contextTimeOut)
	defer cancel()

	res, err := service.productRepository.AddCart(ctx, prodCreate)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (service *productService) Deleted(ctx context.Context, prod entities.Products) error {
	return service.productRepository.Deleted(ctx, prod)
}

func (service *productService) GetCarts(ctx context.Context) (*[]entities.Products, error) {
	ctx, cancel := context.WithTimeout(ctx, service.contextTimeOut)
	defer cancel()
	res, err := service.productRepository.GetCarts(ctx)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (service *productService) Paginantion(ctx *gin.Context, paginat *entities.Pagination) (helpers.Response, error) {
	operationResult, totalPages := service.productRepository.Pagination(paginat)

	if operationResult.Error != nil {
		return helpers.Response{Success: true, Message: operationResult.Error.Error()}, nil
	}

	var data = operationResult.Result.(*entities.Pagination)

	urlPath := ctx.Request.URL.Path

	searchQueryParams := ""

	for _, search := range paginat.Searchs {
		searchQueryParams += fmt.Sprintf("&%s.%s=%s", search.Column, search.Action, search.Query)
	}

	data.FirstPage = fmt.Sprintf("%s?limit=%d&page=%d&sort_product=%s", urlPath, paginat.Limit, 0, paginat.SortName) + searchQueryParams
	data.LastPage = fmt.Sprintf("%s?limit=%d&page=%d&sort_product=%s", urlPath, paginat.Limit, totalPages, paginat.SortName) + searchQueryParams

	if data.Page > 0 {
		// set previous page pagination response
		data.PreviousPage = fmt.Sprintf("%s?limit=%d&page=%d&sort_product=%s", urlPath, paginat.Limit, data.Page-1, paginat.SortName) + searchQueryParams
	}

	if data.Page < totalPages {
		// set next page pagination response
		data.NextPage = fmt.Sprintf("%s?limit=%d&page=%d&sort_product=%s", urlPath, paginat.Limit, data.Page+1, paginat.SortName) + searchQueryParams
	}

	if data.Page > totalPages {
		// reset previous page
		data.PreviousPage = ""
	}
	return helpers.BuildResponse(true, "Ok", data), nil
}

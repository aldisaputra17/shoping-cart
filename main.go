package main

import (
	"fmt"
	"time"

	"github.com/aldisaputra17/shoping-cart/controllers"
	"github.com/aldisaputra17/shoping-cart/database"
	"github.com/aldisaputra17/shoping-cart/repositories"
	"github.com/aldisaputra17/shoping-cart/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	contextTimeOut    time.Duration                  = 10 * time.Second
	db                *gorm.DB                       = database.ConnectDB()
	productRepository repositories.ProductRepository = repositories.NewProductRepository(db)
	productService    services.ProductService        = services.NewProductService(productRepository, contextTimeOut)
	productController controllers.ProductController  = controllers.NewProductController(productService)
)

func main() {
	fmt.Println("Starting Server")
	defer database.CloseDatabaseConnection(db)

	r := gin.Default()

	api := r.Group("api")

	productRoutes := api.Group("/product")
	{
		productRoutes.GET("/", productController.GetCarts)
		productRoutes.POST("/", productController.AddCart)
		productRoutes.DELETE("/:id", productController.Deleted)
		productRoutes.GET("/search", productController.Pagination)
	}
	r.Run()
}

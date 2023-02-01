package repositories

import (
	"context"
	"fmt"
	"math"
	"strings"

	"github.com/aldisaputra17/shoping-cart/entities"
	"github.com/aldisaputra17/shoping-cart/helpers"
	"gorm.io/gorm"
)

type ProductRepository interface {
	AddCart(ctx context.Context, prod *entities.Products) (*entities.Products, error)
	Update(ctx context.Context, prod *entities.Products) (*entities.Products, error)
	GetProdId(ctx context.Context, id int) ([]*entities.Products, error)
	GetCarts(ctx context.Context) (*[]entities.Products, error)
	Deleted(ctx context.Context, prod entities.Products) error
	Pagination(pagination *entities.Pagination) (helpers.PaginationResult, int)
}

type productConnection struct {
	connection *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productConnection{
		connection: db,
	}
}

func (db *productConnection) GetProdId(ctx context.Context, id int) ([]*entities.Products, error) {
	var prod []*entities.Products

	res := db.connection.WithContext(ctx).Where("id = ?", id).Find(&prod)
	if res.Error != nil {
		return nil, res.Error
	}
	return prod, nil
}

func (db *productConnection) AddCart(ctx context.Context, prod *entities.Products) (*entities.Products, error) {
	prodItem, err := db.GetProdId(ctx, prod.ID)
	if err != nil {
		return nil, err
	}

	if len(prodItem) > 0 {
		fmt.Println("Item already exists in cart.")
		prod.Quantity = prodItem[0].Quantity + 1
	}

	res := db.connection.WithContext(ctx).Create(&prod)
	if res.Error != nil {
		return nil, res.Error
	}
	return prod, nil
}

func (db *productConnection) Update(ctx context.Context, prod *entities.Products) (*entities.Products, error) {
	res := db.connection.WithContext(ctx).Model(&prod).Updates(entities.Products{
		Quantity: prod.Quantity,
	})
	if res.Error != nil {
		return nil, res.Error
	}
	return prod, nil
}

func (db *productConnection) Deleted(ctx context.Context, prod entities.Products) error {
	res := db.connection.WithContext(ctx).Delete(prod)
	if res.Error != nil {
		return nil
	}
	return nil
}

func (db *productConnection) GetCarts(ctx context.Context) (*[]entities.Products, error) {
	var prods *[]entities.Products
	res := db.connection.WithContext(ctx).Find(&prods)
	if res.Error != nil {
		return nil, res.Error
	}
	return prods, nil
}

func (db *productConnection) Pagination(pagination *entities.Pagination) (helpers.PaginationResult, int) {
	var prod []entities.Products

	var (
		totalRows  int64
		totalPages int
		fromRow    int
		toRow      int
	)

	offset := pagination.Page * pagination.Limit

	find := db.connection.Limit(pagination.Limit).Offset(offset).Order(pagination.SortName)

	searchs := pagination.Searchs

	for _, value := range searchs {
		column := value.Column
		action := value.Action
		query := value.Query

		switch action {
		case "equals":
			whereQuery := fmt.Sprintf("%s = ?", column)
			find = find.Where(whereQuery, query)
		case "contains":
			whereQuery := fmt.Sprintf("%s LIKE ?", column)
			find = find.Where(whereQuery, "%"+query+"%")
		case "in":
			whereQuery := fmt.Sprintf("%s IN (?)", column)
			queryArray := strings.Split(query, ",")
			find = find.Where(whereQuery, queryArray)
		}
	}

	find = find.Find(&prod)

	errFind := find.Error

	if errFind != nil {
		return helpers.PaginationResult{Error: errFind}, totalPages
	}

	pagination.Rows = prod

	errCount := db.connection.Model(&entities.Products{}).Count(&totalRows).Error

	if errCount != nil {
		return helpers.PaginationResult{Error: errCount}, totalPages
	}

	pagination.TotalRows = int(totalRows)

	totalPages = int(math.Ceil(float64(totalRows)/float64(pagination.Limit))) - 1

	if pagination.Page == 0 {
		fromRow = 1
		toRow = pagination.Limit
	} else {
		if pagination.Page <= totalPages {
			fromRow = pagination.Page*pagination.Limit + 1
			toRow = (pagination.Page + 1) * pagination.Limit
		}
	}

	if toRow > int(totalRows) {
		toRow = int(totalRows)
	}

	pagination.FromRow = fromRow
	pagination.ToRow = toRow

	return helpers.PaginationResult{Result: pagination}, totalPages
}

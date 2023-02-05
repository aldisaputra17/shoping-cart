package mock

import (
	"context"

	"github.com/aldisaputra17/shoping-cart/entities"
	"github.com/stretchr/testify/mock"
)

type ProductRepository struct {
	mock.Mock
}

func (_m *ProductRepository) Deleted(ctx context.Context, prod entities.Products) {

}

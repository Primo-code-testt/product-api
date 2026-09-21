package repository

import (
	"context"

	"github.com/mahasachan/kkp-pre-test/internal/domain"
)

type ProductRepository interface {
	Create(ctx context.Context, input domain.CreateProduct) (domain.Product, error)
	Patch(ctx context.Context, id int64, patch domain.ProductPatch) (domain.Product, error)
}

package usecase

import (
	"context"

	"github.com/mahasachan/kkp-pre-test/internal/domain"
	"github.com/mahasachan/kkp-pre-test/internal/repository"
)

type ProductService struct {
	repository repository.ProductRepository
}

func NewProductService(repository repository.ProductRepository) *ProductService {
	return &ProductService{repository: repository}
}

func (s *ProductService) Create(ctx context.Context, input domain.CreateProduct) (domain.Product, error) {
	if err := input.Validate(); err != nil {
		return domain.Product{}, err
	}
	return s.repository.Create(ctx, input)
}

func (s *ProductService) Patch(ctx context.Context, id int64, patch domain.ProductPatch) (domain.Product, error) {
	if id <= 0 {
		return domain.Product{}, domain.ErrValidation
	}
	if err := patch.Validate(); err != nil {
		return domain.Product{}, err
	}
	return s.repository.Patch(ctx, id, patch)
}

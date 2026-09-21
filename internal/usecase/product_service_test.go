package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/mahasachan/kkp-pre-test/internal/domain"
)

type repositoryStub struct {
	created domain.CreateProduct
	patched domain.ProductPatch
}

func (r *repositoryStub) Create(_ context.Context, input domain.CreateProduct) (domain.Product, error) {
	r.created = input
	return domain.Product{ID: 1, Name: input.Name, Price: input.Price}, nil
}

func (r *repositoryStub) Patch(_ context.Context, _ int64, patch domain.ProductPatch) (domain.Product, error) {
	r.patched = patch
	return domain.Product{ID: 1}, nil
}

func TestProductServiceCreate(t *testing.T) {
	repository := &repositoryStub{}
	service := NewProductService(repository)

	product, err := service.Create(context.Background(), domain.CreateProduct{Name: "Keyboard", Price: 1200})
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}
	if product.ID != 1 || repository.created.Name != "Keyboard" {
		t.Fatalf("unexpected product or repository input: %+v %+v", product, repository.created)
	}
}

func TestProductServiceRejectsInvalidInputBeforeRepository(t *testing.T) {
	repository := &repositoryStub{}
	service := NewProductService(repository)

	_, err := service.Create(context.Background(), domain.CreateProduct{Name: "", Price: 1200})
	if !errors.Is(err, domain.ErrValidation) {
		t.Fatalf("expected validation error, got %v", err)
	}
	if repository.created.Name != "" || repository.created.Price != 0 {
		t.Fatal("repository should not be called for invalid input")
	}
}

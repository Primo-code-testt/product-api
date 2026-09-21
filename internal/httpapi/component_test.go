package httpapi

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/mahasachan/kkp-pre-test/internal/domain"
	"github.com/mahasachan/kkp-pre-test/internal/usecase"
)

type memoryRepository struct {
	mu      sync.Mutex
	nextID  int64
	product map[int64]domain.Product
}

func newMemoryRepository() *memoryRepository {
	return &memoryRepository{nextID: 1, product: make(map[int64]domain.Product)}
}

func (r *memoryRepository) Create(_ context.Context, input domain.CreateProduct) (domain.Product, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	now := time.Now().UTC()
	product := domain.Product{
		ID: r.nextID, Name: input.Name, Description: input.Description,
		Price: input.Price, SalePrice: input.SalePrice, CreatedAt: now, UpdatedAt: now,
	}
	r.nextID++
	r.product[product.ID] = product
	return product, nil
}

func (r *memoryRepository) Patch(_ context.Context, id int64, patch domain.ProductPatch) (domain.Product, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	product, exists := r.product[id]
	if !exists {
		return domain.Product{}, domain.ErrNotFound
	}
	if patch.Name.Set {
		product.Name = *patch.Name.Value
	}
	if patch.Description.Set {
		product.Description = patch.Description.Value
	}
	if patch.Price.Set {
		product.Price = *patch.Price.Value
	}
	if patch.SalePrice.Set {
		product.SalePrice = patch.SalePrice.Value
	}
	product.UpdatedAt = time.Now().UTC()
	r.product[id] = product
	return product, nil
}

func TestProductComponentCreateThenPatchOnlyProvidedFields(t *testing.T) {
	repository := newMemoryRepository()
	router := NewRouter(NewHandler(usecase.NewProductService(repository)))

	create := httptest.NewRequest(http.MethodPost, "/product", bytes.NewBufferString(
		`{"name":"Keyboard","description":"Mechanical","price":1200,"sale_price":999}`,
	))
	create.Header.Set("Content-Type", "application/json")
	createResponse := httptest.NewRecorder()
	router.ServeHTTP(createResponse, create)
	if createResponse.Code != http.StatusCreated {
		t.Fatalf("create status = %d, body = %s", createResponse.Code, createResponse.Body.String())
	}

	patch := httptest.NewRequest(http.MethodPatch, "/product/1", bytes.NewBufferString(
		`{"description":null,"price":1100}`,
	))
	patch.Header.Set("Content-Type", "application/json")
	patchResponse := httptest.NewRecorder()
	router.ServeHTTP(patchResponse, patch)
	if patchResponse.Code != http.StatusOK {
		t.Fatalf("patch status = %d, body = %s", patchResponse.Code, patchResponse.Body.String())
	}

	var result struct {
		Successful bool           `json:"successful"`
		Data       domain.Product `json:"data"`
	}
	if err := json.Unmarshal(patchResponse.Body.Bytes(), &result); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if !result.Successful || result.Data.Name != "Keyboard" || result.Data.Price != 1100 {
		t.Fatalf("unexpected patched product: %+v", result.Data)
	}
	if result.Data.Description != nil {
		t.Fatalf("description should be null: %+v", result.Data.Description)
	}
	if result.Data.SalePrice == nil || *result.Data.SalePrice != 999 {
		t.Fatalf("omitted sale_price should be unchanged: %+v", result.Data.SalePrice)
	}
}

func TestPatchRejectsNullForNonNullableField(t *testing.T) {
	repository := newMemoryRepository()
	router := NewRouter(NewHandler(usecase.NewProductService(repository)))

	request := httptest.NewRequest(http.MethodPatch, "/product/1", bytes.NewBufferString(`{"price":null}`))
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	if response.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, body = %s", response.Code, response.Body.String())
	}
}

func TestCreateRequiresPrice(t *testing.T) {
	repository := newMemoryRepository()
	router := NewRouter(NewHandler(usecase.NewProductService(repository)))

	request := httptest.NewRequest(http.MethodPost, "/product", bytes.NewBufferString(`{"name":"Keyboard"}`))
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	if response.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, body = %s", response.Code, response.Body.String())
	}
}

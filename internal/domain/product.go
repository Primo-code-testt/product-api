package domain

import (
	"errors"
	"strings"
	"time"
)

var (
	ErrNotFound   = errors.New("product not found")
	ErrValidation = errors.New("validation failed")
)

type Product struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	Price       float64   `json:"price"`
	SalePrice   *float64  `json:"sale_price"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateProduct struct {
	Name        string
	Description *string
	Price       float64
	SalePrice   *float64
}

// Field represents PATCH semantics: Set=false means absent; Set=true and
// Value=nil means explicit null.
type Field[T any] struct {
	Set   bool
	Value *T
}

type ProductPatch struct {
	Name        Field[string]
	Description Field[string]
	Price       Field[float64]
	SalePrice   Field[float64]
}

func (c CreateProduct) Validate() error {
	if strings.TrimSpace(c.Name) == "" || c.Price < 0 {
		return ErrValidation
	}
	if c.SalePrice != nil && *c.SalePrice < 0 {
		return ErrValidation
	}
	return nil
}

func (p ProductPatch) Validate() error {
	if !p.Name.Set && !p.Description.Set && !p.Price.Set && !p.SalePrice.Set {
		return ErrValidation
	}
	if p.Name.Set && (p.Name.Value == nil || strings.TrimSpace(*p.Name.Value) == "") {
		return ErrValidation
	}
	if p.Price.Set && (p.Price.Value == nil || *p.Price.Value < 0) {
		return ErrValidation
	}
	if p.SalePrice.Set && p.SalePrice.Value != nil && *p.SalePrice.Value < 0 {
		return ErrValidation
	}
	return nil
}

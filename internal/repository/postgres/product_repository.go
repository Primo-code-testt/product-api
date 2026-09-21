package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mahasachan/kkp-pre-test/internal/domain"
)

type ProductRepository struct {
	pool *pgxpool.Pool
}

func NewProductRepository(pool *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{pool: pool}
}

func (r *ProductRepository) Create(ctx context.Context, input domain.CreateProduct) (domain.Product, error) {
	const query = `
		INSERT INTO products (name, description, price, sale_price)
		VALUES ($1, $2, $3, $4)
		RETURNING id, name, description, price, sale_price, created_at, updated_at`

	return scanProduct(r.pool.QueryRow(ctx, query,
		input.Name, input.Description, input.Price, input.SalePrice,
	))
}

func (r *ProductRepository) Patch(ctx context.Context, id int64, patch domain.ProductPatch) (domain.Product, error) {
	const query = `
		UPDATE products SET
			name = CASE WHEN $2 THEN $3 ELSE name END,
			description = CASE WHEN $4 THEN $5 ELSE description END,
			price = CASE WHEN $6 THEN $7 ELSE price END,
			sale_price = CASE WHEN $8 THEN $9 ELSE sale_price END,
			updated_at = NOW()
		WHERE id = $1
		RETURNING id, name, description, price, sale_price, created_at, updated_at`

	product, err := scanProduct(r.pool.QueryRow(ctx, query,
		id,
		patch.Name.Set, patch.Name.Value,
		patch.Description.Set, patch.Description.Value,
		patch.Price.Set, patch.Price.Value,
		patch.SalePrice.Set, patch.SalePrice.Value,
	))
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Product{}, domain.ErrNotFound
	}
	return product, err
}

type rowScanner interface {
	Scan(dest ...any) error
}

func scanProduct(row rowScanner) (domain.Product, error) {
	var product domain.Product
	err := row.Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.SalePrice,
		&product.CreatedAt,
		&product.UpdatedAt,
	)
	return product, err
}

//go:build integration

package postgres

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mahasachan/kkp-pre-test/internal/domain"
)

func TestProductRepositoryCreateAndPatch(t *testing.T) {
	databaseURL := os.Getenv("TEST_DATABASE_URL")
	if databaseURL == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()

	_, err = pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS products (
			id BIGSERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			description TEXT NULL,
			price DOUBLE PRECISION NOT NULL CHECK (price >= 0),
			sale_price DOUBLE PRECISION NULL CHECK (sale_price IS NULL OR sale_price >= 0),
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _, _ = pool.Exec(ctx, "TRUNCATE products RESTART IDENTITY") })

	repository := NewProductRepository(pool)
	description := "Mechanical"
	product, err := repository.Create(ctx, domain.CreateProduct{
		Name: "Keyboard", Description: &description, Price: 1200,
	})
	if err != nil {
		t.Fatalf("create product: %v", err)
	}

	newPrice := 1100.0
	product, err = repository.Patch(ctx, product.ID, domain.ProductPatch{
		Description: domain.Field[string]{Set: true, Value: nil},
		Price:       domain.Field[float64]{Set: true, Value: &newPrice},
	})
	if err != nil {
		t.Fatalf("patch product: %v", err)
	}
	if product.Description != nil || product.Price != 1100 {
		t.Fatalf("unexpected product: %+v", product)
	}
}

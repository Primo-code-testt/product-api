package domain

import (
	"errors"
	"testing"
)

func TestCreateProductValidate(t *testing.T) {
	tests := []struct {
		name  string
		input CreateProduct
		valid bool
	}{
		{name: "valid", input: CreateProduct{Name: "Keyboard", Price: 1200}, valid: true},
		{name: "blank name", input: CreateProduct{Name: " ", Price: 1200}},
		{name: "negative price", input: CreateProduct{Name: "Keyboard", Price: -1}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.input.Validate()
			if test.valid && err != nil {
				t.Fatalf("expected valid input, got %v", err)
			}
			if !test.valid && !errors.Is(err, ErrValidation) {
				t.Fatalf("expected validation error, got %v", err)
			}
		})
	}
}

func TestProductPatchValidateNullableFields(t *testing.T) {
	patch := ProductPatch{Description: Field[string]{Set: true, Value: nil}}
	if err := patch.Validate(); err != nil {
		t.Fatalf("explicit null description should be valid: %v", err)
	}

	if err := (ProductPatch{}).Validate(); !errors.Is(err, ErrValidation) {
		t.Fatalf("empty patch should fail validation, got %v", err)
	}
}

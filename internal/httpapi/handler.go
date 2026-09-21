package httpapi

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/mahasachan/kkp-pre-test/internal/domain"
	"github.com/mahasachan/kkp-pre-test/internal/usecase"
)

type Handler struct {
	service *usecase.ProductService
}

type response struct {
	Successful bool   `json:"successful"`
	ErrorCode  string `json:"error_code"`
	Data       any    `json:"data,omitempty"`
}

func NewHandler(service *usecase.ProductService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Name        string   `json:"name"`
		Description *string  `json:"description"`
		Price       *float64 `json:"price"`
		SalePrice   *float64 `json:"sale_price"`
	}
	if err := decodeJSON(w, r, &request); err != nil || request.Price == nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST")
		return
	}

	product, err := h.service.Create(r.Context(), domain.CreateProduct{
		Name: request.Name, Description: request.Description,
		Price: *request.Price, SalePrice: request.SalePrice,
	})
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, response{Successful: true, ErrorCode: "", Data: product})
}

func (h *Handler) PatchProduct(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST")
		return
	}

	var request struct {
		Name        optional[string]  `json:"name"`
		Description optional[string]  `json:"description"`
		Price       optional[float64] `json:"price"`
		SalePrice   optional[float64] `json:"sale_price"`
	}
	if err := decodeJSON(w, r, &request); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST")
		return
	}

	product, err := h.service.Patch(r.Context(), id, domain.ProductPatch{
		Name:        domain.Field[string]{Set: request.Name.set, Value: request.Name.value},
		Description: domain.Field[string]{Set: request.Description.set, Value: request.Description.value},
		Price:       domain.Field[float64]{Set: request.Price.set, Value: request.Price.value},
		SalePrice:   domain.Field[float64]{Set: request.SalePrice.set, Value: request.SalePrice.value},
	})
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, response{Successful: true, ErrorCode: "", Data: product})
}

func decodeJSON(w http.ResponseWriter, r *http.Request, destination any) error {
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(destination); err != nil {
		return err
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return errors.New("request body must contain exactly one JSON value")
	}
	return nil
}

func writeDomainError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrValidation):
		writeError(w, http.StatusBadRequest, "VALIDATION_ERROR")
	case errors.Is(err, domain.ErrNotFound):
		writeError(w, http.StatusNotFound, "PRODUCT_NOT_FOUND")
	default:
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR")
	}
}

func writeError(w http.ResponseWriter, status int, code string) {
	writeJSON(w, status, response{Successful: false, ErrorCode: code})
}

func writeJSON(w http.ResponseWriter, status int, body response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

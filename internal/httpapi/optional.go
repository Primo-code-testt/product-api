package httpapi

import (
	"bytes"
	"encoding/json"
)

type optional[T any] struct {
	set   bool
	value *T
}

func (o *optional[T]) UnmarshalJSON(data []byte) error {
	o.set = true
	if bytes.Equal(bytes.TrimSpace(data), []byte("null")) {
		o.value = nil
		return nil
	}
	var value T
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	o.value = &value
	return nil
}

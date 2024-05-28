package json

import (
	"encoding/json"
	"fmt"
	"log"
	"service/service/internal/lib/model"

	"github.com/go-playground/validator/v10"
)

// ValidateData валидирует данные, поступившие из NATS-Streaming и переданные в []byte.
// data в функции десериализируется в model.Order.
func ValidateData(data []byte) (*model.Order, error) {
	const op = "json.json.GetOrderUIDFromJson"

	var order model.Order

	err := json.Unmarshal(data, &order)
	if err != nil {
		return &order, fmt.Errorf("%s: failed unmarshal json: %w", op, err)
	}

	validate := validator.New()
	if err := validate.Struct(order); err != nil {
		errs := err.(validator.ValidationErrors)
		for _, fieldErr := range errs {
			log.Printf("op:[%v] field:[%s] %s\n", op, fieldErr.Field(), fieldErr.Tag())

		}
		return &order, fmt.Errorf("%s: failed validate json: %w", op, err)
	}

	return &order, nil
}

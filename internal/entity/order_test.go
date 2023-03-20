package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func PrototypeOrder() Order {
	return Order{
		ID:         "any_id",
		Price:      float64(190),
		Tax:        float64(2),
		FinalPrice: float64(380),
	}
}

func TestOrder_Validate(t *testing.T) {
	orderMissingID := PrototypeOrder()
	orderMissingID.ID = ""

	orderInvalidPrice := PrototypeOrder()
	orderInvalidPrice.Price = -10

	orderInvalidTax := PrototypeOrder()
	orderInvalidTax.Tax = -2

	testCases := []struct {
		Name          string
		Order         Order
		ExpectedError bool
	}{
		{
			Name:          "All valid params",
			Order:         PrototypeOrder(),
			ExpectedError: false,
		}, {
			Name:          "Error - Missing ID",
			Order:         orderMissingID,
			ExpectedError: true,
		}, {
			Name:          "Error - Invalid Price",
			Order:         orderInvalidPrice,
			ExpectedError: true,
		}, {
			Name:          "Error - Invalid Tax",
			Order:         orderInvalidPrice,
			ExpectedError: true,
		},
	}

	for _, tc := range testCases {
		err := tc.Order.Validate()
		if tc.ExpectedError {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
		}
	}
}

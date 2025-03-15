package entities

import productValueobjects "github.com/baothaihcmut/Ecommerce-go/products/internal/core/command/domain/aggregates/products/value_objects"

type VariationValue struct {
	VariationId productValueobjects.VariationId
	Value       string
}

func NewVariationValue(variationId productValueobjects.VariationId, value string) *VariationValue {
	return &VariationValue{
		VariationId: variationId,
		Value:       value,
	}
}
func (v *VariationValue) IsEqual(o *VariationValue) bool {
	return v.VariationId.IsEqual(o.VariationId)
}

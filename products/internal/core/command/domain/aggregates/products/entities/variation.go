package entities

import valueobjects "github.com/baothaihcmut/Ecommerce-Go/products/internal/core/command/domain/aggregates/products/value_objects"

type Variation struct {
	Id valueobjects.VariationId
}

func NewVariation(id valueobjects.VariationId) *Variation {
	return &Variation{
		Id: id,
	}
}
func (v *Variation) IsEqual(o *Variation) bool {
	return v.Id.IsEqual(o.Id)
}

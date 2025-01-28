package productitems

import (
	"github.com/baothaihcmut/Ecommerce-Go/products/internal/core/command/domain/aggregates/product_items/entities"
	valueobjects "github.com/baothaihcmut/Ecommerce-Go/products/internal/core/command/domain/aggregates/product_items/value_objects"
	productValueobjects "github.com/baothaihcmut/Ecommerce-Go/products/internal/core/command/domain/aggregates/products/value_objects"
	"github.com/baothaihcmut/Ecommerce-Go/products/internal/core/command/domain/exceptions"
)

type VariationValueArg struct {
	VariationId productValueobjects.VariationId
	Value       string
}

type ProductItem struct {
	Id              valueobjects.ProductItemId
	Sku             valueobjects.SKU
	Price           valueobjects.ProductPrice
	Quantity        valueobjects.ProductQuantity
	ProductId       productValueobjects.ProductId
	VariationValues []*entities.VariationValue
}

func NewProductItem(
	id valueobjects.ProductItemId,
	sku valueobjects.SKU,
	price valueobjects.ProductPrice,
	quantity valueobjects.ProductQuantity,
	productId productValueobjects.ProductId,
	variationArgs []VariationValueArg,
) (*ProductItem, error) {
	//check if product item have two same variationvalues of 1 variation
	mapVariation := make(map[string]struct{}, len(variationArgs))
	variationValues := make([]*entities.VariationValue, len(variationArgs))
	for idx, variation := range variationArgs {
		key := variation.Value + variation.VariationId.Name
		if _, exist := mapVariation[key]; exist {
			return nil, exceptions.ErrDuplicateVariation
		} else {
			mapVariation[key] = struct{}{}
		}
		variationValues[idx] = entities.NewVariationValue(variation.VariationId, variation.Value)
	}
	return &ProductItem{
		Id:              id,
		Sku:             sku,
		Price:           price,
		Quantity:        quantity,
		ProductId:       productId,
		VariationValues: variationValues,
	}, nil
}
func (p *ProductItem) AddVariationValues(variationArgs []VariationValueArg) error {
	//check for duplicate variation value
	for _, variationArg := range variationArgs {
		variationValue := entities.NewVariationValue(variationArg.VariationId, variationArg.Value)
		for _, productVariation := range p.VariationValues {
			if productVariation.VariationId.IsEqual(variationValue.VariationId) && productVariation.Value == variationValue.Value {
				return exceptions.ErrDuplicateVariation
			}
		}
		p.VariationValues = append(p.VariationValues, variationValue)
	}
	return nil
}

func (p *ProductItem) IncreProductItemQuantity(quantity valueobjects.ProductQuantity) {
	p.Quantity = p.Quantity.IncreQuantity(quantity)
}

func (p *ProductItem) DecreProductItemQuantity(quantity valueobjects.ProductQuantity) error {
	newQuantity, err := p.Quantity.DecreQuantity(quantity)
	if err != nil {
		return err
	}
	p.Quantity = newQuantity
	return nil
}

package valueobject

type Address struct {
	Street   string
	Town     string
	City     string
	Province string
}

func NewAddress(
	street string,
	town string,
	city string,
	province string,
) *Address {
	return &Address{
		Town:     town,
		City:     city,
		Street:   street,
		Province: province,
	}
}
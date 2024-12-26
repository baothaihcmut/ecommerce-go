package valueobject

type Address struct {
	Street   string
	Town     string
	City     string
	Province string
	Country  string
}

func NewAddress(
	street string,
	town string,
	city string,
	province string,
	country string,
) *Address {
	return &Address{
		Town:     town,
		City:     city,
		Street:   street,
		Province: province,
		Country:  country,
	}
}

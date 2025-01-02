package valueobject

type Address struct {
	Priority int
	Street   string
	Town     string
	City     string
	Province string
}

func NewAddress(
	priority int,
	street string,
	town string,
	city string,
	province string,
) *Address {
	return &Address{
		Priority: priority,
		Town:     town,
		City:     city,
		Street:   street,
		Province: province,
	}
}

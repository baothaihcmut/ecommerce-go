package commands

type SignUpAddress struct {
	Priority int
	Street   string
	Town     string
	City     string
	Province string
}

type SignUpCommand struct {
	Email       string
	Password    string
	PhoneNumber string
	Addresses   []SignUpAddress
	FirstName   string
	LastName    string
}

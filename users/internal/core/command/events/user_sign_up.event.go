package events

type UserSignUpEvent struct {
	Email          string `json:"email"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	PhoneNumber    string `json:"phone_number"`
	VericationCode string `json:"verification_code"`
}

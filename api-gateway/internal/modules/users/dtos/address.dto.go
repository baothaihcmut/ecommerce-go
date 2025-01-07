package dtos

type Address struct {
	Priority int    `json:"priority"`
	Town     string `json:"town"`
	Street   string `json:"street"`
	City     string `json:"city"`
	Province string `json:"province"`
}

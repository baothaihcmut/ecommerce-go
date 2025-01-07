package enums

type TokenType string

const (
	REFRESH_TOKEN TokenType = "REFRESH_TOKEN"
	ACCESS_TOKEN  TokenType = "ACCESS_TOKEN"
	OAUTH_TOKEN   TokenType = "OAUTH_TOKEN"
)

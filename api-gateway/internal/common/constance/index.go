package constance

type PayloadKey string

type ErrorKey string

type AccessTokenKey string
type RefreshTokenKey string

type UserContextKey string
type TokenContextKey string

const PayloadContext PayloadKey = "payload"

const ErrorContext ErrorKey = "error"

const AccessToken AccessTokenKey = "access_token"

const RefreshToken RefreshTokenKey = "refresh_token"

const UserContext UserContextKey = "user_context"

const TokenContext TokenContextKey = "token_context"

package constance

type PayloadKey string

type ErrorKey string

type AccessTokenKey string
type RefresTokenKey string

type UserContextKey string

const PayloadContext PayloadKey = "payload"

const ErrorContext ErrorKey = "error"

const AccessToken AccessTokenKey = "access_token"

const RefreshToken RefresTokenKey = "refresh_token"

const UserContext UserContextKey = "user_context"

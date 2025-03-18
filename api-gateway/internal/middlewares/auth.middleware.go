package middlewares

import (
	"encoding/json"
	"errors"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/baothaihcmut/Ecommerce-go/api-gateway/internal/config"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrTokenExpire = errors.New("token expire")
)

type accessTokenClaims struct {
	*jwt.RegisteredClaims
	UserId 			  uuid.UUID `json:"user_id"`
	IsShopOwnerActive bool `json:"is_shop_owner_active"`
}

func validateToken(jwtConfig *config.JwtConfig,tokenString string) (accessTokenClaims,error) {
	var claims accessTokenClaims

	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtConfig.AccessToken.Secret, nil
	})

	if err != nil || !token.Valid {
		return claims, ErrInvalidToken
	}

	// Check if the token is expired
	if claims.ExpiresAt.Time.Before(time.Now()) {
		return claims, ErrTokenExpire
	}
	return claims, nil
}


func AuthMiddleware(jwtConfig *config.JwtConfig, webconfig *config.WebConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				uri := strings.TrimPrefix(r.RequestURI,webconfig.Prefix)
				if slices.Contains(webconfig.Public, uri) {
					next.ServeHTTP(w,r)
					return
				}
				authHeader := r.Header.Get("Authorization")
				
				claims,err:= validateToken(jwtConfig,strings.TrimPrefix(authHeader,"Bearer "))
				if err != nil{
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusUnauthorized)
					json.NewEncoder(w).Encode(map[string]any{
						"success": false,
						"message": err.Error(),
					})
					return
				}
				isShopOwnerActive:= "false"
				if claims.IsShopOwnerActive {
				isShopOwnerActive ="true"
				} 
				md:= metadata.Pairs(
					"user_id",claims.UserId.String(),
					"is_shop_owner_active", isShopOwnerActive,
				)
				ctx := metadata.NewIncomingContext(r.Context(),md)
				r = r.WithContext(ctx)
				next.ServeHTTP(w,r)
			},
		)
	}
}


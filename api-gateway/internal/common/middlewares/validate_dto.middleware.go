package middlewares

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/common/constance"
	"github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/common/utils"
	"github.com/go-playground/validator/v10"
)

func ValidateDTOMiddleware[T any]() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var dest T
			//decode json
			if err := json.NewDecoder(r.Body).Decode(&dest); err != nil {
				utils.WriteResponseErr(w, http.StatusBadRequest, []string{"Invalid json format"})
				return
			}
			//validate json
			if err := validator.New().Struct(dest); err != nil {
				validationErrors := err.(validator.ValidationErrors)

				// Collect error messages
				errorMessages := make([]string, len(validationErrors))
				for idx, fieldError := range validationErrors {
					errorMessages[idx] = fmt.Sprintf("Error with field %s: %s", fieldError.Field(), fieldError.Tag())
				}
				//if err response
				utils.WriteResponseErr(w, http.StatusBadRequest, errorMessages)
				return
			}
			ctx := context.WithValue(r.Context(), constance.PayloadContext, &dest)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

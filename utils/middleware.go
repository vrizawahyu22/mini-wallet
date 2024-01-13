package utils

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
)

func SetupMiddleware(route *chi.Mux, config *BaseConfig) {
	route.Use(Recovery)
}

func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err != nil {
				var statusCode int

				appErr, isAppErr := err.(AppError)
				validationErr, isValidationErr := err.(ValidationErrors)

				if isAppErr {
					statusCode = appErr.StatusCode
					messages := strings.Split(appErr.Message, "|")
					errorMsgs := messages[1]

					if appErr.StatusCode >= 500 {
						message := messages[0]
						LogError(fmt.Sprintf("APP ERROR (PANIC) %s", message))
					}

					if appErr.StatusCode >= 400 {
						LogWarning(fmt.Sprintf("APP ERROR (PANIC) %s", messages[0]))
					}

					GenerateErrorResp(w, errorMsgs, statusCode)
				} else if isValidationErr {
					errorMsgs := make(map[string]interface{})
					LogWarning(fmt.Sprintf("VALIDATION ERROR (PANIC) %v", validationErr))

					for _, err := range validationErr.Errors {
						errorMsgs[err.Key] = []string{err.Message}
					}

					statusCode = validationErr.StatusCode

					GenerateErrorResp(w, errorMsgs, statusCode)
				} else {
					LogError(fmt.Sprintf("UNKNOWN ERROR (PANIC) %v", err))
					errorMsgs := []map[string]interface{}{
						{"message": "internal server error"},
					}

					statusCode = 500

					GenerateErrorResp(w, errorMsgs, statusCode)
				}
			}
		}()
		next.ServeHTTP(w, r)
	})
}

type AuthPayload struct {
	UserID string `json:"userId"`
	Token  string `json:"token"`
}
type key string

const UserSession key = "user-session"

func AppendRequestCtx(r *http.Request, ctxKey key, input interface{}) context.Context {
	return context.WithValue(r.Context(), ctxKey, input)
}

func CheckIsAuthenticated(handler func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")

		if header == "" || !strings.Contains(header, "Token ") {
			PanicIfError(CustomError("unauthorized", 401))
		}
		authToken := strings.Split(header, " ")[1]

		var authPayload AuthPayload
		authPayload.Token = authToken
		authCtx := AppendRequestCtx(r, UserSession, &authPayload)

		handler(w, r.WithContext(authCtx))
	}
}

func GetRequestCtx(ctx context.Context, ctxKey key) *AuthPayload {
	ctxVal := ctx.Value(ctxKey)
	if ctxVal != nil {
		return ctxVal.(*AuthPayload)
	}

	return &AuthPayload{}
}

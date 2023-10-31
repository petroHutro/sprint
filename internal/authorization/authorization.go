package authorization

import (
	"fmt"
	"net/http"
	"sprint/internal/logger"
	"sprint/internal/utils"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID string
}

func buildJWTString(secretKey string, tokenEXP time.Duration) (string, error) {
	id := utils.GenerateString()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenEXP)),
		},
		UserID: id,
	})

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("cannot get token: %v", err)
	}

	return tokenString, nil
}

func getUserID(secretKey, tokenString string) (string, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(secretKey), nil
		})
	if err != nil {
		return "", fmt.Errorf("cannot pars: %v", err)
	}

	if !token.Valid {
		return "", fmt.Errorf("token is not valid: %v", err)
	}

	return claims.UserID, nil
}

func setAuthorization(w *http.ResponseWriter, secretKey string, tokenEXP time.Duration) *http.Cookie {
	token, _ := buildJWTString(secretKey, tokenEXP)
	cookie := http.Cookie{Name: "Authorization", Value: token}
	http.SetCookie(*w, &cookie)
	return &cookie
}

func AuthorizationMiddleware(secretKey string, tokenEXP time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("Authorization")
			if err != nil {
				logger.Error("cookies do not contain a token: %v", err)
				cookie = setAuthorization(&w, secretKey, tokenEXP)
			}
			id, err := getUserID(secretKey, cookie.Value)
			if err != nil {
				logger.Error("token does not pass validation")
				setAuthorization(&w, secretKey, tokenEXP)
			}
			r.Header.Set("User_id", id)
			next.ServeHTTP(w, r)
		})
	}
}

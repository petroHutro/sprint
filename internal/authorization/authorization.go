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

const TokenEXP = time.Hour * 3
const SecretKey = "supersecretkey"

func buildJWTString() (string, error) {
	id := utils.GenerateString()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenEXP)),
		},
		UserID: id,
	})

	tokenString, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		return "", fmt.Errorf("cannot get token: %v", err)
	}

	return tokenString, nil
}

func getUserID(tokenString string) (string, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(SecretKey), nil
		})
	if err != nil {
		return "", fmt.Errorf("cannot pars: %v", err)
	}

	if !token.Valid {
		return "", fmt.Errorf("token is not valid: %v", err)
	}

	return claims.UserID, nil
}

func setAuthorization(w *http.ResponseWriter) *http.Cookie {
	token, _ := buildJWTString()
	cookie := http.Cookie{Name: "Authorization", Value: token}
	http.SetCookie(*w, &cookie)
	return &cookie
}

func AuthorizationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("Authorization")
		if err != nil {
			logger.Error("cookies do not contain a token: %v", err)
			cookie = setAuthorization(&w)
		}
		id, err := getUserID(cookie.Value)
		if err != nil {
			logger.Error("token does not pass validation")
			setAuthorization(&w)
		}
		r.Header.Set("User_id", id)
		next.ServeHTTP(w, r)
	})
}

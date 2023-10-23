package authorization

import (
	"fmt"
	"math/rand"
	"net/http"
	"sprint/internal/logger"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID int
}

const TokenEXP = time.Hour * 3
const SecretKey = "supersecretkey"

func buildJWTString() (string, error) {
	source := rand.NewSource(time.Now().UnixNano())
	generator := rand.New(source)
	id := generator.Intn(2000000)
	// id := uuid.New() // исправить проверка на наличие
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenEXP)),
		},
		UserID: id,
	})

	tokenString, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func getUserID(tokenString string) (int, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(SecretKey), nil
		})
	if err != nil {
		return -1, fmt.Errorf("cannot pars: %v", err)
	}

	if !token.Valid {
		return -1, fmt.Errorf("token is not valid: %v", err)
	}

	return claims.UserID, nil
}

func setAuthorization(w *http.ResponseWriter) {
	token, _ := buildJWTString()
	cookie := http.Cookie{Name: "Authorization", Value: token}
	http.SetCookie(*w, &cookie)
}

func AuthorizationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("Authorization")
		if err != nil {
			logger.Error("cookies do not contain a token: %v", err)
			token, _ := buildJWTString()
			cookie = &http.Cookie{Name: "Authorization", Value: token}
			http.SetCookie(w, cookie)
			// setAuthorization(&w)
			// w.WriteHeader(http.StatusUnauthorized)
			// return
		}
		id, err := getUserID(cookie.Value)
		if err != nil {
			logger.Error("token does not pass validation")
			setAuthorization(&w)
			// w.WriteHeader(http.StatusUnauthorized)
			// return
		}
		r.Header.Set("User_id", strconv.Itoa(id))
		next.ServeHTTP(w, r)
	})
}

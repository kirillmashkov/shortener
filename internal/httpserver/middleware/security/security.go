package security

import (
	"context"
	"fmt"
	"net/http"

	"math/rand"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/kirillmashkov/shortener.git/internal/app"
	"go.uber.org/zap"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID int
}

type UserIDType string

const tokenExp = time.Hour * 3
const secretKey = "supersecretkey"

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "debug") {
			next.ServeHTTP(w, r)
			return
		}

		cookie, err := r.Cookie("token")

		if err != nil {
			cookie = nil
		}

		jwtToken, userID, err, newToken := getJWT(cookie)
		if err != nil {
			app.Log.Error("Error get token", zap.Error(err))
			http.Error(w, "Something went wrong", http.StatusBadRequest)
			return
		}

		u := UserIDType("userID")
		c := context.WithValue(r.Context(), u, userID)

		if newToken {
			resCookie := http.Cookie{Name: "token", Value: jwtToken}
			http.SetCookie(w, &resCookie)
		}

		next.ServeHTTP(w, r.WithContext(c))
	})
}

func getJWT(cookie *http.Cookie) (string, int, error, bool) {
	if cookie == nil {
		tokenString, userID, err := buildJWTString()
		return tokenString, userID, err, true
	}

	checkJWT, userID := сheckJWT(cookie)

	if checkJWT {
		return cookie.Value, userID, nil, false
	}

	tokenString, userID, err := buildJWTString()
	return tokenString, userID, err, true
}

func buildJWTString() (string, int, error) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	userID := r.Int()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenExp)),
		},
		UserID: userID,
	})

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", 0, err
	}

	return tokenString, userID, nil
}

func сheckJWT(cookie *http.Cookie) (bool, int) {
	if cookie == nil {
		app.Log.Warn("Token is empty")
		return true, -1
	}

	claims := &Claims{UserID: -1}
	token, err := jwt.ParseWithClaims(cookie.Value, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(secretKey), nil
		})

	if err != nil {
		app.Log.Warn("Can't parse token")
		return false, 0
	}

	if !token.Valid {
		app.Log.Warn("Token is not valid")
		return false, 0
	}

	if claims.UserID == -1 {
		app.Log.Warn("Token doesn't contain UserID")
		return false, 0
	}

	app.Log.Info("Token is valid", zap.Int("UserID", claims.UserID))
	return true, claims.UserID
}

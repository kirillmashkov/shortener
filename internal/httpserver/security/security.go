package security

import (
	"fmt"
	"net/http"

	"math/rand"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/kirillmashkov/shortener.git/internal/app"
	"go.uber.org/zap"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID int
}

const tokenExp = time.Hour * 3
const secretKey = "supersecretkey"

func GetJWT(cookie *http.Cookie) (string, error) {
	checkJWT := CheckJWT(cookie)

	if checkJWT {
		return cookie.Value, nil
	}

	return buildJWTString()
}

func buildJWTString() (string, error) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenExp)),
		},
		UserID: r.Int(),
	})

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func CheckJWT(cookie *http.Cookie) bool {
	if cookie == nil {
		app.Log.Warn("Token is empty")
		return true
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
		return false
	}

	if !token.Valid {
		app.Log.Warn("Token is not valid")
		return false
	}

	if claims.UserID == -1 {
		app.Log.Warn("Token doesn't contain UserID")
		return false
	}

	app.Log.Info("Token is valid", zap.Int("UserID", claims.UserID))
	return true
}

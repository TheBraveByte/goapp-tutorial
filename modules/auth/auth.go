package auth

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/akinbyte/go-app/modules/config"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"os"
	"time"
)

var app config.GoAppTools

type GoAppClaims struct {
	jwt.RegisteredClaims
	Email string
	ID    primitive.ObjectID
}

var secretKey = os.Getenv("GOAPP_KEY")

func Generate(email string, id primitive.ObjectID) (string, string, error) {
	goappClaims := GoAppClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:   "goAppUser",
			IssuedAt: jwt.NewNumericDate(time.Now()),
			ExpiresAt: &jwt.NumericDate{
				Time: time.Now().Add(24 * time.Hour),
			},
		},
		Email: email,
		ID:    id,
	}

	newGoAppClaims := &jwt.RegisteredClaims{
		Issuer:   "goAppUser",
		IssuedAt: jwt.NewNumericDate(time.Now()),
		ExpiresAt: &jwt.NumericDate{
			Time: time.Now().Add(48 * time.Hour),
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, goappClaims).SignedString([]byte(secretKey))
	if err != nil {
		return "", "", err
	}
	newToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, newGoAppClaims).SignedString([]byte(secretKey))
	if err != nil {
		return "", "", err
	}
	return token, newToken, nil
}

func Parse(tokenString string) (*GoAppClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &GoAppClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		app.ErrorLogger.Fatalf("error while parsing token with it claims %v", err)
	}
	claims, ok := token.Claims.(*GoAppClaims)
	if !ok {
		app.ErrorLogger.Fatalf("error %v controller not authorized access", http.StatusUnauthorized)
	}
	if err := claims.Valid(); err != nil {
		app.ErrorLogger.Fatalf("error %v %s", http.StatusUnauthorized, err)
	}
	return claims, nil
}

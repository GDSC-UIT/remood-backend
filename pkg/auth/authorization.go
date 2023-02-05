package auth

import (
	"strings"
	"errors"
	"crypto/rand"
	"log"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Claims struct {
	ID primitive.ObjectID `json:"_id" bson:"_id"`
	jwt.StandardClaims
}

var (
	JWTKey []byte // todo: init
)

func GenerateJWTKey() {
	JWTKey = make([]byte, 32)
	rand.Read(JWTKey)
	log.Println(JWTKey)
}

func GetTokenString(ctx *gin.Context) string {
	var tokenString string
	authHeader := ctx.Request.Header.Get("Authorization")
	parts := strings.Split(authHeader, " ")
	if parts[0] == "Bearer" {
		tokenString = parts[1]
	}

	return tokenString
}

func ParseToken(tokenString string) (Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return JWTKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return Claims{}, errors.New("signature invalid")
		} else {
			return Claims{}, errors.New("can't parse token string")
		}
	}

	if !token.Valid {
		return Claims{}, errors.New("token invalid")
	}

	return *claims, nil
}

func GenerateTokenString(claims Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(JWTKey)
	if err != nil {
		return "", errors.New("fail to sign token")
	}

	return tokenString, nil
}

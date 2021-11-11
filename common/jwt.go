package common

import (
	"errors"

	"github.com/dgrijalva/jwt-go"
)

// JWTDecode
func JWTDecode(tokenString string) (*jwt.StandardClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(JWTSignKey), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*jwt.StandardClaims); ok {
		// todo 刷新令牌生存时间
		return claims, nil
	}
	return nil, errors.New("token is not jwt.StandardClaims")

}

// JWTEncode
func JWTEncode(claims jwt.StandardClaims) (string, error) {
	mySigningKey := []byte(JWTSignKey)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(mySigningKey)
}

package pkgjwt

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

//// Create token
//token := jwt.New(jwt.SigningMethodHS256)
//
//// Set claims
//// This is the information which frontend can use
//// The backend can also decode the token and get admin etc.
//claims := token.Claims.(jwt.MapClaims)
//claims["sub"] = 1
//claims["name"] = "Jon Doe"
//claims["admin"] = true
//claims["exp"] = time.Now().Add(time.Minute * 15).Unix()
//
//// Generate encoded token and send it as response.
//// The signing string should be secret (a generated UUID works too)
//t, err := token.SignedString([]byte("secret"))
//if err != nil {
//return nil, err
//}
//
//refreshToken := jwt.New(jwt.SigningMethodHS256)
//rtClaims := refreshToken.Claims.(jwt.MapClaims)
//rtClaims["sub"] = 1
//rtClaims["exp"] = time.Now().Add(time.Hour * 24).Unix()
//
//rt, err := refreshToken.SignedString([]byte("secret"))
//if err != nil {
//return nil, err
//}
//
//return map[string]string{
//"access_token":  t,
//"refresh_token": rt,
//}, nil

func GenerateJWT(mapClaims jwt.MapClaims, secretSign []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, mapClaims)

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(secretSign)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseJWT(tokenString string, secretSign []byte) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")

		return secretSign, nil
	})

	if err != nil {
		return jwt.MapClaims{}, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return jwt.MapClaims{}, errors.New("invalid token")
}

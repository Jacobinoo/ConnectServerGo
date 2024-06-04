package Security

import (
	"os"

	"github.com/golang-jwt/jwt/v5"
)

// type UserClaims struct {
// 	jwt.Claims
// }

func constructJwtToken(secretKey string, claims *jwt.RegisteredClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(secretKey))
	return ss, err
}

func ConstructAccessToken() (string, error) {
	return constructJwtToken(os.Getenv("AT_PRIVATE_B64"), &jwt.RegisteredClaims{
		Issuer: "Connect",
	})
}

func ConstructRefreshToken() (string, error) {
	return constructJwtToken(os.Getenv("RT_PRIVATE_B64"), &jwt.RegisteredClaims{
		Issuer: "Connect",
	})
}

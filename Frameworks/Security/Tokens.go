package Security

import (
	"os"
	"strings"

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

// RetrieveBearerTokenFromAuthHeader derives the actual Bearer token from the provided Authorization header and returns it. If the header value is not properly formatted the function returns an empty string.
func RetrieveBearerTokenFromAuthHeader(rawAuthorizationHeaderValue string) string {
	splitHeader := strings.Split(rawAuthorizationHeaderValue, "Bearer")
	if len(splitHeader) != 2 {
		//Error: Bearer token not properly formatted
		return ""
	}
	obtainedToken := strings.TrimSpace(splitHeader[1])
	return obtainedToken
}

package Security

import (
	"errors"
	"os"
	"strings"
	"time"

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
		Issuer:    "Connect",
		ExpiresAt: jwt.NewNumericDate(time.Now().AddDate(0, 0, 14)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
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

func ValidateRefreshToken(refreshToken string) error {
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("RT_PRIVATE_B64")), nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return errors.New("Invalid refresh token")
	}

	return nil
}

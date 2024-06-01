package Account

import (
	"ConnectServer/Frameworks/Security"
	"errors"
	"log"
)

func GenerateTokenPair() (accessToken, refreshToken string, error error) {
	at, err := Security.ConstructAccessToken()
	if err != nil {
		log.Println(errAccountAccessTokenGenerationFailed)
		return "", "", errAccountAccessTokenGenerationFailed
	}
	rt, err := Security.ConstructRefreshToken()
	if err != nil {
		log.Println(errAccountRefreshTokenGenerationFailed)
		return "", "", errAccountRefreshTokenGenerationFailed
	}
	return at, rt, nil
}

var errAccountAccessTokenGenerationFailed = errors.New("couldn't generate access token")
var errAccountRefreshTokenGenerationFailed = errors.New("couldn't generate refresh token")

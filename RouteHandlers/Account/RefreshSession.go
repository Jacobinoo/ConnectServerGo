package Account

import (
	"ConnectServer/Frameworks/Security"
	"ConnectServer/Helpers"
	"ConnectServer/Types"
	"encoding/json"
	"log"
	"net/http"
)

func RefreshSessionHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	encoder := *json.NewEncoder(writer)

	refreshToken := Security.RetrieveBearerTokenFromAuthHeader(request.Header.Get("Authorization"))
	if refreshToken == "" {
		log.Println("Invalid refresh token")
		Helpers.JSONError(encoder, writer, Types.HttpErrorResponse{
			HttpResponse: Types.HttpResponse{
				Success: false,
			},
			Error: "Invalid refresh token",
		}, http.StatusUnauthorized)
		return
	}

	refreshTokenValidationError := Security.ValidateRefreshToken(refreshToken)

	if refreshTokenValidationError != nil {
		log.Println("Refresh token validation failed:", refreshTokenValidationError)
		Helpers.JSONError(encoder, writer, Types.HttpErrorResponse{
			HttpResponse: Types.HttpResponse{
				Success: false,
			},
			Error: "Refresh token validation failed",
		}, http.StatusUnauthorized)
		return
	}

	at, rt, err := GenerateTokenPair()
	if err != nil {
		log.Println("New token pair couldn't be created")
		Helpers.JSONError(encoder, writer, Types.HttpErrorResponse{
			HttpResponse: Types.HttpResponse{
				Success: false,
			},
			Error: "Token pair could not be generated",
		}, http.StatusInternalServerError)
		return
	}

	response := Types.HttpAuthResponse{
		AccessToken:  at,
		RefreshToken: rt,
		HttpResponse: Types.HttpResponse{
			Success: true,
		},
	}
	encoder.Encode(response)
}

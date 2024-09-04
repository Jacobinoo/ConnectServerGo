package main

import (
	"ConnectServer/Frameworks/CoreData"
	"ConnectServer/Frameworks/Security"
	"ConnectServer/Helpers"
	"ConnectServer/RouteHandlers/Account"
	"ConnectServer/Types"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	appEnvPath := os.Getenv("APP_ENV_PATH")
	if appEnvPath == "" {
		appEnvPath = ".env"
	}
	err := godotenv.Load(appEnvPath)
	if err != nil {
		log.Fatal("error loading .env file")
	}

	CoreData.ConnectUserServices()
	CoreData.ConnectStorageServices()

	router := http.NewServeMux()

	router.HandleFunc("POST /Account/SignIn", Account.SignInHandler)
	router.HandleFunc("POST /Account/SignUp", Account.SignUpHandler)
	router.HandleFunc("GET /Account/RefreshSession", Account.RefreshSessionHandler)

	// router.HandleFunc("GET /Conversations", Conversation.FetchManyConversationsHandler)

	http.ListenAndServe(os.Getenv("NETWORK_ADDR"), router)
}

func AuthGuard(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		encoder := *json.NewEncoder(w)
		log.Println("Before handler")

		accessToken := Security.RetrieveBearerTokenFromAuthHeader(r.Header.Get("Authorization"))
		subject, err := Security.VerifyAccessTokenAndDeriveOwnerId(accessToken)
		if err != nil || subject == "" {
			if err != nil {
				log.Print(err)
			}
			Helpers.JSONError(encoder, w, Types.HttpErrorResponse{
				Error: "Access Token Invalid",
				HttpResponse: Types.HttpResponse{
					Success: false,
				},
			}, http.StatusUnauthorized)
		}

		ctx := context.WithValue(r.Context(), "tokenSubject", subject)
		// ctx = context.WithValue(ctx, "userRole", userRole)
		// ctx = context.WithValue(ctx, "userID", userID)

		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)

		// Logic after the handler
		log.Println("After handler")
	}
}

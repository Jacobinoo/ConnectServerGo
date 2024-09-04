package main

import (
	"ConnectServer/Frameworks/CoreData"
	"ConnectServer/Frameworks/Security"
	"ConnectServer/Helpers"
	"ConnectServer/RouteHandlers/Account"
	"ConnectServer/RouteHandlers/Conversation"
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

	router.HandleFunc("GET /Conversations", AuthGuard(Conversation.FetchManyConversationsHandler))

	http.ListenAndServe(os.Getenv("NETWORK_ADDR"), router)
}

// AuthGuard is a middleware wrapper that guards the wrapped route handler preventing any unauthorized requests. It validates the access token, derives the identity of the token owner, and finally passes it into the request context. On validation failure it returns, thus preventing the call of the route handler function.
//
// Usage: AuthGuard(Account.SignInRouteHandler))
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
			return
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

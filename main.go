package main

import (
	"ConnectServer/Frameworks/CoreData"
	"ConnectServer/RouteHandlers/Account"
	"log"
	"net/http"
	"os"

	// _ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v5"
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

	CoreData.Connect()

	router := http.NewServeMux()

	router.HandleFunc("POST /Account/SignIn", Account.SignInHandler)
	router.HandleFunc("POST /Account/SignUp", Account.SignUpHandler)
	router.HandleFunc("GET /Account/RefreshSession", Account.RefreshSessionHandler)

	http.ListenAndServe(os.Getenv("NETWORK_ADDR"), router)
}

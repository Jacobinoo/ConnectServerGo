package main

import (
	"ConnectServer/Frameworks/CoreData"
	"ConnectServer/RouteHandlers/Account"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("error loading .env file")
	}

	CoreData.Connect()

	router := http.NewServeMux()

	router.HandleFunc("POST /Account/SignIn", Account.SignInHandler)
	router.HandleFunc("POST /Account/SignUp", Account.SignUpHandler)

	http.ListenAndServe("localhost:3000", router)
}

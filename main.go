package main

import (
	"ConnectServer/Frameworks/CoreData"
	"ConnectServer/RouteHandlers/Account"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	CoreData.Connect()

	router := http.NewServeMux()

	router.HandleFunc("POST /Account/SignIn", Account.SignInHandler)

	http.ListenAndServe("localhost:3000", router)
}

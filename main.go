package main

import (
	"ConnectServer/RouteHandlers/Account"
	"net/http"
)

func main() {
	router := http.NewServeMux()

	router.HandleFunc("POST /Account/SignIn", Account.SignInHandler)

	http.ListenAndServe("localhost:3000", router)
}

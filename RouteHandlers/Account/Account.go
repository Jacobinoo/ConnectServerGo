package Account

import (
	"ConnectServer/Helpers"
	"ConnectServer/Types"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func SignInHandler(writer http.ResponseWriter, request *http.Request) {
	var account Types.AccountLoginData

	err := Helpers.DecodeJSONBody(writer, request, &account)
	if err != nil {
		var malformedReq *Helpers.MalformedRequest
		if errors.As(err, &malformedReq) {
			http.Error(writer, malformedReq.Msg, malformedReq.Status)
		} else {
			log.Print(err.Error())
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}
	connectToDatabase()
}

func connectToDatabase() {
	fmt.Print(("s"))
	db, err := sql.Open("mysql", "mysql://root:SMJyflsGzCIWuqmGSPtmcHFZxCLxQAsX@roundhouse.proxy.rlwy.net:11308/railway")
	if err != nil {
		log.Fatal(err)
	}
	if err = db.Ping(); err != nil {
		log.Fatalf("Cannot ping database because %s", err)
	}
	log.Println("Successfully connected to database and pinged it")
}

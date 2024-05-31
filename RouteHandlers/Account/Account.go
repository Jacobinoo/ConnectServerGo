package Account

import (
	"ConnectServer/Frameworks/CoreData"
	"ConnectServer/Frameworks/Security"
	"ConnectServer/Helpers"
	"ConnectServer/Types"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
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

	passwordHash, err := fetchPasswordHashMatchingEmail(&account)
	if err != nil {
		if errors.Is(err, errAccountEmailNotFound) {
			http.Error(writer, "Invalid credentials", http.StatusUnauthorized)
			return
		}
		http.Error(writer, Helpers.InternalServerErrorHttpResponseMessage, http.StatusInternalServerError)
		return
	}
	if passwordHash == "" {
		http.Error(writer, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	if !Security.VerifyPassword(account.Password, passwordHash) {
		log.Println(errAccountWrongPassword)
		http.Error(writer, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	log.Println("successful login for", account.Email)
	fmt.Fprintln(writer, "Successful login!")
}

func fetchPasswordHashMatchingEmail(account *Types.AccountLoginData) (accountPasswordHash string, error error) {
	var row Types.AccountLoginData

	err := CoreData.DatabaseInstance.QueryRow("SELECT email,password FROM `Accounts` WHERE email=? LIMIT 1", account.Email).Scan(&row.Email, &row.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println(errAccountEmailNotFound)
			return "", errAccountEmailNotFound
		}
		log.Println(err)
		return "", err
	}
	return row.Password, nil
}

var errAccountEmailNotFound = errors.New("email doesn't exist in db")
var errAccountWrongPassword = errors.New("passwords hashes don't match")

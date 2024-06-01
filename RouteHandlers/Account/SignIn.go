package Account

import (
	"ConnectServer/Frameworks/CoreData"
	"ConnectServer/Frameworks/Security"
	"ConnectServer/Helpers"

	"ConnectServer/Types"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

func SignInHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(writer)
	var response Types.HttpSignInResponse

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
		if errors.Is(err, Helpers.ErrSignInEmailNotFound) {
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
		log.Println(Helpers.ErrSignInWrongPassword)
		http.Error(writer, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	log.Println("successful login for", account.Email)

	at, rt, err := GenerateTokenPair()
	if err != nil {
		log.Println("token pair generation failed")
		http.Error(writer, "Authentication token pair could not be generated", http.StatusInternalServerError)
		return
	}

	response.Success = true
	response = Types.HttpSignInResponse{
		AccessToken:  at,
		RefreshToken: rt,
		HttpResponse: Types.HttpResponse{
			Success: true,
		},
	}
	encoder.Encode(response)

}

func fetchPasswordHashMatchingEmail(account *Types.AccountLoginData) (accountPasswordHash string, error error) {
	var row Types.AccountLoginData

	err := CoreData.DatabaseInstance.QueryRow("SELECT email,password FROM `Accounts` WHERE email=? LIMIT 1", account.Email).Scan(&row.Email, &row.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println(Helpers.ErrSignInEmailNotFound)
			return "", Helpers.ErrSignInEmailNotFound
		}
		log.Println(err)
		return "", err
	}
	return row.Password, nil
}

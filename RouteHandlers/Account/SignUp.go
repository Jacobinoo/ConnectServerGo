package Account

import (
	"ConnectServer/Helpers"
	"ConnectServer/Types"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

func SignUpHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(writer)
	var response Types.HttpSignInResponse

	var accountToRegister Types.AccountRegisterData

	err := Helpers.DecodeJSONBody(writer, request, &accountToRegister)
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

	registerDataValid, err := validateRegistrationFormData(&accountToRegister)
	fmt.Println(registerDataValid, err)

	log.Println("successful register for", accountToRegister.Email)

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

func validateRegistrationFormData(registerFormData *Types.AccountRegisterData) (formDataValid bool, error error) {
	if registerFormData.Email == "" {
		log.Println()
	}
	return true, nil
}

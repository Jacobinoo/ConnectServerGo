package Account

import (
	"ConnectServer/Helpers"
	"ConnectServer/Types"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"unicode/utf8"
)

func SignUpHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(writer)
	var response Types.HttpSignInResponse

	var accountToRegister Types.AccountRegisterData

	decodingError := Helpers.DecodeJSONBody(writer, request, &accountToRegister)
	if decodingError != nil {
		var malformedReq *Helpers.MalformedRequest
		if errors.As(decodingError, &malformedReq) {
			http.Error(writer, malformedReq.Msg, malformedReq.Status)
		} else {
			log.Print(decodingError.Error())
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	validationError := validateRegistrationFormData(&accountToRegister)
	if validationError != nil {
		log.Println(validationError)
		switch {
		case errors.Is(validationError, Helpers.ErrSignUpInvalidEmail):
			http.Error(writer, "Provided email is invalid", http.StatusBadRequest)
			return
		case errors.Is(validationError, Helpers.ErrSignUpPasswordTooLong):
			http.Error(writer, "Password is too long", http.StatusBadRequest)
			return
		case errors.Is(validationError, Helpers.ErrSignUpPasswordTooShort):
			http.Error(writer, "Password is too short", http.StatusBadRequest)
			return
		case errors.Is(validationError, Helpers.ErrSignUpPasswordNoDigit):
			http.Error(writer, "Password does not contain a digit", http.StatusBadRequest)
			return
		case errors.Is(validationError, Helpers.ErrSignUpPasswordNoLower):
			http.Error(writer, "Password does not contain a lowercase letter", http.StatusBadRequest)
			return
		case errors.Is(validationError, Helpers.ErrSignUpPasswordNoUpper):
			http.Error(writer, "Password does not contain an uppercase letter", http.StatusBadRequest)
			return
		case errors.Is(validationError, Helpers.ErrSignUpPasswordNoSpecial):
			http.Error(writer, "Password does not contain a special character", http.StatusBadRequest)
			return
		case errors.Is(validationError, Helpers.ErrSignUpPasswordMismatch):
			http.Error(writer, "Provided password and confirm password are different", http.StatusBadRequest)
			return
		case errors.Is(validationError, Helpers.ErrSignUpInvalidFirstName):
			http.Error(writer, "Provided first name is invalid", http.StatusBadRequest)
			return
		case errors.Is(validationError, Helpers.ErrSignUpInvalidLastName):
			http.Error(writer, "Provided last name is invalid", http.StatusBadRequest)
			return
		case errors.Is(validationError, Helpers.ErrSignUpNameContainsWhitespace):
			http.Error(writer, "Names can't contain whitespaces", http.StatusBadRequest)
			return
		}
	}

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

func validateRegistrationFormData(registerFormData *Types.AccountRegisterData) error {
	if !Helpers.EmailRegexp.MatchString(registerFormData.Email) {
		return Helpers.ErrSignUpInvalidEmail
	}

	_, err := Helpers.PasswordValid(registerFormData.Password)
	if err != nil {
		return err
	}

	if registerFormData.Password != registerFormData.ConfirmPassword {
		return Helpers.ErrSignUpPasswordMismatch
	}

	if utf8.RuneCountInString(registerFormData.FirstName) > 35 || utf8.RuneCountInString(registerFormData.FirstName) <= 0 {
		return Helpers.ErrSignUpInvalidFirstName
	}

	if utf8.RuneCountInString(registerFormData.LastName) > 35 || utf8.RuneCountInString(registerFormData.LastName) <= 0 {
		return Helpers.ErrSignUpInvalidLastName
	}

	if Helpers.ContainsWhitespaceRegexp.MatchString(registerFormData.FirstName) || Helpers.ContainsWhitespaceRegexp.MatchString(registerFormData.LastName) {
		return Helpers.ErrSignUpNameContainsWhitespace
	}

	return nil
}

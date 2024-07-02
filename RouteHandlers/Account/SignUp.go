package Account

import (
	"ConnectServer/Frameworks/CoreData"
	"ConnectServer/Frameworks/Security"
	"ConnectServer/Helpers"
	"ConnectServer/Types"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"unicode/utf8"
)

func SignUpHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	encoder := *json.NewEncoder(writer)

	var accountToRegister Types.AccountRegisterData

	decodingError := Helpers.DecodeJSONBody(writer, request, &accountToRegister)
	if decodingError != nil {
		var malformedReq *Helpers.MalformedRequest
		if errors.As(decodingError, &malformedReq) {
			Helpers.JSONError(encoder, writer, Types.HttpErrorResponse{
				HttpResponse: Types.HttpResponse{
					Success: false,
				},
				Error: malformedReq.Msg,
			}, malformedReq.Status)
		} else {
			log.Print(decodingError.Error())
			Helpers.JSONError(encoder, writer, Types.HttpErrorResponse{
				HttpResponse: Types.HttpResponse{
					Success: false,
				},
				Error: Helpers.InternalServerErrorHttpResponseMessage,
			}, http.StatusInternalServerError)
		}
		return
	}

	validationError := validateRegistrationFormData(&accountToRegister)
	if validationError != nil {
		log.Println(validationError)
		switch {
		case errors.Is(validationError, Helpers.ErrSignUpInvalidEmail):
			Helpers.JSONError(encoder, writer, Types.HttpErrorResponse{
				HttpResponse: Types.HttpResponse{
					Success: false,
				},
				Error: "Provided email is invalid",
			}, http.StatusBadRequest)
			return
		case errors.Is(validationError, Helpers.ErrSignUpPasswordTooLong):
			Helpers.JSONError(encoder, writer, Types.HttpErrorResponse{
				HttpResponse: Types.HttpResponse{
					Success: false,
				},
				Error: "Password is too long",
			}, http.StatusBadRequest)
			return
		case errors.Is(validationError, Helpers.ErrSignUpPasswordTooShort):
			Helpers.JSONError(encoder, writer, Types.HttpErrorResponse{
				HttpResponse: Types.HttpResponse{
					Success: false,
				},
				Error: "Password is too short",
			}, http.StatusBadRequest)
			return
		case errors.Is(validationError, Helpers.ErrSignUpPasswordNoDigit):
			Helpers.JSONError(encoder, writer, Types.HttpErrorResponse{
				HttpResponse: Types.HttpResponse{
					Success: false,
				},
				Error: "Password does not contain a digit",
			}, http.StatusBadRequest)
			return
		case errors.Is(validationError, Helpers.ErrSignUpPasswordNoLower):
			Helpers.JSONError(encoder, writer, Types.HttpErrorResponse{
				HttpResponse: Types.HttpResponse{
					Success: false,
				},
				Error: "Password does not contain a lowercase letter",
			}, http.StatusBadRequest)
			return
		case errors.Is(validationError, Helpers.ErrSignUpPasswordNoUpper):
			Helpers.JSONError(encoder, writer, Types.HttpErrorResponse{
				HttpResponse: Types.HttpResponse{
					Success: false,
				},
				Error: "Password does not contain an uppercase letter",
			}, http.StatusBadRequest)
			return
		case errors.Is(validationError, Helpers.ErrSignUpPasswordNoSpecial):
			Helpers.JSONError(encoder, writer, Types.HttpErrorResponse{
				HttpResponse: Types.HttpResponse{
					Success: false,
				},
				Error: "Password does not contain a special character",
			}, http.StatusBadRequest)
			return
		case errors.Is(validationError, Helpers.ErrSignUpPasswordMismatch):
			Helpers.JSONError(encoder, writer, Types.HttpErrorResponse{
				HttpResponse: Types.HttpResponse{
					Success: false,
				},
				Error: "Provided password and confirm password are different",
			}, http.StatusBadRequest)
			return
		case errors.Is(validationError, Helpers.ErrSignUpInvalidFirstName):
			Helpers.JSONError(encoder, writer, Types.HttpErrorResponse{
				HttpResponse: Types.HttpResponse{
					Success: false,
				},
				Error: "Provided first name is invalid",
			}, http.StatusBadRequest)
			return
		case errors.Is(validationError, Helpers.ErrSignUpInvalidLastName):
			Helpers.JSONError(encoder, writer, Types.HttpErrorResponse{
				HttpResponse: Types.HttpResponse{
					Success: false,
				},
				Error: "Provided last name is invalid",
			}, http.StatusBadRequest)
			return
		case errors.Is(validationError, Helpers.ErrSignUpNameContainsWhitespace):
			Helpers.JSONError(encoder, writer, Types.HttpErrorResponse{
				HttpResponse: Types.HttpResponse{
					Success: false,
				},
				Error: "Names can't contain whitespaces",
			}, http.StatusBadRequest)
			return
		}
	}

	insertionErr := insertAccountToDb(&accountToRegister)
	if insertionErr != nil {
		log.Println(insertionErr)
		Helpers.JSONError(encoder, writer, Types.HttpErrorResponse{
			HttpResponse: Types.HttpResponse{
				Success: false,
			},
			Error: "Account creation failed due to an internal server error.",
		}, http.StatusInternalServerError)
		return
	}

	log.Println("successful register for", accountToRegister.Email)

	at, rt, err := GenerateTokenPair()
	if err != nil {
		log.Println("token pair generation failed")
		Helpers.JSONError(encoder, writer, Types.HttpErrorResponse{
			HttpResponse: Types.HttpResponse{
				Success: false,
			},
			Error: "Authentication token pair could not be generated",
		}, http.StatusInternalServerError)
		return
	}

	response := Types.HttpAuthResponse{
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

func insertAccountToDb(data *Types.AccountRegisterData) error {
	passwordHash, hasherErr := Security.HashPassword(data.Password)
	if hasherErr != nil {
		log.Print(hasherErr)
		return Helpers.ErrHasherHashNew
	}

	query := "INSERT INTO `Accounts` (`id`, `password`, `email`, `firstName`, `lastName`) VALUES (?, ?, ?, ?, ?)"

	newUUID, UUIDErr := Security.GenerateUUID()
	if UUIDErr != nil {
		log.Print(UUIDErr)
		return Helpers.ErrUUIDGenerationFailed
	}

	_, err := CoreData.DatabaseInstance.ExecContext(context.Background(), query, newUUID, passwordHash, data.Email, data.FirstName, data.LastName)
	if err != nil {
		log.Print(err)
		return Helpers.ErrInsertionFailed
	}

	return nil
}

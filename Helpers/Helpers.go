package Helpers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"unicode"
)

type MalformedRequest struct {
	Status int
	Msg    string
}

func (mr *MalformedRequest) Error() string {
	return mr.Msg
}

func DecodeJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	ct := r.Header.Get("Content-Type")
	if ct != "" {
		mediaType := strings.ToLower(strings.TrimSpace(strings.Split(ct, ";")[0]))
		if mediaType != "application/json" {
			msg := "Content-Type header is not application/json"
			return &MalformedRequest{Status: http.StatusUnsupportedMediaType, Msg: msg}
		}
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(&dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
			return &MalformedRequest{Status: http.StatusBadRequest, Msg: msg}

		case errors.Is(err, io.ErrUnexpectedEOF):
			msg := "Request body contains badly-formed JSON"
			return &MalformedRequest{Status: http.StatusBadRequest, Msg: msg}

		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			return &MalformedRequest{Status: http.StatusBadRequest, Msg: msg}

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
			return &MalformedRequest{Status: http.StatusBadRequest, Msg: msg}

		case errors.Is(err, io.EOF):
			msg := "Request body must not be empty"
			return &MalformedRequest{Status: http.StatusBadRequest, Msg: msg}

		case err.Error() == "http: request body too large":
			msg := "Request body must not be larger than 1MB"
			return &MalformedRequest{Status: http.StatusRequestEntityTooLarge, Msg: msg}

		default:
			return err
		}
	}

	err = dec.Decode(&struct{}{})
	if !errors.Is(err, io.EOF) {
		msg := "Request body must only contain a single JSON object"
		return &MalformedRequest{Status: http.StatusBadRequest, Msg: msg}
	}

	return nil
}

const InternalServerErrorHttpResponseMessage string = "An unknown error occured on our side. We're sorry for the in"

// Validation

var EmailRegexp *regexp.Regexp = regexp.MustCompile(`^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$`)
var ContainsWhitespaceRegexp *regexp.Regexp = regexp.MustCompile(`\s`)

// Returns 'true <nil>' if password string is valid (meets our password requirements)
func PasswordValid(s string) (bool, error) {
	var oneDigitOk, oneUpperOk, oneLowerOk, oneSpecialOk bool
	characters := 0
	for _, c := range s {
		switch {
		case unicode.IsNumber(c):
			oneDigitOk = true
		case unicode.IsUpper(c):
			oneUpperOk = true
		case unicode.IsLower(c):
			oneLowerOk = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c) || c == ' ':
			oneSpecialOk = true
		}
		characters += 1
	}

	tooShort := characters < 8
	tooLong := characters > 255

	switch {
	case tooShort:
		return false, ErrSignUpPasswordTooShort
	case tooLong:
		return false, ErrSignUpPasswordTooLong
	case !oneLowerOk:
		return false, ErrSignUpPasswordNoLower
	case !oneUpperOk:
		return false, ErrSignUpPasswordNoUpper
	case !oneDigitOk:
		return false, ErrSignUpPasswordNoDigit
	case !oneSpecialOk:
		return false, ErrSignUpPasswordNoSpecial
	}

	return true, nil
}

var ErrSignInEmailNotFound = errors.New("email doesn't exist in db")
var ErrSignInWrongPassword = errors.New("passwords hashes don't match")

var ErrSignUpInvalidEmail = errors.New("email is invalid")
var ErrSignUpPasswordMismatch = errors.New("password and confirmPassword are not equal")
var ErrSignUpInvalidFirstName = errors.New("first name is invalid")
var ErrSignUpInvalidLastName = errors.New("last name is invalid")
var ErrSignUpNameContainsWhitespace = errors.New("names cannot contain whitespaces")

var ErrSignUpPasswordTooShort = errors.New("password too short")
var ErrSignUpPasswordTooLong = errors.New("password too long")
var ErrSignUpPasswordNoUpper = errors.New("password needs to have at least one uppercase")
var ErrSignUpPasswordNoLower = errors.New("password needs to have at least one lowercase")
var ErrSignUpPasswordNoSpecial = errors.New("password needs to have at least one special")
var ErrSignUpPasswordNoDigit = errors.New("password needs to have at least one digit")

var ErrInsertionFailed = errors.New("row insertion failed")
var ErrHasherHashNew = errors.New("security framework's hasher failed")
var ErrUUIDGenerationFailed = errors.New("security framework's hasher failed")

type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

func JSONError(encoder json.Encoder, writer http.ResponseWriter, err interface{}, code int) {
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	writer.Header().Set("X-Content-Type-Options", "nosniff")
	writer.WriteHeader(code)
	json.NewEncoder(writer).Encode(err)
}

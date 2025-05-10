package helper

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"regexp"
	"slices"
	"unicode/utf8"
)

func CheckPassword(password string) error {
	hasLower, _ := regexp.MatchString(`[a-z]`, password)
	hasUpper, _ := regexp.MatchString(`[A-Z]`, password)
	hasDigit, _ := regexp.MatchString(`\d`, password)
	hasSpecial, _ := regexp.MatchString(`[^\da-zA-Z]`, password)
	hasLength := utf8.RuneCountInString(password) >= 8

	if !(hasLower && hasUpper && hasDigit && hasSpecial && hasLength) {
		return errors.New("password is invalid, it must be min 8 chars long, contain 1 special char, 1 upper and lower char and 1 digit")
	}

	return nil
}

// helper to reject specified method that are not allowed,
// return false if a request use a not allowed method else true
func MethodsAllowed(w http.ResponseWriter, r *http.Request, allowed ...string) bool {

	// if request method is in allowed methods then return
	if slices.Contains(allowed, r.Method) {
		return true
	}

	JsonResponse(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})

	return false
}

// helper to encode and write an http respose as json
func JsonResponse(w http.ResponseWriter, httpStatus int, jsonContent any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	json.NewEncoder(w).Encode(jsonContent)
}

// helper to read body and unmarshal it to a struct
func ReadBody(w http.ResponseWriter, r *http.Request, payloadStruct any) any {

	// check body content type
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		return nil
	}

	defer r.Body.Close()

	// read body
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return nil
	}

	// check if body is empty
	if len(bodyBytes) == 0 {
		return nil
	}

	// unmashall it and check if it errored
	err = json.Unmarshal(bodyBytes, &payloadStruct)
	if err != nil {
		return nil
	}

	return payloadStruct
}

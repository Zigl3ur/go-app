package helper

import (
	"encoding/json"
	"io"
	"net/http"
	"slices"
)

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

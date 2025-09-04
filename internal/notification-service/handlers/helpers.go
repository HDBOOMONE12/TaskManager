package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strings"
)

var (
	ErrBodyTooLarge  = errors.New("body too large")
	ErrInvalidSyntax = errors.New("invalid json syntax")
	ErrWrongType     = errors.New("wrong json type")
	ErrUnknownField  = errors.New("unknown field")
	ErrEmptyBody     = errors.New("empty body")
	ErrTrailingData  = errors.New("multiple json values")
)

func ParseTgMsg(r *http.Request) {

}

func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("writeJSON encode error: %v", err)
	}
}

func errorJSON(w http.ResponseWriter, status int, msg string) {
	err := struct {
		Error string `json:"error"`
		Code  string `json:"code"`
	}{
		Error: msg,
		Code:  http.StatusText(status)}
	WriteJSON(w, status, err)
}

func decodeJSON(w http.ResponseWriter, r *http.Request, dst any, maxBytes int64) error {

	r.Body = http.MaxBytesReader(w, r.Body, maxBytes)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(dst); err != nil {

		if err == io.EOF {
			return ErrEmptyBody
		}

		var se *json.SyntaxError
		if errors.As(err, &se) || errors.Is(err, io.ErrUnexpectedEOF) {
			return ErrInvalidSyntax
		}

		var te *json.UnmarshalTypeError
		if errors.As(err, &te) {
			return ErrWrongType
		}

		if strings.Contains(err.Error(), "unknown field") {
			return ErrUnknownField
		}

		if strings.Contains(err.Error(), "request body too large") {
			return ErrBodyTooLarge
		}

		return err
	}

	var extra any
	if err := dec.Decode(&extra); err != io.EOF {
		return ErrTrailingData
	}

	return nil
}

func respondDecodeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrBodyTooLarge):
		errorJSON(w, http.StatusRequestEntityTooLarge, "request body too large")
	case errors.Is(err, ErrInvalidSyntax):
		errorJSON(w, http.StatusBadRequest, "invalid JSON")
	case errors.Is(err, ErrWrongType):
		errorJSON(w, http.StatusBadRequest, "wrong type")
	case errors.Is(err, ErrUnknownField):
		errorJSON(w, http.StatusBadRequest, "unknown field")
	case errors.Is(err, ErrEmptyBody):
		errorJSON(w, http.StatusBadRequest, "empty body")
	case errors.Is(err, ErrTrailingData):
		errorJSON(w, http.StatusBadRequest, "body must contain a single JSON value")
	default:
		errorJSON(w, http.StatusBadRequest, "invalid JSON")
	}
}

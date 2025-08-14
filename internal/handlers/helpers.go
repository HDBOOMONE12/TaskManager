package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var (
	errBadPath   = errors.New("bad user path")
	errBadID     = errors.New("bad user id")
	errBadTaskID = errors.New("bad task id")
)

func parseUserID(r *http.Request) (int, error) {
	path := r.URL.Path
	if !strings.HasPrefix(path, "/users/") {
		return 0, errBadPath
	}

	parts := strings.Split(path, "/")

	if !(len(parts) == 3 || (len(parts) == 4 && parts[3] == "")) {
		return 0, errBadPath
	}

	idStr := parts[2]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, errBadID
	}
	return id, nil
}

func writeJSON(w http.ResponseWriter, status int, v any) {
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
	writeJSON(w, status, err)
}

var (
	ErrBodyTooLarge  = errors.New("body too large")
	ErrInvalidSyntax = errors.New("invalid json syntax")
	ErrWrongType     = errors.New("wrong json type")
	ErrUnknownField  = errors.New("unknown field")
	ErrEmptyBody     = errors.New("empty body")
	ErrTrailingData  = errors.New("multiple json values")
)

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

func parseUserTasksListPath(r *http.Request) (int, error) {
	path := r.URL.Path
	parts := strings.Split(path, "/")

	if !(len(parts) == 4 || (len(parts) == 5 && parts[4] == "")) {
		return 0, errBadPath
	}
	if parts[1] != "users" || parts[3] != "tasks" {
		return 0, errBadPath
	}

	uid, err := strconv.Atoi(parts[2])
	if err != nil {
		return 0, errBadID
	}
	return uid, nil
}

func parseUserTaskDetailPath(r *http.Request) (int, int, error) {
	parts := strings.Split(r.URL.Path, "/")

	if !((len(parts) == 5 && parts[1] == "users" && parts[3] == "tasks") ||
		(len(parts) == 6 && parts[1] == "users" && parts[3] == "tasks" && parts[5] == "")) {
		return 0, 0, errBadPath
	}

	uid, err := strconv.Atoi(parts[2])
	if err != nil {
		return 0, 0, errBadID
	}
	tid, err := strconv.Atoi(parts[4])
	if err != nil {
		return 0, 0, errBadTaskID
	}
	return uid, tid, nil
}

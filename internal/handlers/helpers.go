package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var (
	errBadPath = errors.New("bad user path")
	errBadID   = errors.New("bad user id")
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

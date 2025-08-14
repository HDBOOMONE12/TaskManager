package handlers

import (
	"net/http"
	"strings"
)

func UsersSubtreeHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	parts := strings.Split(path, "/")

	if len(parts) < 2 || parts[1] != "users" {
		http.NotFound(w, r)
		return
	}

	if len(parts) == 3 && parts[2] == "" {
		UsersHandler(w, r)
		return
	}

	if (len(parts) == 3 && parts[2] != "") ||
		(len(parts) == 4 && parts[2] != "" && parts[3] == "") {
		UserDetailHandler(w, r)
		return
	}

	if (len(parts) == 4 && parts[2] != "" && parts[3] == "tasks") ||
		(len(parts) == 5 && parts[2] != "" && parts[3] == "tasks" && parts[4] == "") {
		UserTasksHandler(w, r)
		return
	}

	if (len(parts) == 5 && parts[2] != "" && parts[3] == "tasks") ||
		(len(parts) == 6 && parts[2] != "" && parts[3] == "tasks" && parts[5] == "") {
		UserTaskDetailHandler(w, r)
		return
	}

	http.NotFound(w, r)
}

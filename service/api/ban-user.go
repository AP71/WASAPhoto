package api

import (
	"encoding/json"
	"net/http"
	"strings"
	"wasa-photo/service/api/errors"
	"wasa-photo/service/api/structures"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) banUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, user structures.User) {
	w.Header().Set("content-type", "application/json")

	username := ps.ByName("username")

	if ps.ByName("byUsername") != user.Username.Value {
		errors.WriteResponse(rt.baseLogger, w, "Operation not permitted", http.StatusForbidden, "Unauthorized access: Operation not permitted")
		return
	}

	err := rt.db.BanUser(username, user)
	if err != nil && err.Error() == "user not found" {
		errors.WriteResponse(rt.baseLogger, w, "User not found", http.StatusNotFound, "User not found")
		return
	} else if err != nil && strings.Contains(err.Error(), "UNIQUE constraint failed") {
		errors.WriteResponse(rt.baseLogger, w, "Conflict error", http.StatusConflict, "You have already banned this user")
		return
	} else if err != nil {
		errors.WriteResponse(rt.baseLogger, w, "Database error: "+err.Error(), http.StatusInternalServerError, "Database error")
		return
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(errors.JSONMsg{Message: "Resource created"})
	if err != nil {
		errors.WriteResponse(rt.baseLogger, w, "banUser return an error.", http.StatusInternalServerError, "Internal server error")
		return
	}
}

package api

import (
	"net/http"
	"strings"
	"wasa-photo/service/api/errors"
	"wasa-photo/service/api/structures"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) unbanUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, user structures.User) {
	w.Header().Set("content-type", "application/json")

	username := ps.ByName("username")

	if ps.ByName("byUsername") != user.Username.Value {
		errors.WriteResponse(rt.baseLogger, w, "Operation not permitted", http.StatusForbidden, "Unauthorized access: Operation not permitted")
		return
	}

	err := rt.db.UnbanUser(username, user)
	if err != nil && err.Error() == "user not found" {
		errors.WriteResponse(rt.baseLogger, w, "User not found", http.StatusNotFound, "User not found")
		return
	} else if err != nil && strings.Contains(err.Error(), "relationship not found") {
		errors.WriteResponse(rt.baseLogger, w, "Conflict error", http.StatusConflict, "You have not banned this user")
		return
	} else if err != nil {
		errors.WriteResponse(rt.baseLogger, w, "Database error: "+err.Error(), http.StatusInternalServerError, "Database error")
		return
	}

	w.WriteHeader(http.StatusNoContent)
	if err != nil {
		errors.WriteResponse(rt.baseLogger, w, "unbanUser return an error.", http.StatusInternalServerError, "Internal server error")
		return
	}
}

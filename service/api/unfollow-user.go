package api

import (
	"net/http"
	"strings"
	"wasa-photo/service/api/auth"
	"wasa-photo/service/api/errors"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) unfollowUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	var username string
	var byUsername string

	//Getting user data
	res, user := auth.CheckAuth(rt.db, r)

	username = ps.ByName("username")
	byUsername = ps.ByName("byUsername")

	if !res {
		errors.WriteResponse(rt.baseLogger, w, "Authentication failed", http.StatusUnauthorized, "Unauthorized access")
		return
	}

	if byUsername != user.Username.Value {
		errors.WriteResponse(rt.baseLogger, w, "Operation not permitted", http.StatusForbidden, "Unauthorized access: Operation not permitted")
		return
	}

	err := rt.db.UnfollowUser(username, byUsername)
	if err != nil && err.Error() == "user not found" {
		errors.WriteResponse(rt.baseLogger, w, "User not found", http.StatusNotFound, "User not found")
		return
	} else if err != nil && strings.Contains(err.Error(), "relationship not found") {
		errors.WriteResponse(rt.baseLogger, w, "Conflict error", http.StatusConflict, "You have not follwed this user")
		return
	} else if err != nil {
		errors.WriteResponse(rt.baseLogger, w, "Database error", http.StatusInternalServerError, "Database error")
		return
	}

	w.WriteHeader(http.StatusNoContent)
	if err != nil {
		errors.WriteResponse(rt.baseLogger, w, "UpdateUsername return an error.", http.StatusInternalServerError, "Internal server error")
		return
	}
}

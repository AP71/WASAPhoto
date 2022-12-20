package api

import (
	"net/http"
	"wasa-photo/service/api/errors"
	"wasa-photo/service/api/structures"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) followStatus(w http.ResponseWriter, r *http.Request, ps httprouter.Params, user structures.User) {
	w.Header().Set("content-type", "application/json")

	var username structures.User

	username.Username.Value = ps.ByName("username")

	if ps.ByName("byUsername") != user.Username.Value {
		errors.WriteResponse(rt.baseLogger, w, "Operation not permitted", http.StatusForbidden, "Unauthorized access: Operation not permitted")
		return
	}

	err := rt.db.GetFollowStatus(username, user)
	if err != nil {
		errors.WriteResponse(rt.baseLogger, w, "Interaction not found (follow)", http.StatusNotFound, "Interaction not found.")
		return
	} else {
		w.WriteHeader(http.StatusNoContent)
	}

}

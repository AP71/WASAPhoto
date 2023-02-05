package api

import (
	"encoding/json"
	"net/http"
	"wasa-photo/service/api/errors"
	"wasa-photo/service/api/structures"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) banStatus(w http.ResponseWriter, r *http.Request, ps httprouter.Params, user structures.User) {
	w.Header().Set("content-type", "application/json")

	var username structures.User
	var status structures.Status
	username.Username.Value = ps.ByName("username")

	if ps.ByName("byUsername") != user.Username.Value {
		errors.WriteResponse(rt.baseLogger, w, "Operation not permitted", http.StatusForbidden, "Unauthorized access: Operation not permitted")
		return
	}

	err := rt.db.GetBanStatus(username, user)
	w.WriteHeader(http.StatusOK)
	if err != nil {
		status.Status = false
	} else {
		status.Status = true
	}
	err = json.NewEncoder(w).Encode(status)
	if err != nil {
		errors.WriteResponse(rt.baseLogger, w, "banStatus return an error.", http.StatusInternalServerError, "Internal server error")
		return
	}

}

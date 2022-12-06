package api

import (
	"encoding/json"
	"net/http"
	"wasa-photo/service/api/errors"
	"wasa-photo/service/api/structures"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) setMyUsername(w http.ResponseWriter, r *http.Request, ps httprouter.Params, user structures.User) {
	w.Header().Set("content-type", "application/json")

	var new structures.NewUsername

	if ps.ByName("username") != user.Username.Value {
		errors.WriteResponse(rt.baseLogger, w, "Operation not permitted", http.StatusForbidden, "Unauthorized access: Operation not permitted")
		return
	}

	err := json.NewDecoder(r.Body).Decode(&new)
	_ = r.Body.Close()

	if err != nil {
		errors.WriteResponse(rt.baseLogger, w, "Wrong JSON received", http.StatusBadRequest, "Wrong JSON received")
		return
	}

	if !structures.CheckUsername(new.Value) {
		errors.WriteResponse(rt.baseLogger, w, "Invalid username received", http.StatusBadRequest, "Invalid username received")
		return
	}

	user.Username.Value, err = rt.db.UpdateUsername(user, new)
	if err != nil {
		errors.WriteResponse(rt.baseLogger, w, "Database error", http.StatusInternalServerError, "Internal server error")
		return
	}

	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(user.Username)
	if err != nil {
		errors.WriteResponse(rt.baseLogger, w, "UpdateUsername return an error.", http.StatusInternalServerError, "Internal server error")
		return
	}
}

package api

import (
	"encoding/json"
	"net/http"
	"wasa-photo/service/api/errors"
	"wasa-photo/service/api/structures"

	"github.com/julienschmidt/httprouter"
)

// getHelloWorld is an example of HTTP endpoint that returns "Hello world!" as a plain text
func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	var user structures.User

	err := json.NewDecoder(r.Body).Decode(&user.Username)
	_ = r.Body.Close()
	//Checking JSON
	if err != nil {
		errors.WriteResponse(rt.baseLogger, w, "Wrong JSON received", http.StatusBadRequest, "Wrong JSON received")
		return
	}
	//Checking username
	if !user.IsValid() {
		errors.WriteResponse(rt.baseLogger, w, "Invalid username received", http.StatusBadRequest, "Invalid username received")
		return
	}
	//Getting username identifier
	user.Id.Value, err = rt.db.GetUserId(user.Username.Value)
	if user.Id.Value == "" || err != nil {
		user.Id.Value, err = rt.db.CreateUser(user.Username.Value)
		if err != nil {
			errors.WriteResponse(rt.baseLogger, w, "Database error", http.StatusInternalServerError, "Internal server error")
			return
		}
	}
	w.WriteHeader(http.StatusCreated)
	//Returning identifier
	err = json.NewEncoder(w).Encode(user.Id)
	if err != nil {
		errors.WriteResponse(rt.baseLogger, w, "DoLogin return an error.", http.StatusInternalServerError, "Internal server error")
		return
	}
}

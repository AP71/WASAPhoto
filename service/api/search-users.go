package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"wasa-photo/service/api/errors"
	"wasa-photo/service/api/structures"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getUsers(w http.ResponseWriter, r *http.Request, ps httprouter.Params, user structures.User) {
	w.Header().Set("content-type", "application/json")
	var users structures.Users
	var userToSearch string
	var pageId int64
	var err error

	if r.URL.Query().Has("userToSearch") {
		userToSearch = r.URL.Query().Get("userToSearch")
	}
	if r.URL.Query().Has("pageId") {
		pageId, err = strconv.ParseInt(r.URL.Query().Get("pageId"), 10, 64)
		if err != nil {
			errors.WriteResponse(rt.baseLogger, w, "Incorrect pageId number", http.StatusBadRequest, "Incorrect pageId number")
			return
		}
	} else {
		pageId = 0
	}

	users, err = rt.db.GetUsers(userToSearch, pageId, user.Username.Value)
	if err != nil {
		errors.WriteResponse(rt.baseLogger, w, "Database error: "+err.Error(), http.StatusInternalServerError, "Internal server error")
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(&users)
	if err != nil {
		errors.WriteResponse(rt.baseLogger, w, "getUsers return an error.", http.StatusInternalServerError, "Internal server error")
		return
	}

}

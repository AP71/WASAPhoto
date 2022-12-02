package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"wasa-photo/service/api/auth"
	"wasa-photo/service/api/errors"
	"wasa-photo/service/api/structures"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getUserProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")
	var user structures.UserPage
	var pageId int64
	var err error

	res, _ := auth.CheckAuth(rt.db, r)
	user.Username = ps.ByName("username")

	if !res {
		errors.WriteResponse(rt.baseLogger, w, "Authentication failed", http.StatusUnauthorized, "Unauthorized access")
		return
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

	user, err = rt.db.GetUserPage(user.Username, pageId)
	if err != nil && strings.Contains(err.Error(), "sql: no rows in result set") {
		errors.WriteResponse(rt.baseLogger, w, "File not found", http.StatusNotFound, "Image not found")
		return
	} else if err != nil {
		errors.WriteResponse(rt.baseLogger, w, "Database error", http.StatusInternalServerError, "Database error")
		return
	}

	err = json.NewEncoder(w).Encode(&user)
	if err != nil {
		errors.WriteResponse(rt.baseLogger, w, "DoLogin return an error.", http.StatusInternalServerError, "Internal server error")
		return
	}

}

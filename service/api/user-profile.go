package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"wasa-photo/service/api/errors"
	"wasa-photo/service/api/structures"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getUserProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params, user structures.User) {
	w.Header().Set("content-type", "application/json")
	var userPage structures.UserPage
	var pageId int64
	var err error

	userPage.Username = ps.ByName("username")

	if r.URL.Query().Has("pageId") {
		pageId, err = strconv.ParseInt(r.URL.Query().Get("pageId"), 10, 64)
		if err != nil {
			errors.WriteResponse(rt.baseLogger, w, "Incorrect pageId number", http.StatusBadRequest, "Incorrect pageId number")
			return
		}
	} else {
		pageId = 0
	}

	userPage, err = rt.db.GetUserPage(userPage.Username, pageId)
	if err != nil && strings.Contains(err.Error(), `sql: Scan error on column index 0, name "Id": converting NULL to string is unsupported`) {
		errors.WriteResponse(rt.baseLogger, w, "User not found", http.StatusNotFound, "User not found")
		return
	} else if err != nil {
		errors.WriteResponse(rt.baseLogger, w, "Database error", http.StatusInternalServerError, "Database error")
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(&userPage)
	if err != nil {
		errors.WriteResponse(rt.baseLogger, w, "getUserProfile return an error.", http.StatusInternalServerError, "Internal server error")
		return
	}

}

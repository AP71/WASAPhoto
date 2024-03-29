package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"wasa-photo/service/api/errors"
	"wasa-photo/service/api/structures"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getMyStream(w http.ResponseWriter, r *http.Request, ps httprouter.Params, user structures.User) {
	w.Header().Set("content-type", "application/json")
	var pageId int64
	var photos structures.Photos
	var err error

	if r.URL.Query().Has("pageId") {
		pageId, err = strconv.ParseInt(r.URL.Query().Get("pageId"), 10, 64)
		if err != nil {
			errors.WriteResponse(rt.baseLogger, w, "Incorrect pageId number", http.StatusBadRequest, "Incorrect pageId number")
			return
		}
	} else {
		pageId = 0
	}

	photos, err = rt.db.GetFeed(user, pageId)
	if err != nil {
		errors.WriteResponse(rt.baseLogger, w, "Database error: "+err.Error(), http.StatusInternalServerError, "Database error")
		return
	}

	if len(photos.Post) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(&photos)
	if err != nil {
		errors.WriteResponse(rt.baseLogger, w, "getMyStream return an error", http.StatusInternalServerError, "Internal server error")
		return
	}

}

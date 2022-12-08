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

func (rt *_router) getComments(w http.ResponseWriter, r *http.Request, ps httprouter.Params, user structures.User) {
	w.Header().Set("content-type", "application/json")
	var photoId structures.PhotoID
	var comments structures.Comments
	var pageId int64
	var err error

	photoId.Value, err = strconv.ParseInt(ps.ByName("photoId"), 10, 64)
	if err != nil {
		errors.WriteResponse(rt.baseLogger, w, "PhotoId is not valid", http.StatusBadRequest, "PhotoId is not valid")
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

	comments, err = rt.db.GetComments(photoId, pageId, user)
	if err != nil && strings.Contains(err.Error(), "image not found") {
		errors.WriteResponse(rt.baseLogger, w, "File not found", http.StatusNotFound, "Image not found")
		return
	} else if err != nil {
		errors.WriteResponse(rt.baseLogger, w, "Database error", http.StatusInternalServerError, "Internal server error")
		return
	}

	if len(comments.Comments) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(&comments)
	if err != nil {
		errors.WriteResponse(rt.baseLogger, w, "getComments return an error.", http.StatusInternalServerError, "Internal server error")
		return
	}

}

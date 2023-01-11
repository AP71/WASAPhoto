package api

import (
	"net/http"
	"strconv"
	"strings"
	"wasa-photo/service/api/errors"
	"wasa-photo/service/api/structures"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) uncommentPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, user structures.User) {
	w.Header().Set("content-type", "application/json")

	var photoId structures.PhotoID
	var commentId structures.CommentId
	var err error

	photoId.Value, err = strconv.ParseInt(ps.ByName("photoId"), 10, 64)
	if err != nil {
		errors.WriteResponse(rt.baseLogger, w, "Bad request: photoId not valid", http.StatusBadRequest, "Bad request: photoId not valid")
		return
	}

	if ps.ByName("username") != user.Username.Value {
		errors.WriteResponse(rt.baseLogger, w, "Operation not permitted", http.StatusForbidden, "Unauthorized access: Operation not permitted")
		return
	}

	if r.URL.Query().Has("idCommento") {
		commentId.Value, err = strconv.ParseInt(r.URL.Query().Get("idCommento"), 10, 64)
		if err != nil {
			errors.WriteResponse(rt.baseLogger, w, "Incorrect pageId number", http.StatusBadRequest, "Incorrect pageId number")
			return
		}
	} else {
		errors.WriteResponse(rt.baseLogger, w, "Bad request: comment id requested", http.StatusBadRequest, "Bad request:  comment id requested")
		return
	}

	if err != nil {
		errors.WriteResponse(rt.baseLogger, w, "Wrong JSON received", http.StatusBadRequest, "Wrong JSON received")
		return
	}

	err = rt.db.DeleteComment(commentId)
	if err != nil && strings.Contains(err.Error(), "0 changes") {
		errors.WriteResponse(rt.baseLogger, w, "Comment not foun", http.StatusNotFound, "Comment not found")
		return
	} else if err != nil {
		errors.WriteResponse(rt.baseLogger, w, "Database error: "+err.Error(), http.StatusInternalServerError, "Internal server error")
		return
	}

	w.WriteHeader(http.StatusNoContent)

}

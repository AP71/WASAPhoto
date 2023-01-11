package api

import (
	"net/http"
	"strconv"
	"strings"
	"wasa-photo/service/api/errors"
	"wasa-photo/service/api/structures"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) unlikePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, user structures.User) {
	w.Header().Set("content-type", "application/json")

	var photoId structures.PhotoID
	var byUsername string
	var err error

	photoId.Value, err = strconv.ParseInt(ps.ByName("photoId"), 10, 64)
	byUsername = ps.ByName("username")

	if byUsername != user.Username.Value {
		errors.WriteResponse(rt.baseLogger, w, "Operation not permitted", http.StatusForbidden, "Unauthorized access: Operation not permitted")
		return
	}

	if err != nil {
		errors.WriteResponse(rt.baseLogger, w, "Bad request: photoId not valid", http.StatusBadRequest, "Bad request: photoId not valid")
		return
	}

	err = rt.db.RemoveLike(photoId, user)
	if err != nil && strings.Contains(err.Error(), "image not found") {
		errors.WriteResponse(rt.baseLogger, w, "Image not found", http.StatusNotFound, "Image not found")
		return
	} else if err != nil && strings.Contains(err.Error(), "0 changes") {
		errors.WriteResponse(rt.baseLogger, w, "Conflict error", http.StatusConflict, "You never placed a like on this photo")
		return
	} else if err != nil {
		errors.WriteResponse(rt.baseLogger, w, "Database error: "+err.Error(), http.StatusInternalServerError, "Internal server error")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

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

func (rt *_router) commentPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, user structures.User) {
	w.Header().Set("content-type", "application/json")

	var photoId structures.PhotoID
	var comment structures.Comment
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

	err = json.NewDecoder(r.Body).Decode(&comment)
	_ = r.Body.Close()

	if err != nil {
		errors.WriteResponse(rt.baseLogger, w, "Wrong JSON received", http.StatusBadRequest, "Wrong JSON received")
		return
	}

	err = rt.db.WriteComment(photoId, user, comment)
	if err != nil && strings.Contains(err.Error(), "image not found") {
		errors.WriteResponse(rt.baseLogger, w, "Image not found", http.StatusNotFound, "Image not found")
		return
	} else if err != nil && strings.Contains(err.Error(), "UNIQUE constraint failed:") {
		errors.WriteResponse(rt.baseLogger, w, "Conflict error", http.StatusConflict, "You already placed a like on this photo")
		return
	} else if err != nil {
		errors.WriteResponse(rt.baseLogger, w, "Database error", http.StatusInternalServerError, "Internal server error")
		return
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(errors.JSONMsg{Message: "Comment placed successfully"})
	if err != nil {
		errors.WriteResponse(rt.baseLogger, w, "commentPhoto return an error.", http.StatusInternalServerError, "Internal server error")
		return
	}

}

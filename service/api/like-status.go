package api

import (
	"net/http"
	"strconv"
	"wasa-photo/service/api/errors"
	"wasa-photo/service/api/structures"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) likeStatus(w http.ResponseWriter, r *http.Request, ps httprouter.Params, user structures.User) {
	w.Header().Set("content-type", "application/json")

	var photoId structures.PhotoID

	id, err := strconv.ParseInt(ps.ByName("photoId"), 10, 64)
	if err != nil {
		errors.WriteResponse(rt.baseLogger, w, "Error during conversion from string to int64", http.StatusBadRequest, "Bad request: id not valid")
		return
	}
	photoId.Value = id

	if ps.ByName("username") != user.Username.Value {
		errors.WriteResponse(rt.baseLogger, w, "Operation not permitted", http.StatusForbidden, "Unauthorized access: Operation not permitted")
		return
	}

	err = rt.db.GetLikeStatus(user, photoId)
	if err != nil {
		errors.WriteResponse(rt.baseLogger, w, "Interaction not found (ban)", http.StatusNotFound, "Interaction not found.")
		return
	} else {
		w.WriteHeader(http.StatusNoContent)
	}

}

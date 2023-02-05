package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"wasa-photo/service/api/errors"
	"wasa-photo/service/api/structures"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) likeStatus(w http.ResponseWriter, r *http.Request, ps httprouter.Params, user structures.User) {
	w.Header().Set("content-type", "application/json")

	var photoId structures.PhotoID
	var status structures.Status

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
	w.WriteHeader(http.StatusOK)
	if err != nil {
		status.Status = false
	} else {
		status.Status = true
	}
	err = json.NewEncoder(w).Encode(status)
	if err != nil {
		errors.WriteResponse(rt.baseLogger, w, "likeStatus return an error.", http.StatusInternalServerError, "Internal server error")
		return
	}

}

package api

import (
	"net/http"
	"strconv"
	"strings"
	"wasa-photo/service/api/errors"
	"wasa-photo/service/api/structures"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) deletePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, user structures.User) {
	w.Header().Set("content-type", "application/json")

	var file structures.PhotoID

	if ps.ByName("username") != user.Username.Value {
		errors.WriteResponse(rt.baseLogger, w, "Operation not permitted", http.StatusForbidden, "Unauthorized access: Operation not permitted")
		return
	}

	id, err := strconv.ParseInt(ps.ByName("photoId"), 10, 64)
	if err != nil {
		errors.WriteResponse(rt.baseLogger, w, "Error during conversion from string to int64", http.StatusBadRequest, "Bad request: id not valid")
		return
	}
	file.Value = id

	err = rt.db.DeleteFile(file)
	if err != nil && strings.Contains(err.Error(), "sql: no rows in result set") {
		errors.WriteResponse(rt.baseLogger, w, "File not found", http.StatusNotFound, "Image not found")
		return
	} else if err != nil {
		errors.WriteResponse(rt.baseLogger, w, "Database error.", http.StatusInternalServerError, "Internal server error")
		return
	}

	w.WriteHeader(http.StatusNoContent)
	if err != nil {
		errors.WriteResponse(rt.baseLogger, w, "deletePhoto return an error.", http.StatusInternalServerError, "Internal server error")
		return
	}

}

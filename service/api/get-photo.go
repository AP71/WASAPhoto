package api

import (
	"net/http"
	"strconv"
	"strings"
	"wasa-photo/service/api/auth"
	"wasa-photo/service/api/errors"
	"wasa-photo/service/api/structures"

	"github.com/julienschmidt/httprouter"
)

// getHelloWorld is an example of HTTP endpoint that returns "Hello world!" as a plain text
func (rt *_router) getPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var photoId int64
	var image structures.Image

	res, _ := auth.CheckAuth(rt.db, r)

	if !res {
		errors.WriteResponse(rt.baseLogger, w, "Authentication failed", http.StatusUnauthorized, "Unauthorized access")
		return
	}

	photoId, err := strconv.ParseInt(ps.ByName("photoId"), 10, 64)
	if err != nil {
		errors.WriteResponse(rt.baseLogger, w, "Error during conversion from string to int64", http.StatusBadRequest, "Bad request: photoId not valid")
		return
	}

	err = rt.db.GetPhoto(photoId, &image)
	if err != nil && strings.Contains(err.Error(), "sql: no rows in result set") {
		errors.WriteResponse(rt.baseLogger, w, "File not found", http.StatusNotFound, "Image not found")
		return
	} else if err != nil {
		errors.WriteResponse(rt.baseLogger, w, "Database error", http.StatusInternalServerError, "Internal server error")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("content-type", "image/*")
	_, err = w.Write(image.Value)
	if err != nil {
		errors.WriteResponse(rt.baseLogger, w, "getPhoto return an error.", http.StatusInternalServerError, "Internal server error")
		return
	}

}

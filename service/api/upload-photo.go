package api

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"wasa-photo/service/api/errors"
	"wasa-photo/service/api/structures"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) uploadPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, user structures.User) {
	w.Header().Set("content-type", "application/json")

	var file structures.Image

	if ps.ByName("username") != user.Username.Value {
		errors.WriteResponse(rt.baseLogger, w, "Operation not permitted", http.StatusForbidden, "Unauthorized access: Operation not permitted")
		return
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		errors.WriteResponse(rt.baseLogger, w, "Wrong image received", http.StatusBadRequest, "Wrong image received")
		return
	}
	file.Value = data

	if !file.IsValid() {
		errors.WriteResponse(rt.baseLogger, w, "File to large", http.StatusRequestEntityTooLarge, "File to large")
		return
	} else if !strings.HasPrefix(file.Extension(), "image/") {
		errors.WriteResponse(rt.baseLogger, w, "File not supported", http.StatusUnsupportedMediaType, "File not supported")
		return
	}

	err = rt.db.UploadFile(file, user.Id.Value)
	if err != nil {
		errors.WriteResponse(rt.baseLogger, w, "Database error", http.StatusInternalServerError, "Internal server error")
		return
	}

	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(errors.JSONMsg{Message: "Resource created"})
	if err != nil {
		errors.WriteResponse(rt.baseLogger, w, "uploadPhoto return an error.", http.StatusInternalServerError, "Internal server error")
		return
	}

}

package api

import (
	"encoding/json"
	"io"
	"net/http"
	"wasa-photo/service/api/auth"
	"wasa-photo/service/api/errors"
	"wasa-photo/service/api/structures"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) uploadPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	var file structures.Image

	//Getting user data
	res, user := auth.CheckAuth(rt.db, r)

	if !res {
		errors.WriteResponse(rt.baseLogger, w, "Authentication failed", http.StatusUnauthorized, "Unauthorized access")
		return
	}

	if ps.ByName("username") != user.Username.Value {
		errors.WriteResponse(rt.baseLogger, w, "Operation not permitted", http.StatusForbidden, "Unauthorized access: Operation not permitted")
		return
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		errors.WriteResponse(rt.baseLogger, w, "Wrong image received", http.StatusUnsupportedMediaType, "Wrong image received")
		return
	}
	file.Value = data

	err = rt.db.UploadFile(file, user.Id.Value)
	if err != nil {
		errors.WriteResponse(rt.baseLogger, w, "Database error", http.StatusInternalServerError, "Internal server error")
		return
	}

	//returning Message
	err = json.NewEncoder(w).Encode(errors.JSONMsg{Message: "Resource created"})
	if err != nil {
		errors.WriteResponse(rt.baseLogger, w, "uploadPhoto return an error.", http.StatusInternalServerError, "Internal server error")
		return
	}

}

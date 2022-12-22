package errors

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

type JSONMsg struct {
	Message string `json:"message"`
}

func WriteResponse(logger logrus.FieldLogger, w http.ResponseWriter, logMessage string, httpStatus int, message string) {
	w.Header().Set("content-type", "application/json")
	if !strings.Contains(logMessage, "Interaction not found") {
		logger.Warning(logMessage)
	}
	w.WriteHeader(httpStatus)
	_ = json.NewEncoder(w).Encode(JSONMsg{Message: message})
}

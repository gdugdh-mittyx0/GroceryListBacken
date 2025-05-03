package response

import (
	"encoding/json"
	"glbackend/internal/errorsStatus"
	"net/http"
)

type Error struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

func NewError(w http.ResponseWriter, err error, status int) error {
	response := Error{
		Error:   err.Error(),
		Message: errorsStatus.Message(err),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(response)
}

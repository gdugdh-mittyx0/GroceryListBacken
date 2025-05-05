package response

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func NewSuccess(w http.ResponseWriter, result interface{}, status int) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if result != nil {
		return json.NewEncoder(w).Encode(result)
	}
	return nil
}

type SuccessList struct {
	Result interface{} `json:"result"`
	Total  int         `json:"total"`
}

func NewSuccessList(w http.ResponseWriter, result interface{}, total int, status int) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	response := SuccessList{
		Result: result,
		Total:  total,
	}

	if result != nil {
		return json.NewEncoder(w).Encode(response)
	}
	return nil
}

func NewSuccessFile(w http.ResponseWriter, filename string, content []byte) error {
	contentType := http.DetectContentType(content)
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.Header().Set("Content-Length", strconv.Itoa(len(content)))

	_, err := w.Write(content)
	return err
}

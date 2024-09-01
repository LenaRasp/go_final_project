package utils

import (
	"encoding/json"
	"net/http"
)

var TimeLayout = "20060102"

func ResponseError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(code)

	errorMsg := map[string]string{"error": message}

	resp, err := json.Marshal(errorMsg)
	if err != nil {
		return
	}

	w.Write(resp)
}

func ResponseSuccess(w http.ResponseWriter, response any, code int) {
	resp, err := json.Marshal(response)
	if err != nil {
		ResponseError(w, "Ошибка сериализации", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(resp)
}

package response

import (
	"encoding/json"
	"net/http"
)

func Error(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(code)

	errorMsg := map[string]string{"error": message}

	resp, err := json.Marshal(errorMsg)
	if err != nil {
		return
	}

	_ , _ = w.Write(resp)
}

func Success(w http.ResponseWriter, response any, code int) {
	resp, err := json.Marshal(response)
	if err != nil {
		Error(w, "Ошибка сериализации", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ , _ = w.Write(resp)
}

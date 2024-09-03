package handlers

import (
	"net/http"
	"os"
)

func GetStatic(w http.ResponseWriter, req *http.Request) {
	webDir := "web"

	if _, err := os.Stat(webDir + req.RequestURI); os.IsNotExist(err) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else {
		http.FileServer(http.Dir(webDir)).ServeHTTP(w, req)
	}
}

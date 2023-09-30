package router

import (
	"net/http"

	"github.com/Leonardo-Antonio/ms-storage/internal/files"
)

func init() {
	handler := files.New()
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	router.HandleFunc("/v1/directory/structure/users/{userId}", handler.GetDirectoryStructure).Methods(http.MethodGet)
	router.HandleFunc("/v1/files/upload/{userId}", handler.UploadFiles).Methods(http.MethodPost)
}

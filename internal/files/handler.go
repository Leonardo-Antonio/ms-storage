package files

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	utils_files "github.com/Leonardo-Antonio/ms-storage/pkg/files"
	"github.com/Leonardo-Antonio/ms-storage/pkg/req"
	"github.com/Leonardo-Antonio/ms-storage/pkg/response"
	"github.com/gorilla/mux"
)

type handler struct{}

func New() *handler { return &handler{} }

func (h *handler) GetDirectoryStructure(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userId"]

	folderPath := fmt.Sprintf("static/%s", userID)
	root := &utils_files.TreeNode{
		Name:  folderPath,
		IsDir: true,
	}

	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		parts := strings.Split(path, string(filepath.Separator))
		node := root
		url := fmt.Sprintf("%s://%s/static", req.IsTLS(r), r.Host)
		for _, part := range parts[1:] {
			child := utils_files.FindChild(node, part)
			if child == nil {
				if info.IsDir() {
					child = &utils_files.TreeNode{
						Name:  part,
						IsDir: info.IsDir(),
					}
					url += "/" + part
				} else {
					child = &utils_files.TreeNode{
						Name:  url + "/" + part,
						IsDir: info.IsDir(),
					}
				}
				node.Children = append(node.Children, child)
			}
			node = child
			url += "/" + child.Name
		}

		return nil
	})

	if err != nil {
		fmt.Println(err)
		return
	}

	response.Json(w, response.Response{
		Success:   true,
		Data:      root,
		ItemFound: true,
	})
}

func (h *handler) UploadFiles(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userId"]

	// Parsea el formulario con las imágenes
	err := r.ParseMultipartForm(10 << 20) // Establece un límite de 10MB para los archivos
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	// Obtiene los archivos desde el formulario
	files := r.MultipartForm.File["images"]

	for _, file := range files {
		// Abre el archivo
		src, err := file.Open()
		if err != nil {
			http.Error(w, "Error opening file", http.StatusInternalServerError)
			return
		}
		defer src.Close()

		pathSaveFile := fmt.Sprintf("static/%s/", userID)
		if utils_files.IsVideo(file.Filename) {
			pathSaveFile += "videos/" + file.Filename
		} else {
			pathSaveFile += "images/" + file.Filename
		}

		dst, err := os.Create(pathSaveFile)
		if err != nil {
			http.Error(w, "Error creating destination file", http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		// Copia el contenido del archivo al destino
		_, err = io.Copy(dst, src)
		if err != nil {
			http.Error(w, "Error copying file", http.StatusInternalServerError)
			return
		}
	}

	// Responde al cliente con una confirmación
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Images uploaded successfully")
}

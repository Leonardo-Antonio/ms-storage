package files

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	utils_files "github.com/Leonardo-Antonio/ms-storage/pkg/files"
	"github.com/Leonardo-Antonio/ms-storage/pkg/req"
	"github.com/Leonardo-Antonio/ms-storage/pkg/response"
	"github.com/gorilla/mux"
)

type handler struct {
	pathVideo string
	pathImage string
}

func New() *handler {
	return &handler{
		pathVideo: "static/%s/videos/%s",
		pathImage: "static/%s/images/%s",
	}
}

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
	}, http.StatusOK)
}

func (h *handler) UploadFiles(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userId"]

	// Parsea el formulario con las imágenes
	err := r.ParseMultipartForm(300 << 20) // Establece un límite de 10MB para los archivos
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	// Obtiene los archivos desde el formulario
	files := r.MultipartForm.File["images"]

	// Crear directorio base
	baseDir := fmt.Sprintf("static/%s", userID)
	if err := os.MkdirAll(baseDir, os.ModePerm); err != nil {
		log.Println("Error creating base directory:", err)
		http.Error(w, "Error creating base directory", http.StatusInternalServerError)
		return
	}

	// Limitar la cantidad de trabajadores concurrentes
	concurrency := 4 // Puedes ajustar este valor según tus necesidades
	var wg sync.WaitGroup
	workQueue := make(chan *multipart.FileHeader, len(files))

	// Inicializar trabajadores
	for i := 0; i < concurrency; i++ {
		go func() {
			for file := range workQueue {
				processAndSaveFile(file, baseDir)
				wg.Done()
			}
		}()
	}

	// Agregar trabajos al canal de trabajo
	for _, file := range files {
		wg.Add(1)
		workQueue <- file
	}

	// Cerrar el canal de trabajo y esperar a que todos los trabajadores finalicen
	close(workQueue)
	wg.Wait()

	// Responder al cliente con una confirmación
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Images uploaded successfully")
}

func (h *handler) RemoveFiles(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var path string
	if utils_files.IsVideo(vars["name"]) {
		path = fmt.Sprintf(h.pathVideo, vars["userId"], vars["name"])
	} else {
		path = fmt.Sprintf(h.pathImage, vars["userId"], vars["name"])
	}

	if err := os.Remove(path); err != nil {
		response.Json(w, response.Response{
			Success:   false,
			Data:      path,
			ItemFound: false,
			TimeStamp: uint64(time.Now().UnixMilli()),
		}, http.StatusConflict)
		return
	}

	response.Json(w, response.Response{
		Success:   true,
		Data:      path,
		ItemFound: true,
		TimeStamp: uint64(time.Now().UnixMilli()),
	}, http.StatusOK)
}

func processAndSaveFile(file *multipart.FileHeader, baseDir string) {
	// Abre el archivo
	src, err := file.Open()
	if err != nil {
		log.Println("Error opening file:", err)
		return
	}
	defer src.Close()

	// Determina la ruta del archivo
	var fileTypeDir string
	if utils_files.IsVideo(file.Filename) {
		fileTypeDir = "videos"
	} else {
		fileTypeDir = "images"
	}

	ext := filepath.Ext(file.Filename)
	fileName := fmt.Sprintf("%d%s", time.Now().UnixMilli(), ext)
	pathSaveFile := filepath.Join(baseDir, fileTypeDir, fileName)

	// Crea el directorio del tipo de archivo si no existe
	if err := os.MkdirAll(filepath.Join(baseDir, fileTypeDir), os.ModePerm); err != nil {
		log.Println("Error creating file type directory:", err)
		return
	}

	dst, err := os.Create(pathSaveFile)
	if err != nil {
		log.Println("Error creating destination file:", err)
		return
	}
	defer dst.Close()

	// Copia el contenido del archivo al destino
	_, err = io.Copy(dst, src)
	if err != nil {
		log.Println("Error copying file:", err)
		return
	}

	log.Printf("success => %s => upload file => %s\n", fileName, file.Filename)
}

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func main() {
	fmt.Print("Ingresa el user y path: [leo/images]: ")
	var pathname string
	fmt.Scan(&pathname)
	// Especifica la ruta de la carpeta que deseas explorar
	folderPath := "static/" + pathname

	fmt.Printf("Path a modificar => %s", folderPath)

	// Llama a la función para obtener los archivos
	files, err := getFilesInFolder(folderPath)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Imprime los nombres de los archivos
	for _, file := range files {
		time.Sleep(time.Second * 1)
		if err := os.Rename(file, folderPath+fmt.Sprint(time.Now().UnixMilli())+filepath.Ext(file)); err != nil {
			fmt.Printf("error => %s\n", err.Error())
		}
	}
}

func getFilesInFolder(folderPath string) ([]string, error) {
	var files []string

	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}

package files

import "strings"

func FindChild(node *TreeNode, name string) *TreeNode {
	for _, child := range node.Children {
		if child.Name == name {
			return child
		}
	}
	return nil
}

func IsVideo(filename string) bool {
	// Lista de extensiones de archivos de video comunes
	videoExtensions := []string{".mp4", ".avi", ".mov", ".mkv", ".wmv", ".mp3"}

	// Obtener la extensión del archivo
	fileExt := strings.ToLower(filename[strings.LastIndex(filename, "."):])

	// Comprobar si la extensión del archivo está en la lista de extensiones de video
	for _, ext := range videoExtensions {
		if ext == fileExt {
			return true
		}
	}
	return false
}

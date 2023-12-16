package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func uploadHandler(w http.ResponseWriter, r *http.Request) {


	r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile("file")
	if err != nil {
		fmt.Println("Error al obtener el archivo:", err)
		http.Error(w, "Error al obtener el archivo", http.StatusBadRequest)
		return
	}
	defer file.Close()

	fmt.Printf("Archivo recibido: %+v\n", handler.Filename)

	dst, err := os.Create("./uploads/" + handler.Filename)
	if err != nil {
		fmt.Println("Error al crear el archivo en el servidor:", err)
		http.Error(w, "Error al crear el archivo en el servidor", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		fmt.Println("Error al copiar el contenido del archivo:", err)
		http.Error(w, "Error al copiar el contenido del archivo", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Archivo recibido: %s", handler.Filename)
}


func listFilesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	fileList, err := listFiles("./uploads")
	if err != nil {
		fmt.Println("Error al obtener la lista de archivos:", err)
		http.Error(w, "Error al obtener la lista de archivos", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(fileList)
}

func listFiles(directory string) ([]os.DirEntry, error) {
	files, err := os.ReadDir(directory)
	if err != nil {
		return nil, err
	}

	return files, nil
}



func main() {
	os.Mkdir("./uploads", os.ModePerm)

	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/listFiles", listFilesHandler)

	fmt.Println("Servidor escuchando en http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

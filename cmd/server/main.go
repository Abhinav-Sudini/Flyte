package main

import (
	"fmt"
	_ "io"
	_ "log"
	"net/http"
	// "net/http/httputil"
	_ "net/url"
	"os"
	_ "path/filepath"
	_ "strings"

	"Flyte/handlers"
)

const (
	PORT = 8000
)

const (
	BASE_URL = "/"
	FS_URL   = "/api/fs/"
)

const (
	STATIC_FILES_PATH = "./static/"
)

func main() {
	fmt.Println("Now Listening on 80")
	http.HandleFunc(BASE_URL, handlers.ServeFiles)
	http.HandleFunc(FS_URL, handlers.ServeFS)
	http.ListenAndServe(":8000", nil)
}

func OSReadDir(root string) ([]string, error) {
	var files []string
	f, err := os.Open(root)
	if err != nil {
		return files, err
	}
	fileInfo, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		return files, err
	}

	for _, file := range fileInfo {
		files = append(files, file.Name())
	}
	return files, nil
}

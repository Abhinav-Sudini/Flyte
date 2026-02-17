package main

import (
	"fmt"
	_ "io"
	_ "log"
	"net/http"
	"strconv"

	// "net/http/httputil"
	_ "net/url"
	"os"
	_ "path/filepath"
	_ "strings"

	"Flyte/config"
	"Flyte/handlers"
)

var (
	PORT = config.PORT
)

const (
	BASE_URL = "/"
	FS_URL   = "/api/fs/"
)

func main() {
	if len(os.Args)==2{
		p,err := strconv.Atoi(os.Args[1])
		if err == nil {
			PORT = p
		}
	}
	fmt.Println("Starting server")
	fmt.Println("sering on http://localhost:", PORT)

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

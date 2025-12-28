package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"net/url"
	"strings"
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
	http.HandleFunc(BASE_URL, serveFiles)
	http.HandleFunc(FS_URL, serveFS)
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
func handleUpload(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20) //10 MB
	operation := r.FormValue("operation")
	USER_SERVER_FS_ROOT := "./ALL_FILES/" + Int_to_string(Get_user_id(r)) + "/"

	server_upload_dir := USER_SERVER_FS_ROOT + r.FormValue("upload-loc")
	if Dir_not_exist(server_upload_dir) {
		log.Println("error dir does not exist", server_upload_dir)
		http.Error(w, "dir does not exist", http.StatusBadRequest)
		return
	}

	if operation == "new-file" {

		file, handler, err := r.FormFile("file")
		if err != nil {
			log.Println("error retrieving file", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer file.Close()

		dst, err := os.Create(server_upload_dir + handler.Filename)
		if err != nil {
			log.Println("error creating file", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		if _, err := io.Copy(dst, file); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else if operation == "new-folder" {
		os.Mkdir(server_upload_dir+"NewFolder", 0755)
	}
	fmt.Println("done bro")
	fmt.Fprintf(w, "uploaded file")
}

func handleDownload(w http.ResponseWriter, r *http.Request) {
	USER_SERVER_FS_ROOT := "./ALL_FILES/" + Int_to_string(Get_user_id(r)) 

	decoded_path, _ := url.PathUnescape(r.URL.Path)
	path, _ := strings.CutPrefix(decoded_path, FS_URL)
	fmt.Println(path)
	if path == "Home/" {
		json_topo, err := DirToJSON(USER_SERVER_FS_ROOT + "/Home/")
		if err != nil {
			fmt.Println("topo not generated")
			return
		}
		fmt.Fprintf(w, "%s", json_topo)
	} else {
		server_fs_path := USER_SERVER_FS_ROOT +"/"+ path

		if File_not_exist(server_fs_path) {
			log.Println("error dir does not exist", server_fs_path)
			http.Error(w, "dir does not exist", http.StatusBadRequest)
			return
		}
		http.ServeFile(w, r, server_fs_path)
	}

}

func handlePUT(w http.ResponseWriter, r *http.Request) {

}
func handleDelete(w http.ResponseWriter, r *http.Request) {
	// r.ParseMultipartForm(10 << 20) //10 MB
	// operation := r.FormValue("operation")
	// USER_SERVER_FS_ROOT := "./ALL_FILES/" + Int_to_string(Get_user_id(r)) + "/"
	//
	// server_upload_dir := USER_SERVER_FS_ROOT + r.FormValue("upload-loc")
	// if Dir_not_exist(server_upload_dir) {
	// 	log.Println("error dir does not exist", server_upload_dir)
	// 	http.Error(w, "dir does not exist", http.StatusBadRequest)
	// 	return
	// }

}
func serveFS(w http.ResponseWriter, r *http.Request) {
	ShowReq(r)

	if r.Method == "POST" {
		handleUpload(w, r)
	} else if r.Method == "GET" {
		handleDownload(w, r)
	} else if r.Method == "PUT" {
		handlePUT(w, r)
	} else if r.Method == "DELETE" {
		handleDelete(w, r)
	}

}

func serveFiles(w http.ResponseWriter, r *http.Request) {
	dump, _ := httputil.DumpRequest(r, false) // true = include body

	fmt.Println("----- HTTP REQUEST -----")
	fmt.Println(string(dump))
	fmt.Println("------------------------")
	path, _ := strings.CutPrefix(r.URL.Path, "")
	// path := r.URL.Path
	p := STATIC_FILES_PATH + path
	fmt.Println(p)
	// if p == "./" {
	//     p = "./static/index.html"
	// }
	http.ServeFile(w, r, p)
}

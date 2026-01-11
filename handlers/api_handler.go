package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"net/url"
	"strings"

	"Flyte/lib/fs"
	"Flyte/util"
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

func HandleUpload(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20) //10 MB
	operation := r.FormValue("operation")
	USER_SERVER_FS_ROOT := "./ALL_FILES/" + util.Int_to_string(util.Get_user_id(r)) + "/"

	server_upload_dir := USER_SERVER_FS_ROOT + r.FormValue("upload-loc")
	if util.Dir_not_exist(server_upload_dir) {
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

		if handler.Size > (10 << 20) {
			log.Println("file to big - ", handler.Size)
			http.Error(w, "to big of a file", http.StatusBadRequest)
			return
		}

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
	} else if operation == "rename" {
		old_file_name := r.FormValue("old_file_name")
		new_file_name := r.FormValue("new_file_name")
		server_file_path := filepath.Join(server_upload_dir, old_file_name)
		if util.File_not_exist(server_file_path) && util.Dir_not_exist(server_file_path) {
			log.Println("file does not exist")
			http.Error(w, "file does not exist", http.StatusBadRequest)
			return
		}
		err := os.Rename(server_file_path, filepath.Join(server_upload_dir, new_file_name))
		if err != nil {
			log.Println("rename failed -", err)
			return
		}

	} else if operation == "delete" {
		file_name := r.FormValue("file_name")
		server_file_path := filepath.Join(server_upload_dir, file_name)
		if util.File_not_exist(server_file_path) && util.Dir_not_exist(server_file_path) {
			log.Println("file does not exist")
			http.Error(w, "file does not exist", http.StatusBadRequest)
			return
		}
		err := os.RemoveAll(server_file_path)
		if err != nil {
			log.Println("failed to del", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	}
	fmt.Println("done bro")
	fmt.Fprintf(w, "uploaded file")
}

func HandleDownload(w http.ResponseWriter, r *http.Request) {
	USER_SERVER_FS_ROOT := "./ALL_FILES/" + util.Int_to_string(util.Get_user_id(r))

	decoded_path, _ := url.PathUnescape(r.URL.Path)
	path, _ := strings.CutPrefix(decoded_path, FS_URL)
	if path == "Home/" {
		json_topo, err := fs.DirToJSON(USER_SERVER_FS_ROOT + "/Home/")
		if err != nil {
			return
		}
		fmt.Fprintf(w, "%s", json_topo)
	} else {
		server_fs_path := USER_SERVER_FS_ROOT + "/" + path

		if util.File_not_exist(server_fs_path) {
			log.Println("error dir does not exist", server_fs_path)
			http.Error(w, "dir does not exist", http.StatusBadRequest)
			return
		}
		http.ServeFile(w, r, server_fs_path)
	}

}

func HandlePUT(w http.ResponseWriter, r *http.Request) {

}
func HandleDelete(w http.ResponseWriter, r *http.Request) {
	// r.ParseMultipartForm(10 << 20) //10 MB
	// operation := r.FormValue("operation")
	// USER_SERVER_FS_ROOT := "./ALL_FILES/" + util.Int_to_string(Get_user_id(r)) + "/"
	//
	// server_upload_dir := USER_SERVER_FS_ROOT + r.FormValue("upload-loc")
	// if util.Dir_not_exist
	// (server_upload_dir) {
	// 	log.Println("error dir does not exist", server_upload_dir)
	// 	http.Error(w, "dir does not exist", http.StatusBadRequest)
	// 	return
	// }

}
func ServeFS(w http.ResponseWriter, r *http.Request) {
	// ShowReq(r)

	if r.Method == "POST" {
		HandleUpload(w, r)
	} else if r.Method == "GET" {
		HandleDownload(w, r)
	} else if r.Method == "PUT" {
		HandlePUT(w, r)
	} else if r.Method == "DELETE" {
		HandleDelete(w, r)
	}

}

func ServeFiles(w http.ResponseWriter, r *http.Request) {
	// dump, _ := httputil.DumpRequest(r, false) // true = include body

	// fmt.Println("----- HTTP REQUEST -----")
	// fmt.Println(string(dump))
	// fmt.Println("------------------------")
	path, _ := strings.CutPrefix(r.URL.Path, "")
	// path := r.URL.Path
	p := STATIC_FILES_PATH + path
	// fmt.Println(p)
	// if p == "./" {
	//     p = "./static/index.html"
	// }
	http.ServeFile(w, r, p)
}

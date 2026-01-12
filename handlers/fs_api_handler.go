package handlers

import (
	"fmt"
	"log"
	"net/http"

	"net/url"
	"strings"

	"Flyte/lib/auth"
	"Flyte/lib/fs"
	"Flyte/store"
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
	STATIC_FILES_PATH = "./static/html/"
)

func ServeFS(w http.ResponseWriter, r *http.Request) {
	// ShowReq(r)
	if r.Method == "POST" {
		HandlePost(w, r)
	} else if r.Method == "GET" {
		HandleGet(w, r)
	} else if r.Method == "PUT" {
		HandlePUT(w, r)
	} else if r.Method == "DELETE" {
		HandleDelete(w, r)
	}

}

func HandlePost(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20) //10 MB
	upload_directory := r.FormValue("upload-loc")
	operation := r.FormValue("operation")

	if upload_directory == "" || (operation != "new-file" && operation != "new-folder") {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	if operation == "new-folder" {
		err := store.SaveNewFolder(auth.GetUserId(r), upload_directory)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "file could not be created", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}

	new_file, handler, err := r.FormFile("file")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "failed parsing file", http.StatusInternalServerError)
		return
	}
	defer new_file.Close()

	err = store.SaveNewFile(auth.GetUserId(r), upload_directory, new_file, handler)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "unable to save file", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)

	fmt.Println("done bro")
	fmt.Fprintf(w, "uploaded file")
}

func HandleGet(w http.ResponseWriter, r *http.Request) {
	USER_SERVER_FS_ROOT := "./ALL_FILES/" + util.Int_to_string(auth.GetUserId(r))

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
	r.ParseMultipartForm(10 << 20) //10 MB
	upload_directory := r.FormValue("upload-loc")
	operation := r.FormValue("operation")

	if upload_directory == "" || (operation != "rename" ){
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	old_file_name := r.FormValue("old_file_name")
	new_file_name := r.FormValue("new_file_name")
	err := store.RenameFile(auth.GetUserId(r),upload_directory,old_file_name,new_file_name)
	if err != nil {
		log.Println("rename failed -", err)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
func HandleDelete(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20) //10 MB
	upload_directory := r.FormValue("upload-loc")
	operation := r.FormValue("operation")
	file_name := r.FormValue("file_name")
	
	if operation != "delete"{
		http.Error(w,"bad req",http.StatusBadRequest)
		return
	}

	err := store.DeleteFile(auth.GetUserId(r),upload_directory,file_name)
	if err != nil {
		log.Println("failed to del", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func ServeFiles(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	path, _ := strings.CutPrefix(r.URL.Path, "")
	p := STATIC_FILES_PATH + path
	http.ServeFile(w, r, p)
}

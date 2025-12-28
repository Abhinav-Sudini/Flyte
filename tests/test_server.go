package test_server

import (
	"fmt"
	"net/http"
	"strings"
)

func run_server() {
    fmt.Println("Now Listening on 80")
    http.HandleFunc("/api/", serveFiles)
    http.ListenAndServe(":8000", nil)
}

func serveFiles(w http.ResponseWriter, r *http.Request) {
	fmt.Println("new req")
    fmt.Println(r.URL.Path)
    fmt.Println(r.URL.Host)
    fmt.Println(r.URL.RawPath)
    fmt.Println(r.URL.RawQuery)
	path,_ := strings.CutPrefix(r.URL.Path,"/api")
    p := "." + path
    // if p == "./" {
    //     p = "./static/index.html"
    // }
    http.ServeFile(w, r, p)
}

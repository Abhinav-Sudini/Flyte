package util

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"
)

func ShowReq(r *http.Request) {
	dump, _ := httputil.DumpRequest(r, false) // true = include body

	fmt.Println("----- HTTP REQUEST -----")
	fmt.Println(string(dump))
	fmt.Println("------------------------")
}
func Int_to_string(val int) string {
	return fmt.Sprintf("%d", val)
}
func Get_user_id(r *http.Request) int {
	return 0
}
func Dir_not_exist(path string) bool {
	if stat, err := os.Stat(path); err == nil && stat.IsDir() {
		return false
	} else {
		return true
	}
}
func File_not_exist(path string) bool {
	if stat, err := os.Stat(path); err == nil && stat.IsDir() == false {
		return false
	} else {
		return true
	}
}

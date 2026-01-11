package util

import (
	"fmt"
	"net/http"
	"net/http/httputil"
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


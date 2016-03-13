package main

import (
	"mpm/web"
	"net/http"
)

func main() {
	RunWeb()
	//RunCmd()
}

func RunWeb() {
	http.HandleFunc("/", web.API)
	http.HandleFunc("/web", web.WebMgr)
	http.ListenAndServe(":10107", nil)
}

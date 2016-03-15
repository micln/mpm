package main

import (
	"mpm/web"
	"net/http"
)

func main() {

	//	cmd.RunCmd()

	RunWeb()
}

func RunWeb() {
	http.HandleFunc("/", web.API)
	http.HandleFunc("/web", web.WebMgr)
	http.ListenAndServe(":10107", nil)
}

package main

import (
	//. "mpm/passer"

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
	http.HandleFunc("/sync", web.SyncQr)

	http.ListenAndServe(":10107", nil)
}

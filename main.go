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
	http.HandleFunc("/", web.Public)
	http.HandleFunc("/web", web.WebMgr)
	http.HandleFunc("/sync", web.SyncQr)

	http.ListenAndServe(":10107", nil)
}

package web

import (
	"encoding/json"
	"log"
	"mpm/passer"
	"net/http"
	"text/template"
)

var P *passer.Passer

func init() {
	P = passer.Pr
}

func API(w http.ResponseWriter, r *http.Request) {

	action := r.URL.Query().Get("action")
	query := r.URL.Query().Get("label")
	switch action {
	case "get":
		result := P.Get(query)
		b, _ := json.Marshal(result)
		w.Write(b)
	case "gen":
		P.Gen(query)
		result := P.Get(query)
		b, _ := json.Marshal(result)
		w.Write(b)
	}
}

func WebMgr(w http.ResponseWriter, r *http.Request) {

	data := make(map[string]interface{})

	data[`list`] = P.Get()

	tp, e := template.ParseFiles(`web/web.html`)
	if e != nil {
		log.Println(e)
		return
	}

	e = tp.Execute(w, data)
	if e != nil {
		log.Println(e)
		return
	}
}

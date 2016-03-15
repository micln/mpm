package web

import (
	"encoding/json"
	"log"
	"mpm/passer"
	"net/http"
	"strconv"
	"text/template"
)

var P *passer.PManager

func init() {
	P = passer.Pr
}

func API(w http.ResponseWriter, r *http.Request) {
	action := r.URL.Query().Get("action")
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	label := r.URL.Query().Get("label")
	password := r.URL.Query().Get("password")
	note := r.URL.Query().Get("note")

	switch action {
	case "get":
	case "gen":
		p := &passer.Password{
			Id:       id,
			Label:    label,
			Password: password,
			Note:     note,
		}
		p.Save()
	case "del":
		passer.Password{Id: id}.Remove()
	}

	result := P.Get(label)
	b, _ := json.Marshal(result)
	w.Write(b)
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

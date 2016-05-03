package web

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"mpm/passer"
	"net/http"
	"strconv"
	"text/template"
	"time"
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

	responseHTML(w, `web/web.html`, data)

	return

}

func SyncQr(w http.ResponseWriter, r *http.Request) {
	salt := sha1.Sum([]byte(fmt.Sprint(time.Now().Nanosecond(), random())))
	log.Println(salt)
	responseHTML(w, `web/qr.html`, map[string]string{})
	return
}

func Public(w http.ResponseWriter, r *http.Request) {
	b, e := ioutil.ReadFile(`web/public` + r.URL.String())
	if e != nil {
		w.Write([]byte(e.Error()))
	}

	w.Write(b)
}

func responseHTML(w http.ResponseWriter, filename string, data interface{}) {

	tp, e := template.ParseFiles(filename)
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

func random() int64 {
	rand.Seed(int64(time.Now().Nanosecond()))
	return rand.Int63()
}

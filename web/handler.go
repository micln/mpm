package web

import (
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"mpm/pm"
	"net/http"
	"strconv"
	"text/template"
	"time"

	"github.com/chanxuehong/wechat/json"
)

var P *pm.PManager

func init() {
	P = pm.Pr
}

func API(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL)

	action := r.URL.Query().Get("action")
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	title := r.URL.Query().Get("title")
	site := r.URL.Query().Get("site")
	account := r.URL.Query().Get("account")
	password := r.URL.Query().Get("password")
	note := r.URL.Query().Get("note")

	switch action {
	case "get":
	case "gen":
		p := &pm.Password{
			Id:       id,
			Title:    title,
			Site:     site,
			Account:  account,
			Password: password,
			Note:     note,
		}
		p.Save()
	case "del":
		pm.Password{Id: id}.Remove()
	}

	//	Filter
	query := r.URL.Query().Get("q")
	origins := P.Get(query)
	results := []interface{}{}
	for i := range origins {
		results = append(results, origins[i].GetFake())
	}
	b, _ := json.Marshal(results)
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

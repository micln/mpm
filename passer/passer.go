package passer

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

const DB_FILE_NAME = `mp.dbp`

var Pr *Passer

func init() {
	Pr = New()
}

type Passer struct {
	Skey string
	List map[string]string
}

func New() *Passer {
	passer := Passer{}
	passer.List = make(map[string]string)
	passer.Load()
	return &passer
}

func (this *Passer) GetSkey() string {
	return `4589968`
}

func (this *Passer) SetSkey(newSK string) error {
	//	@todo 需要修改所有的密码
	return nil
}

func (this *Passer) Put(k, v string) {
	this.List[k] = v
	this.Save()
}

func (this *Passer) Get(greps ...string) map[string]string {
	result := make(map[string]string)

	grep := ""
	if len(greps) > 0 {
		grep = greps[0]
	}

	for k, v := range this.List {
		if len(grep) > 0 && !strings.Contains(k, grep) {
			continue
		}
		result[k] = v
	}
	return result
}

func (this *Passer) GetToCLI(greps ...string) {
	result := this.Get(greps...)
	for k, v := range result {
		fmt.Printf("%-15s: %s\n", k, v)
	}
}

func randomPassowrd() string {
	rand.Seed(time.Now().UnixNano())
	all := `0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ`
	s := ``
	for i := 1; i <= 12; i++ {
		j := rand.Intn(len(all))
		s += all[j : j+1]
		if i%3 == 0 && i < 12 {
			s += `-`
		}
	}
	return s
}

func (this *Passer) Gen(label string) {
	pwd := randomPassowrd()
	this.Put(label, pwd)
	this.Save()
}

func (this *Passer) Save() {
	list := make(map[string][]byte)
	for k, v := range this.List {
		list[k] = this.encode(v)
	}
	b, _ := json.Marshal(list)
	ioutil.WriteFile(DB_FILE_NAME, b, os.ModePerm)
}

func (this *Passer) Load() {
	b, e := ioutil.ReadFile(DB_FILE_NAME)

	//	首次运行
	if os.IsNotExist(e) {
		this.List = make(map[string]string)
		return
	}

	if e != nil {
		log.Fatalln(`读文件错误`, e)
	}

	//	没有数据了
	if len(b) < 3 {
		this.List = make(map[string]string)
		return
	}
	list := make(map[string][]byte)
	e = json.Unmarshal(b, &list)
	if e != nil {
		log.Fatalln(`解析数据出错`, e)
	}

	for k, v := range list {
		this.List[k] = this.decode(v)
	}
}

func (this *Passer) encode(pwd string) []byte {
	s, _ := RsaEncrypt([]byte(this.GetSkey() + pwd))
	return s
}

func (this *Passer) decode(pwd []byte) string {
	sb, e := RsaDecrypt(pwd)
	if e != nil {
		log.Fatalln(e)
	}
	s := string(sb)

	if strings.Index(s, this.GetSkey()) == 0 {
		return s[len(this.GetSkey()):]
	}
	return ""
}

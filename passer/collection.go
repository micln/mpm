package passer

import (
	"encoding/base64"
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

var Pr *PManager

func init() {
	Pr = NewPasser()
}

type PManager struct {
	Skey  string
	List  map[string]string
	Datas []*Password

	nextId int
}

type Password struct {
	Id       int
	Label    string
	Password string
	Note     string
}

func NewPasser() *PManager {
	passer := PManager{}
	passer.List = make(map[string]string)
	passer.LoadFromDB()
	return &passer
}

func GetOne(id int) *Password {
	return Pr.Datas[id]
}

func (p *Password) Save() {

	if len(p.Password) < 1 {
		p.Password = randomPassowrd()
	}

	if p.Exist() {
		//	新建
		for i := range Pr.Datas {
			if Pr.Datas[i].Label == p.Label {
				p.Id = Pr.NextId()
				Pr.Datas[i] = p
				break
			}
		}
	} else {
		p.Id = Pr.NextId()
		Pr.Datas = append(Pr.Datas, p)
	}
	go Pr.SaveToDB()
}

func (p Password) Remove() {
	for i := range Pr.Datas {
		if Pr.Datas[i].Id == p.Id {
			Pr.Datas = append(Pr.Datas[:i], Pr.Datas[i+1:]...)
			return
		}
	}
}

func (this *Password) serialize() []byte {
	m := make(map[string]interface{})
	m["id"] = this.Id
	m["label"] = this.Label
	m["note"] = this.Note
	m["password"] = encodePwd(this.Password)
	b, _ := json.Marshal(m)
	return b
}

func (this *Password) unSerialize(b []byte) {
	m := make(map[string]interface{})
	err := json.Unmarshal(b, &m)
	if err != nil {
		log.Fatalln(`数据加载错误：`, err)
	}
	this.Id = int(m["id"].(float64))
	this.Label = m["label"].(string)
	this.Note = m["note"].(string)

	decoded, _ := base64.StdEncoding.DecodeString(m["password"].(string))
	this.Password = decodePwd(decoded)
}

func (p *Password) Exist() bool {
	return Pr.Exist(p)
}

func (this *PManager) GetSkey() string {
	return ``
}

func (this *PManager) Len() int {
	return len(this.Datas)
}

func (this *PManager) NextId() int {
	this.nextId += 1
	return this.nextId
}

//	所有判重使用Label属性
func (this *PManager) Exist(p *Password) bool {
	for i := range this.Datas {
		if this.Datas[i].Label == p.Label {
			return true
		}
	}
	return false
}

//func (this *Passer) Add(label, pwd, note string) *Password {
//
//	if this.Exist(label) {
//		return nil
//	}
//
//	p := &Password{
//		Id:       this.NextId(),
//		Label:    label,
//		Password: pwd,
//		Note:     note,
//	}
//
//	this.Datas = append(this.Datas, p)
//
//	return p
//}

func (this *PManager) Get(greps ...string) (result []*Password) {

	grep := ""
	if len(greps) > 0 {
		grep = greps[0]
	}

	for _, v := range this.Datas {
		if len(grep) > 0 && !strings.Contains(v.Label, grep) {
			continue
		}
		result = append(result, v)
	}

	return result
}

func (this *PManager) GetToCLI(greps ...string) {
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

func (this *PManager) Gen(label, note string) {
	pwd := randomPassowrd()
	p := &Password{
		Label:    label,
		Password: pwd,
		Note:     note,
	}
	p.Save()
}

func (this *PManager) SaveToDB() {
	var ss []interface{}
	for _, v := range this.Datas {
		ss = append(ss, v.serialize())
	}
	b, _ := json.Marshal(ss)
	ioutil.WriteFile(DB_FILE_NAME, b, os.ModePerm)
}

func (this *PManager) LoadFromDB() {
	b, e := ioutil.ReadFile(DB_FILE_NAME)

	//	首次运行
	if os.IsNotExist(e) {
		log.Println(`未找到数据库文件`)
		return
	}

	if e != nil {
		log.Fatalln(`数据库读取错误：`, e)
	}

	//	没有数据了
	if len(b) < 3 {
		log.Println(`初始化数据库`)
		return
	}

	list := []interface{}{}
	e = json.Unmarshal(b, &list)
	if e != nil {
		log.Fatalln(`解析数据出错`, e)
	}

	for _, bi := range list {
		p := &Password{}
		b, _ := base64.StdEncoding.DecodeString(bi.(string))
		p.unSerialize(b)
		this.Datas = append(this.Datas, p)

		//	next Id
		if p.Id > this.nextId {
			this.nextId = p.Id
		}
	}
}

func encodePwd(pwd string) []byte {
	s, _ := RsaEncrypt([]byte(pwd))
	return s
}

func decodePwd(pwd []byte) string {
	sb, e := RsaDecrypt(pwd)
	if e != nil {
		log.Fatalln(e)
	}
	return string(sb)
}

func dump(v ...interface{}) {
	for i := range v {
		fmt.Println(strings.Repeat(`-`, 50))
		fmt.Println(v[i])
	}
}

func dd(v ...interface{}) {
	dump(v)
	os.Exit(1)
}

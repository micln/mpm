package pm

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

const DB_FILE_NAME = `mp.dbp`

var Pr *PManager

func init() {
	Pr = NewPasser()
	Pr.startClearUp()
}

type PManager struct {
	Skey  string
	datas [512]*Password

	lastId int
}

func NewPasser() *PManager {
	provider := PManager{}
	provider.LoadFromDB()
	return &provider
}

func GetOne(id int) *Password {
	return Pr.datas[id]
}

func (this *PManager) GetSkey() string {
	return ``
}

func (this *PManager) Datas() (result []*Password) {
	for i := range this.datas {
		if this.datas[i] != nil {
			result = append(result, this.datas[i])
		}
	}
	return
}

func (this *PManager) NextId() int {
	this.lastId += 1
	return this.lastId
}

func (this *PManager) Exist(p *Password) bool {

	return p.Id > 0

	for i := range this.datas {
		if this.datas[i].Site == p.Site && this.datas[i].Account == p.Account {
			return true
		}
	}
	return false
}

func (this *PManager) Get(greps ...string) (result []*Password) {

	grep := ""
	if len(greps) > 0 {
		grep = greps[0]
	}

	for _, v := range this.Datas() {
		if v != nil &&
			len(grep) > 0 &&
			!strings.Contains(v.Title, grep) &&
			!strings.Contains(v.Site, grep) &&
			!strings.Contains(v.Account, grep) {
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

func (this *PManager) Gen(label, note string) {
	pwd := randomPassword()
	p := &Password{
		Title:    label,
		Password: pwd,
		Note:     note,
	}
	p.Save()
}

func (this *PManager) Remove(id int) {
	this.datas[id] = nil
	this.SaveToDB()
}

func (this *PManager) SaveToDB() {
	for _, v := range this.Datas() {
		v.Serialize()
	}
	b, _ := json.Marshal(this.Datas())
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

	for _, m := range list {
		ms, _ := json.Marshal(m)
		p := &Password{}
		p.UnSerialize(ms)
		this.datas[p.Id] = p

		//	next Id
		if p.Id > this.lastId {
			this.lastId = p.Id
		}
	}
}

func (this *PManager) startClearUp() {
	go func() {
		for {
			time.Sleep(1 * time.Hour)

			l, r := 1, len(this.datas)-1
			for l < r {
				for l < r && this.datas[l] != nil {
					l++
				}

				for l < r && this.datas[r] == nil {
					r--
				}

				if l < r {
					p := this.datas[l]
					this.datas[l] = this.datas[r]
					this.datas[r] = p
					this.datas[l].Id = l
				}

			}
		}
	}()
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

package pm

import (
	"encoding/json"
	"log"
	"math/rand"
	"time"
)

type Password struct {
	Id            int
	Title         string //	类目
	Site          string //	网址
	Account       string //	账号名
	Password      string `json:"-"` // 明文密码
	HashdPassword []byte //	加密密码,用以持久化
	Note          string //	其他说明
}

func (p *Password) Save() {

	if len(p.Password) < 1 {
		p.Password = randomPassword()
	}

	if p.Id == 0 {
		p.Id = Pr.NextId()
	}
	Pr.datas[p.Id] = p

	go Pr.SaveToDB()
}

func (p Password) Remove() {
	Pr.Remove(p.Id)
}

func (this *Password) Serialize() []byte {
	this.HashdPassword = encodePwd(this.Password)
	b, _ := json.Marshal(*this)
	return b
}

func (this *Password) UnSerialize(b []byte) {

	err := json.Unmarshal(b, this)
	if err != nil {
		log.Fatalln(`数据加载错误：`, err)
	}
	this.Password = decodePwd(this.HashdPassword)

	//m := make(map[string]interface{})
	//err := json.Unmarshal(b, &m)
	//if err != nil {
	//	log.Fatalln(`数据加载错误：`, err)
	//}
	//this.Id = int(m["id"].(float64))
	//this.Label = m["label"].(string)
	//this.Note = m["note"].(string)
	//
	//decoded, _ := base64.StdEncoding.DecodeString(m["password"].(string))
	//this.Password = decodePwd(decoded)
}

//	渲染成alfred可识别的数据
func (p *Password) GetFake() map[string]interface{} {
	v := make(map[string]interface{})
	v[`title`] = p.Account
	v[`addon`] = p.Title + ` @` + p.Site
	v[`data`] = p.Password
	return v
}

func (p *Password) Exist() bool {
	return Pr.Exist(p)
}

func randomPassword() string {
	rand.Seed(time.Now().UnixNano())
	all := `0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ,@.!<>`
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

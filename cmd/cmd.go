package cmd

import (
	"fmt"
	"io/ioutil"
	"mpm/passer"
	"os"
)

var P *passer.PManager

func init() {
	P = passer.Pr
}

func RunCmd() {

	if len(os.Args) < 2 {
		showHelp()
		return
	}

	action := os.Args[1]
	args := []string{}

	if len(os.Args) > 2 {
		args = os.Args[2:]
	}

	switch action {
	case `l`, `ls`, `list`:
		grep := ""
		if len(args) > 0 {
			grep = args[0]
		}
		P.GetToCLI(grep)

	case `add`:
		P.Put(args[0], args[1])
		P.GetToCLI()

	default:
		fmt.Println(`参数错误`)
	}

}

func showHelp() {
	b, e := ioutil.ReadFile(`help.txt`)
	if e != nil {
		b = []byte(`帮助文件不存在`)
	}
	fmt.Printf(string(b))
}

package main

import (
	"fmt"

	"github.com/valyala/fasthttp"
)

func main() {
	//Run("11","111910707")
	//Run("1", "801617312")
	//run("12", "9317174")
	getCfg()
	run(cfg.CType, cfg.Oid)
	fmt.Println("*********************************")
	fmt.Printf("去重整合后，共%d人\n", len(set))
	for k, v := range set {
		fmt.Printf("评论者uid: %s 评论者昵称: %s 评论内容: %s\n", k.Mid, k.Uname, v)
	}
}

func demo() {
	url := `http://httpbin.org/get`

	status, resp, err := fasthttp.Get(nil, url)
	if err != nil {
		fmt.Println("请求失败:", err.Error())
		return
	}

	if status != fasthttp.StatusOK {
		fmt.Println("请求没有成功:", status)
		return
	}

	fmt.Println(string(resp))
}

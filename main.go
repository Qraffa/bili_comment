package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	//Run("11","111910707")
	//Run("1", "801617312")
	//run("12", "9317174")
	getCfg()
	run(cfg.CType, cfg.Oid)
	fmt.Println("*********************************")
	fmt.Printf("去重整合后，共%d人\n", len(set))

	printInfo()

	fmt.Println("DONE!!!")

}

func printInfo() {
	f, err := os.OpenFile("out.txt", os.O_APPEND, 0777)
	if err != nil {
		log.Fatal(fmt.Sprintf("open file error ===> %s\n Create the file <out.txt> and try again.", err.Error()))
	}
	buf := bufio.NewWriter(f)
	for k, v := range set {
		buf.WriteString(fmt.Sprintf("评论者uid: %s 评论者昵称: %s 评论内容: %s\n", k.Mid, k.Uname, v))
	}
	buf.Flush()
}
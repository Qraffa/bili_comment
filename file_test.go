package main

import (
	"bili_comment/config"
	"fmt"
	"os"
	"testing"
	"time"
)

func Test_File(t *testing.T) {
	f, err := os.Create("out/comment.txt")
	if err != nil {
		fmt.Println(err)
	}
	f.WriteString("qwe")
	f.WriteString("asd")
	f.Close()
}

func Test_Cfg(t *testing.T) {
	cfg := config.Cfg()
	fmt.Println(cfg)
}

func Test_Out(t *testing.T) {
	//fmt.Printf("uid: %-20s name: %-60s content: %s\n","31760933","长长长长长长长","内容")
	//fmt.Printf("uid: %-20s name: %-60s content: %s\n","186610111","短短短","内容")
	fmt.Printf("out_%d\n", time.Now().Unix())
}

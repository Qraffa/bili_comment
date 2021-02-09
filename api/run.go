package api

import (
	"bili_comment/types"
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"
)

var (
	runWg sync.WaitGroup
)

func Run(hostUid string, page int) {
	dir := mkdir()
	drun := NewDynamicRunner()
	cards := drun.RunDynamic(hostUid, page)
	// 遍历动态列表，获取每个动态的评论区
	for _, v := range cards {
		runWg.Add(1)
		oid := v.Desc.Rid
		// 对于type为 1 | 4 的动态，需要使用DynamicId作为评论区oid
		if v.Desc.DType == 1 || v.Desc.DType == 4 {
			oid = v.Desc.DynamicId
		}
		go func(t, i int) {
			defer runWg.Done()
			crun := NewCommentRunner()
			cType, id := strconv.Itoa(t), strconv.Itoa(i)
			commentRes := crun.RunCommentFromDynamic(cType, id)
			printfComment(cType, id, commentRes, dir)
		}(v.Desc.DType, oid)
	}
	printDynamic(hostUid, cards, dir)
	runWg.Wait()
}

func mkdir() string {
	dir := fmt.Sprintf("out_%d", time.Now().Unix())
	err := os.Mkdir(dir, os.ModePerm)
	if err != nil {
		log.Fatal("create out dir failed\n")
	}
	return dir
}

// 打印评论列表
func printfComment(cType, oid string, set map[types.Member]string, dir string) {
	fileName := fmt.Sprintf("%s/comment_%s_%s.txt", dir, cType, oid)
	f, err := os.Create(fileName)
	if err != nil {
		log.Fatal(fmt.Sprintf("open file error ===> %s\n Create the file <out.txt> and try again.", err.Error()))
	}
	buf := bufio.NewWriter(f)
	for k, v := range set {
		buf.WriteString(fmt.Sprintf("评论者uid: %-20s 评论者昵称: %-20s 评论内容: %-s\n", k.Mid, k.Uname, v))
	}
	buf.Flush()
}

// 打印动态列表
func printDynamic(hostUid string, cards []*types.Card, dir string) {
	fileName := fmt.Sprintf("%s/dynamic_%s.txt", dir, hostUid)
	f, err := os.Create(fileName)
	if err != nil {
		log.Fatal(fmt.Sprintf("open file error ===> %s\n Create the file <out.txt> and try again.", err.Error()))
	}
	buf := bufio.NewWriter(f)
	for _, v := range cards {
		buf.WriteString(fmt.Sprintf("动态类型: %-2d 动态评论区oid: %-20d 动态内容: %s\n", v.Desc.DType, v.Desc.Rid, v.Card))
	}
	buf.Flush()
}

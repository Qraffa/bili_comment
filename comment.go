package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"sync"

	"github.com/valyala/fasthttp"
)

const (
	pageSize = 20
	url      = "http://api.bilibili.com/x/v2/reply"
	nohot    = "1"
	sort     = "1" // TODO 按时间默认排序，数据减少？？
)

var (
	comments    chan *reply
	commentList []*reply
	wg          sync.WaitGroup
	rootNeed    int
	rootGet     int
	set         map[member]string
)

func init() {
	comments = make(chan *reply)
	commentList = make([]*reply, 0)
	set = make(map[member]string)
}

// 查询根评论数
func getPage(cType, oid string) int {
	req := fasthttp.AcquireRequest()
	res := fasthttp.AcquireResponse()

	defer func() {
		fasthttp.ReleaseRequest(req)
		fasthttp.ReleaseResponse(res)
	}()
	req.SetRequestURI(url)
	query := req.URI().QueryArgs()
	query.Add("type", cType)
	query.Add("oid", oid)
	query.Add("nohot", nohot)
	query.Add("sort", sort)
	query.Add("pn", "1")
	query.Add("ps", "1")

	err := fasthttp.Do(req, res)
	if err != nil {
		log.Printf("fasthttp get page error ==> %s\n", err.Error())
		return -1
	}
	replies := &reply{}
	body := res.Body()
	err = json.Unmarshal(body, replies)
	if err != nil {
		log.Printf("json unmarshal page error ==> %s\n", err.Error())
		return -1
	}
	if replies.Code != 0 {
		log.Printf("bili quest comment page error. type: %s, oid: %s, errCode: %d, errMsg: %s\n", cType, oid, replies.Code, replies.Message)
		return -1
	}
	log.Printf("type: %s, oid: %s  ==> get comment ROOTcount: %d, count: %d\n", cType, oid, replies.Data.Page.Count, replies.Data.Page.Acount)
	rootNeed = replies.Data.Page.Count
	return replies.Data.Page.Count
}

// 获取单页评论
func getComment(cType, oid, pageNo string) {
	req := fasthttp.AcquireRequest()
	res := fasthttp.AcquireResponse()

	defer func() {
		fasthttp.ReleaseRequest(req)
		fasthttp.ReleaseResponse(res)
		wg.Done()
	}()
	req.SetRequestURI(url)
	query := req.URI().QueryArgs()
	query.Add("type", cType)
	query.Add("oid", oid)
	query.Add("nohot", nohot)
	query.Add("sort", sort)
	query.Add("pn", pageNo)
	query.Add("ps", strconv.Itoa(pageSize))

	err := fasthttp.Do(req, res)
	if err != nil {
		log.Printf("fasthttp get comment type: %s, oid: %s, pageNo: %s, error ==> %s\n", cType, oid, pageNo, err.Error())
		return
	}
	replies := &reply{}
	body := res.Body()
	err = json.Unmarshal(body, replies)
	if err != nil {
		log.Printf("json unmarshal comment type: %s, oid: %s, pageNo: %s, error ==> %s\n", cType, oid, pageNo, err.Error())
		return
	}
	if replies.Code != 0 {
		log.Printf("bili quest comment error. type: %s, oid: %s, errCode: %d, errMsg: %s\n", cType, oid, replies.Code, replies.Message)
		return
	}
	log.Printf("type: %s, oid: %s, pageSize: %d ==> get comment page %s, get number: %d ===> OK!\n", cType, oid, pageSize, pageNo, len(replies.Data.Replies))
	comments <- replies
}

// 分页查询全部评论
func getAllComment(cType, oid string) {
	rootCount := getPage(cType, oid)
	page := rootCount / pageSize
	if rootCount%pageSize > 0 {
		page++
	}
	for i := 1; i <= page; i++ {
		wg.Add(1)
		page := strconv.Itoa(i)
		go getComment(cType, oid, page)
	}
	wg.Wait()
	close(comments)
}

// 统计所有评论
func countComment() {
	for {
		r, ok := <-comments
		if !ok {
			break
		}
		commentList = append(commentList, r)
	}
	log.Println("get ALL comment ===> OK!")
	// 整合去重
	for _, v := range commentList {
		for _, v := range v.Data.Replies {
			rootGet++
			set[v.Member] = v.Content.Message
		}
	}
}

func run(cType, oid string)  {
	go getAllComment(cType, oid)
	countComment()
	if rootGet != rootNeed {
		log.Printf("GET COMMENT ERR! get: %d, need: %d\n", rootGet, rootNeed)
	}
	log.Println("GET COMMENT DONE!")
}

func check() {
	fmt.Println(comments)
	fmt.Println(commentList)
}

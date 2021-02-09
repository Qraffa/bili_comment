package api

import (
	"bili_comment/types"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"sync"

	"github.com/valyala/fasthttp"
)

var (
	// dynamic类型转换为评论区类型
	// TODO 部分
	dtype2C = map[string]string{
		"64": "12",
		"2":  "11",
		"1":  "17",
		"8":  "1",
		"4":  "17",
	}
)

type CommentRunner struct {
	comments    chan *types.Reply
	commentList []*types.Reply
	wg          sync.WaitGroup
	rootNeed    int
	rootGet     int
	set         map[types.Member]string
	pageSize    int
	url         string
	nohot       string
	sort        string // TODO 按时间默认排序，数据减少？？
}

func NewCommentRunner() *CommentRunner {
	return &CommentRunner{
		comments:    make(chan *types.Reply),
		commentList: make([]*types.Reply, 0),
		wg:          sync.WaitGroup{},
		rootNeed:    0,
		rootGet:     0,
		set:         make(map[types.Member]string),
		pageSize:    20,
		url:         "http://api.bilibili.com/x/v2/reply",
		nohot:       "1",
		sort:        "1",
	}
}

// 查询根评论数
func (c *CommentRunner) getPage(cType, oid string) int {
	req := fasthttp.AcquireRequest()
	res := fasthttp.AcquireResponse()

	defer func() {
		fasthttp.ReleaseRequest(req)
		fasthttp.ReleaseResponse(res)
	}()
	req.SetRequestURI(c.url)
	query := req.URI().QueryArgs()
	query.Add("type", cType)
	query.Add("oid", oid)
	query.Add("nohot", c.nohot)
	query.Add("sort", c.sort)
	query.Add("pn", "1")
	query.Add("ps", "1")

	err := fasthttp.Do(req, res)
	if err != nil {
		log.Printf("fasthttp get page error ==> %s\n", err.Error())
		return -1
	}
	replies := &types.Reply{}
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
	c.rootNeed = replies.Data.Page.Count
	return replies.Data.Page.Count
}

// 获取单页评论
func (c *CommentRunner) getComment(cType, oid, pageNo string) {
	req := fasthttp.AcquireRequest()
	res := fasthttp.AcquireResponse()

	defer func() {
		fasthttp.ReleaseRequest(req)
		fasthttp.ReleaseResponse(res)
		c.wg.Done()
	}()
	req.SetRequestURI(c.url)
	query := req.URI().QueryArgs()
	query.Add("type", cType)
	query.Add("oid", oid)
	query.Add("nohot", c.nohot)
	query.Add("sort", c.sort)
	query.Add("pn", pageNo)
	query.Add("ps", strconv.Itoa(c.pageSize))

	err := fasthttp.Do(req, res)
	if err != nil {
		log.Printf("fasthttp get comment type: %s, oid: %s, pageNo: %s, error ==> %s\n", cType, oid, pageNo, err.Error())
		return
	}
	replies := &types.Reply{}
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
	log.Printf("type: %s, oid: %s, pageSize: %d ==> get comment page %s, get number: %d ===> OK!\n", cType, oid, c.pageSize, pageNo, len(replies.Data.Replies))
	c.comments <- replies
}

// 分页查询全部评论
func (c *CommentRunner) getAllComment(cType, oid string) {
	rootCount := c.getPage(cType, oid)
	page := rootCount / c.pageSize
	if rootCount%c.pageSize > 0 {
		page++
	}
	for i := 1; i <= page; i++ {
		c.wg.Add(1)
		page := strconv.Itoa(i)
		go c.getComment(cType, oid, page)
	}
	c.wg.Wait()
	close(c.comments)
}

// 统计所有评论
func (c *CommentRunner) countComment() {
	for {
		r, ok := <-c.comments
		if !ok {
			break
		}
		c.commentList = append(c.commentList, r)
	}
	log.Println("get ALL comment ===> OK!")
	// 整合去重
	for _, v := range c.commentList {
		for _, v := range v.Data.Replies {
			c.rootGet++
			c.set[v.Member] = v.Content.Message
		}
	}
}

// 获取评论区评论
// 参数：评论区类型，评论区Id
func (c *CommentRunner) RunComment(cType, oid string) map[types.Member]string {
	go c.getAllComment(cType, oid)
	c.countComment()
	if c.rootGet != c.rootNeed {
		log.Printf("GET COMMENT ERR! get: %d, need: %d\n", c.rootGet, c.rootNeed)
	}
	log.Println("GET COMMENT DONE!")
	return c.set
}

// 获取动态的评论区评论
// 参数：动态类型，评论区Id
func (c *CommentRunner) RunCommentFromDynamic(dType, oid string) map[types.Member]string {
	if cType, ok := dtype2C[dType]; !ok {
		log.Printf("Type Not Found!")
		return nil
	} else {
		return c.RunComment(cType, oid)
	}
}

func (c *CommentRunner) check() {
	fmt.Println(c.comments)
	fmt.Println(c.commentList)
}

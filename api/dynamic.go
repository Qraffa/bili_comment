package api

import (
	"bili_comment/types"
	"encoding/json"
	"log"
	"strconv"

	"github.com/valyala/fasthttp"
)

type DynamicRunner struct {
	url          string
	dynamicCards []*types.Card
}

func NewDynamicRunner() *DynamicRunner {
	return &DynamicRunner{
		url:          "https://api.vc.bilibili.com/dynamic_svr/v1/dynamic_svr/space_history",
		dynamicCards: make([]*types.Card, 0),
	}
}

// 获取部分动态列表
// 参数：需要获取的用户的uid，动态列表偏移
// 返回：该用户动态列表的下一次偏移
func (d *DynamicRunner) getDynamic(hostUid, offsetDynamicId string, page int) string {
	req := fasthttp.AcquireRequest()
	res := fasthttp.AcquireResponse()

	defer func() {
		fasthttp.ReleaseRequest(req)
		fasthttp.ReleaseResponse(res)
	}()

	req.SetRequestURI(d.url)
	query := req.URI().QueryArgs()
	query.Add("host_uid", hostUid)
	query.Add("offset_dynamic_id", offsetDynamicId)

	err := fasthttp.Do(req, res)
	if err != nil {
		log.Printf("fasthttp get dynamic error ==> %s\n", err.Error())
		return "0"
	}
	dynamic := &types.Dynamic{}
	body := res.Body()
	err = json.Unmarshal(body, dynamic)
	if err != nil {
		log.Printf("json unmarshal dynamic list error ==> %s\n", err.Error())
		return "0"
	}
	if dynamic.Code != 0 {
		log.Printf("bili quest dynamic list error. uid: %s, offset: %s, errCode: %d, errMsg: %s\n", hostUid, offsetDynamicId, dynamic.Code, dynamic.Message)
		return "0"
	}
	for k := range dynamic.Data.Cards {
		d.dynamicCards = append(d.dynamicCards, &dynamic.Data.Cards[k])
	}
	log.Printf("bili get dynamic UID: %s Page: %d ===> OK!\n", hostUid, page)
	return strconv.Itoa(dynamic.Data.NextOffset)
}

// 获取page页数的动态
func (d *DynamicRunner) getDynamicList(hostUid string, page int) {
	offset := "0"
	// TODO failure retry
	for i := 0; i < page; i++ {
		offset = d.getDynamic(hostUid, offset, page)
	}
}

func (d *DynamicRunner) RunDynamic(hostUid string, page int) []*types.Card {
	d.getDynamicList(hostUid, page)
	log.Printf("bili get all dynamic UID:%s ALL Page: %d OK!\n", hostUid, page)
	return d.dynamicCards
}

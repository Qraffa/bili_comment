package types

type Dynamic struct {
	Code    int         `json:"code"`    // 返回码
	Message string      `json:"message"` // 错误信息
	Data    DynamicData `json:"data"`
}

type DynamicData struct {
	HasMore    int    `json:"has_more"`
	Cards      []Card `json:"cards"`       // 动态列表
	NextOffset int    `json:"next_offset"` // 用于请求下一页动态的偏移
}

type Card struct {
	Desc Desc   `json:"desc"` // 动态信息
	Card string `json:"card"` // 动态内容
}

type Desc struct {
	Uid       int   `json:"uid"`        // 动态发布者uid
	DType     int   `json:"type"`       // 动态类型 ！需要重新对应
	Rid       int   `json:"rid"`        // 动态评论区oid
	DynamicId int   `json:"dynamic_id"` // 用于偏移动态列表 ！对于type为 1 | 4 的动态，需要使用该字段作为评论区oid
	Timestamp int64 `json:"timestamp"`  // 动态发布时间戳
}

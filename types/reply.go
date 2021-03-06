package types

type Reply struct {
	Code    int    `json:"code"`    // 返回码
	Message string `json:"message"` // 错误信息
	Data    Data   `json:"data"`
}

type Data struct {
	Page    Page      `json:"page"`
	Replies []Replies `json:"replies"`
}

type Page struct {
	Num    int `json:"num"`    // 当前页码
	Size   int `json:"size"`   // 每页项数
	Count  int `json:"count"`  // 根评论条数
	Acount int `json:"acount"` // 总计评论条数
}

type Replies struct {
	Oid       int     `json:"oid"`       // 评论区Id
	Type      int     `json:"types"`      // 评论区类型
	Mid       int     `json:"mid"`       // 发送者uid
	FansGrade int     `json:"fansgrade"` // 粉丝标签 0-无，1-有  TODO ？好像数据不对？
	Ctime     int64   `json:"ctime"`     // 评论发送时间戳
	Member    Member  `json:"member"`    // 评论发送者信息
	Content   Content `json:"content"`   // 评论内容
}

type Member struct {
	Mid   string `json:"mid"`   // 评论者uid
	Uname string `json:"uname"` // 评论者昵称
}

type Content struct {
	Message string `json:"message"` // 评论内容
	Plat    int    `json:"plat"`    // 评论发送端 1-web 2-android 3-ios 4-wp
}

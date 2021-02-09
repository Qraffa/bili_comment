`config.yml`配置文件中，uid为目标用户uid，page为需要的动态列表页数

### 使用：
1. 在配置文件中配置好uid和page，run main.go
2. 输出结果保存在`out_xxx`文件夹中，`xxx`对应时间戳

---

最麻烦的点是需要手动去获取评论区的type和oid。

一般是有一个`interaction-status?type=xxx&oid=xxx`的请求，其中的type和oid就是需要的

获取到之后需要手动写到`config.yml`文件的对应键上。

---

F12->Network

在左上方的Filter中输入`interaction-status?type=`一般是可以直接过滤出来的

### 注意

- 从动态列表请求回来的动态card中，type字段与请求评论区时的参数需要转换一下
- 部分类型的动态card评论区oid使用的是`dynamic_id`字段，目前已知type=1|4
- 部分类型的动态card评论区oid使用的是`rid`字段，默认使用该字段
- 如果在请求评论的日志中出现`-404` `啥都没有`等信息，一般可能是需要使用将oid从`rid`字段，改为使用`dynamic_id`字段

| res type   | req type |
| ---------- | -------- |
| 64 专栏    | 12       |
| 2 图片     | 11       |
| 1 分享动态 | 17       |
| 8 视频动态 | 1        |
| 4 动态     | 17       |

### Bug

- 在默认排序方式下部分页码的评论会少于预期。不是偶发。
- 在分页查询情况下，日志输出变成空行？但数据似乎是正常。偶发。
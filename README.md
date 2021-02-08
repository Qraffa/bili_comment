最麻烦的点是需要手动去获取评论区的type和oid。

一般是有一个`interaction-status?type=xxx&oid=xxx`的请求，其中的type和oid就是需要的

获取到之后需要手动写到`config.yml`文件的对应键上。

---

F12->Network

在左上方的Filter中输入`interaction-status?type=`一般是可以直接过滤出来的

### Bug

- 在默认排序方式下部分页码的评论会少于预期？不是偶发。
- 在分页查询情况下，日志输出变成空行？但数据似乎是正常偶发。
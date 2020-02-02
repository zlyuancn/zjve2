# 简单的json值提取器

---

# 获得

`go get -u github.com/zlyuancn/zjve2`

# 示例

# 文档
[godoc](https://godoc.org/github.com/zlyuancn/zjve2)

```
jve := zjve2.LoadString(`
    {
      "a": null,
      "b": true,
      "c": [{"a": "aaa"}, {"b": "bbb"}],
      "d": {"i": 4},
      "e": "hello"
    }`)
fmt.Println(jve.Get("a")) // 获取值
fmt.Println(jve.Get("b"))
fmt.Println(jve.Get("c.#")) // 获取列表的数量
fmt.Println(jve.Get("c.0")) // 获取列表中指定索引的数据
fmt.Println(jve.Get("c.-1")) // 获取列表中倒数第一条数据
fmt.Println(jve.Get("c.1.b"))
fmt.Println(jve.Get("d.i"))
fmt.Println(jve.Get("e"))
```

# 输出结果

```
a: null
b: true
c.#: 2
c.0: {
   "a": "aaa"
}
c.-1: {
   "b": "bbb"
}
c.1.b: "bbb"
d.i: 4
e: "hello"
```
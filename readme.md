# 简单的json值提取器

---

# 获得

`go get -u github.com/zlyuancn/zjve2`

# 文档
[godoc](https://godoc.org/github.com/zlyuancn/zjve2)

# 示例

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
fmt.Println(jve.Get("c.#"))   // 获取列表的数量
fmt.Println(jve.Get("c.[#]")) // 获取列表的数量
fmt.Println(jve.Get("c.0"))   // 获取列表中指定索引的数据
fmt.Println(jve.Get("c.-1"))  // 获取列表中倒数第一条数据
fmt.Println(jve.Get("c.1.b"))
fmt.Println(jve.Get("c.[1]"))   // 获取列表数据另一种语法
fmt.Println(jve.Get("c.[1,2]")) // 获取切片
fmt.Println(jve.Get("c.[1,-1]"))
fmt.Println(jve.Get("d.i"))
fmt.Println(jve.Get("e"))
```

# 输出结果

```
a: null
b: true
c.#: 2
c.[#]: 2
c.0: {
   "a": "aaa"
}
c.-1: {
   "b": "bbb"
}
c.1.b: "bbb"
c.[1]: {
   "b": "bbb"
}
c.[1,2]: err: 索引超出最大数量2
c.[1,-1]: []
d.i: 4
e: "hello"
```

# 舒适的加载方式

```
// 从文件中加载
func LoadFile(filename string) *JVE

// 从Reader中加载
func LoadReader(r io.Reader) *JVE

// 从字符串中加载
func LoadString(s string) *JVE

// 从bytes中加载
func Load(bs [] byte) *JVE
```

# 其他获取结果方法

```
// 获取原始值
func (m *JVE) Raw() interface{}

// 获取错误
func (m *JVE) Err() error

// 获取它的String值, 只有String类型有效
func (m *JVE) Str()

// 获取它的Boolean值, 只有Boolean类型有效
func (m *JVE) Bool()

// 获取它的Float64值, 只有Number类型有效
func (m *JVE) Float64()

// 获取它的Int值, 只有Number类型有效
func (m *JVE) Int()

// 获取它的数量, 只有Array类型有效
func (m *JVE) Count()

// 获取数组指定索引的值, 只有Array类型有效
func (m *JVE) Index(i int)

// 获取数组的切片, 只有Array类型有效
func (m *JVE) Slice(start, end int)

// 判断是否存在某个路径
func (m *JVE) Has(path string)
```

# 将字符串结果转为任何类型值

```
s := jve.MustStr()
s.String()
s.Val()
s.Bytes()
s.Bool()
s.Int()
s.Int8()
s.Int16()
s.Int32()
s.Int64()
s.Uint()
s.Uint8()
s.Uint16()
s.Uint32()
s.Uint64()
s.Float32()
s.Float64()
# 你也可以使用 Scan 将值扫描到任何对象上
```



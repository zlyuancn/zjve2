/*
-------------------------------------------------
   Author :       Zhang Fan
   date：         2020/1/14
   Description :
-------------------------------------------------
*/

package zjve2

import (
    "errors"
    "fmt"
    "io"
    "io/ioutil"
    "os"
    "strconv"
    "strings"

    "github.com/zlyuancn/zstr"
)

var SyntaxErr = errors.New("语法错误")

type JVE struct {
    path string
    t    Type
    raw  interface{}
}

func (m *JVE) makeJve(path string, raw interface{}) *JVE {
    newpath := path
    if m.path != "" {
        newpath = fmt.Sprintf("%s.%s", m.path, path)
    }
    out := &JVE{path: newpath, raw: raw}
    switch raw.(type) {
    case nil:
        out.t = Null
    case error:
        out.t = Error
    case bool:
        out.t = Boolean
    case float64:
        out.t = Number
    case string:
        out.t = String
    case []interface{}:
        out.t = Array
    case map[string]interface{}:
        out.t = Object
    default:
        return m.makeJve(path, fmt.Errorf("未知的类型 %t", raw))
    }
    return out
}

// 获取路径
func (m *JVE) Path() string {
    return m.path
}

// 获取它的类型
func (m *JVE) Type() Type {
    return m.t
}

// -------  获取  ----------

// 根据路径获取结果
func (m *JVE) Get(path string) *JVE {
    if path == "" {
        return m.makeJve("", errors.New("不能有任何空路径"))
    }

    paths := strings.Split(path, ".")

    var out *JVE
    path = paths[0]
    switch m.t {
    case Null:
        return m.makeJve(path, errors.New("Null类型没有子路径了"))
    case Error:
        return m
    case Boolean:
        return m.makeJve(path, errors.New("Boolean类型没有子路径了"))
    case Number:
        return m.makeJve(path, errors.New("Number类型没有子路径了"))
    case String:
        return m.makeJve(path, errors.New("String类型没有子路径了"))
    case Array:
        rawpath := path
        if path[0] == '[' {
            if len(path) < 3 || path[len(path)-1] != ']' {
                return m.makeJve(path, SyntaxErr)
            }
            path = path[1 : len(path)-1]
            rawpath = fmt.Sprintf("[%s]", path)
        }

        if path == "#" {
            count := len(m.raw.([]interface{}))
            return m.makeJve(rawpath, float64(count))
        }

        count := len(m.raw.([]interface{}))

        vs := strings.Split(path, ",")
        if len(vs) > 2 {
            return m.makeJve(rawpath, SyntaxErr)
        }
        array := len(vs) == 2

        var start int
        start, err := strconv.Atoi(vs[0])
        if err != nil {
            return m.makeJve(rawpath, SyntaxErr)
        }
        if start < 0 {
            start += count
        }
        if start < 0 || start >= count {
            return m.makeJve(rawpath, fmt.Errorf("索引超出最大数量%d", count))
        }

        if array {
            end, err := strconv.Atoi(vs[1])
            if err != nil {
                return m.makeJve(rawpath, SyntaxErr)
            }
            if end < 0 {
                end += count
            }
            if end < 0 || end > count {
                return m.makeJve(rawpath, fmt.Errorf("索引超出最大数量%d", count))
            }
            if start > end {
                return m.makeJve(rawpath, errors.New("切片结果不能为负"))
            }
            out = m.makeJve(rawpath, m.raw.([]interface{})[start:end])
        } else {
            out = m.makeJve(rawpath, m.raw.([]interface{})[start])
        }
    case Object:
        v, ok := m.raw.(map[string]interface{})[path]
        if !ok {
            return m.makeJve(path, errors.New("不存在的路径"))
        }
        out = m.makeJve(path, v)
    }

    if len(paths) > 1 {
        return out.Get(strings.Join(paths[1:], "."))
    }

    return out
}

// 获取原始值
// Deprecated: 建议使用 Raw, 因为很容易将Val理解为获取字符串值
func (m *JVE) Val() interface{} {
    return m.raw
}

// 获取原始值
func (m *JVE) Raw() interface{} {
    return m.raw
}

// 获取错误
func (m *JVE) Err() error {
    if m.t != Error {
        return nil
    }
    return m.raw.(error)
}

// 获取它的String值, 只有String类型有效
func (m *JVE) Str() (*zstr.String, error) {
    if m.t != String {
        return nil, fmt.Errorf("需要String, 但它是%s", m.t)
    }
    return zstr.New(m.raw.(string)), nil
}

// 必须转为String值, 否则会panic
func (m *JVE) MustStr() *zstr.String {
    if m.t != String {
        panic(fmt.Errorf("需要String, 但它是%s", m.t))
    }
    return zstr.New(m.raw.(string))
}

// 获取它的Boolean值, 只有Boolean类型有效
func (m *JVE) Bool() (bool, error) {
    if m.t != Boolean {
        return false, fmt.Errorf("需要Boolean, 但它是%s", m.t)
    }
    return m.raw.(bool), nil
}

// 获取它的Float64值, 只有Number类型有效
func (m *JVE) Float64() (float64, error) {
    if m.t != Number {
        return 0, fmt.Errorf("需要Float64, 但它是%s", m.t)
    }
    return m.raw.(float64), nil
}

// 获取它的Int值, 只有Number类型有效
func (m *JVE) Int() (int, error) {
    if m.t != Number {
        return 0, fmt.Errorf("需要Int, 但它是%s", m.t)
    }
    return int(m.raw.(float64)), nil
}

// 获取它的数量, 只有Array类型有效
func (m *JVE) Count() (int, error) {
    if m.t != Array {
        return 0, fmt.Errorf("需要Array, 但它是%s", m.t)
    }
    return len(m.raw.([]interface{})), nil
}

// 获取数组指定索引的值, 只有Array类型有效
func (m *JVE) Index(i int) *JVE {
    return m.Get(strconv.Itoa(i))
}

// 获取数组的切片, 只有Array类型有效
func (m *JVE) Slice(start, end int) *JVE {
    return m.Get(fmt.Sprintf("%d,%d", start, end))
}

// 判断是否存在某个路径
func (m *JVE) Has(path string) bool {
    return m.Get(path).Type() != Error
}

// -------  展示  ----------

// 返回用于展示的字符串
func (m *JVE) String() string {
    return fmt.Sprintf("%s: %s", m.path, m.ToString())
}

// 返回json标准类型格式数据, 注意json没有Error类型, 所以它不能转换为json标准类型格式
func (m *JVE) ToString() string {
    switch m.t {
    case Null:
        return "null"
    case Error:
        return fmt.Sprintf("err: %s", m.raw.(error).Error())
    case Boolean:
        b := m.raw.(bool)
        if b {
            return "true"
        }
        return "false"
    case Number:
        return fmt.Sprintf("%g", m.raw)
    case String:
        return fmt.Sprintf("%q", m.raw)
    case Array, Object:
        out, err := JsonFormatObj(&m.raw, "")
        if err != nil {
            return fmt.Sprintf("不能转换为string: %s", err)
        }
        return out
    }
    return fmt.Sprintf("未知的类型 %d", int(m.t))
}

// -------  加载  ----------

// 从文件中加载
func LoadFile(filename string) *JVE {
    f, err := os.Open(filename)
    if err != nil {
        return &JVE{
            path: "",
            t:    Error,
            raw:  err,
        }
    }
    return LoadReader(f)
}

// 从Reader中加载
func LoadReader(r io.Reader) *JVE {
    bs, err := ioutil.ReadAll(r)
    if err != nil {
        return &JVE{
            path: "",
            t:    Error,
            raw:  err,
        }
    }
    return Load(bs)
}

// 从字符串中加载
func LoadString(s string) *JVE {
    return Load([]byte(s))
}

// 从bytes中加载
func Load(bs [] byte) *JVE {
    jve := new(JVE)
    err := Json.Unmarshal(bs, &jve.raw)
    if err != nil {
        return &JVE{
            path: "",
            t:    Error,
            raw:  err,
        }
    }
    switch jve.raw.(type) {
    case []interface{}:
        jve.t = Array
    case map[string]interface{}:
        jve.t = Object
    default:
        return &JVE{
            path: "",
            t:    Error,
            raw:  fmt.Errorf("未能识别类型: %t", jve.raw),
        }
    }
    return jve
}

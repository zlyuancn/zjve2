/*
-------------------------------------------------
   Author :       Zhang Fan
   date：         2020/1/14
   Description :
-------------------------------------------------
*/

package zjve2

import (
    "bytes"
    "encoding/json"
    "unsafe"

    "github.com/json-iterator/go"
)

var Json = jsoniter.ConfigCompatibleWithStandardLibrary

// 默认换行间隔
var DefaultFormatIndent = "   "

// 将json文本格式化
func JsonFormat(s, indent string) (string, error) {
    return JsonFormatBytes(*(*[]byte)(unsafe.Pointer(&s)), indent)
}

// 将json字节数组格式化
func JsonFormatBytes(bs []byte, indent string) (string, error) {
    if indent == "" {
        indent = DefaultFormatIndent
    }

    var out bytes.Buffer
    // jsoniter当前版本没有Indent
    err := json.Indent(&out, bs, "", indent)
    return out.String(), err
}

// 将对象格式化为json字符串
func JsonFormatObj(v interface{}, indent string) (string, error) {
    if indent == "" {
        indent = DefaultFormatIndent
    }

    // 现在还是使用官方的format, jsoniter当前版本展示方式不是很友好, 或许它以后会解决这个问题?
    bs, err := json.MarshalIndent(v, "", indent)
    return string(bs), err
}

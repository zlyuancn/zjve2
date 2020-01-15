/*
-------------------------------------------------
   Author :       Zhang Fan
   date：         2020/1/14
   Description :
-------------------------------------------------
*/

package zjve2

type Type int

const (
    Null    = Type(iota)
    Error   // 错误类型
    Boolean // 布尔
    Number  // 数字
    String  // 字符串
    Array   // 数组
    Object  // 对象
)

func (m Type) String() string {
    switch m {
    case Error:
        return "Error"
    case Boolean:
        return "Boolean"
    case Number:
        return "Number"
    case String:
        return "String"
    case Array:
        return "Array"
    case Object:
        return "Object"
    }
    return "Null"
}

package kdebug

import (
	"fmt"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/kcolor"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/kstring"

	"github.com/tidwall/pretty"
)

// PrintObject
//  @Description 打印输出对象到控制台
//  @Param message 要打印输出的消息
//  @Param obj 要打印输出的消息体
func PrintObject(message string, obj interface{}) {
	if !IsDevelopmentMode() {
		return
	}
	fmt.Printf("%s => %s\n",
		kcolor.Red(message),
		pretty.Color(
			pretty.Pretty([]byte(kstring.PrettyJson(obj))),
			pretty.TerminalStyle,
		),
	)
}

// DebugBytes
//  @Description interface转换成调试消息
//  @Param obj 要转换的interface
//  @Return 调试的消息
func DebugBytes(obj interface{}) string {
	return string(pretty.Color(pretty.Pretty([]byte(kstring.Json(obj))), pretty.TerminalStyle))
}

// PrintKV
//  @Description 打印输出KV
//  @Param k 要打印输出的Key
//  @Param v 要打印输出的Val
func PrintKV(k string, v string) {
	if !IsDevelopmentMode() {
		return
	}
	fmt.Printf("%-50s => %s\n", kcolor.Red(k), kcolor.Green(v))
}

// PrintKVWithPrefix
//  @Description 打印带前缀输出KV
//  @Param prefix 要打印输出的前缀
//  @Param k 要打印输出的Key
//  @Param v 要打印输出的Val
func PrintKVWithPrefix(prefix string, k string, v string) {
	if !IsDevelopmentMode() {
		return
	}
	fmt.Printf("%-8s]> %-30s => %s\n", prefix, kcolor.Red(k), kcolor.Blue(v))
}

// PrintMap
//  @Description 打印带输出Map
//  @Param m 要打印输出的Map
func PrintMap(m map[string]interface{}) {
	if !IsDevelopmentMode() {
		return
	}
	for key, val := range m {
		fmt.Printf("%-20s : %s\n", kcolor.Red(key), fmt.Sprintf("%+v", val))
	}
}

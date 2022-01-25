package kcolor

import (
	"fmt"
	"math/rand"
	"strconv"
)

// 前景 背景 颜色
// ---------------------------------------
// 30  40  黑色
// 31  41  红色
// 32  42  绿色
// 33  43  黄色
// 34  44  蓝色
// 35  45  紫红色
// 36  46  青蓝色
// 37  47  白色
//
// 模式代码 意义
// -------------------------
//  0  终端默认设置
//  1  高亮显示
//  4  使用下划线
//  5  闪烁
//  7  反白显示
//  8  不可见
// 常规前景色
const (
	TextBlack = iota + 30
	TextRed
	TextGreen
	TextYellow
	TextBlue
	TextMagenta
	TextCyan
	TextWhite
	TextDefault = 39
)

// 扩展前景色 foreground color 90 - 97(非标准)
const (
	TextDarkGray = iota + 90 // 亮黑（灰）
	TextLightRed
	TextLightGreen
	TextLightYellow
	TextLightBlue
	TextLightMagenta
	TextLightCyan
	TextLightWhite
)
const (
	BackBlack = iota + 40
	BackRed
	BackGreen
	BackYellow
	BackBlue
	BackMagenta
	BackCyan
	BackWhite
)

// RandomColor 十六进制随机颜色
//  @Description
//  @Return string
func RandomColor() string {
	return fmt.Sprintf("#%s", strconv.FormatInt(int64(rand.Intn(16777216)), 16))
}

// Yellow
//  @Description 黄色
//  @Param msg 输出文本
//  @Return string
func Yellow(msg string) string {
	//return fmt.Sprintf("\x1b[33m%s\x1b[0m", msg)
	return SetColor(msg, 4, 0, TextYellow)
}

// Red
//  @Description 红色
//  @Param msg 输出文本
//  @Return string
func Red(msg string) string {
	//return fmt.Sprintf("\x1b[31m%s\x1b[0m", msg)
	return SetColor(msg, 5, 0, TextRed)
}

// Redf
//  @Description 红色（格式化）
//  @Param msg 输出文本
//  @Param arg 参数值
//  @Return string
func Redf(msg string, arg interface{}) string {
	return fmt.Sprintf("\x1b[31m%s\x1b[0m %+v\n", msg, arg)
}

// Blue
//  @Description 蓝色
//  @Param msg 输出文本
//  @Return string
func Blue(msg string) string {
	//return fmt.Sprintf("\x1b[34m%s\x1b[0m", msg)
	return SetColor(msg, 0, 0, TextBlue)
}

// Green
//  @Description 绿色
//  @Param msg  输出文本
//  @Return string
func Green(msg string) string {
	//return fmt.Sprintf("\x1b[32m%s\x1b[0m", msg)
	return SetColor(msg, 0, 0, TextGreen)
}

// Greenf
//  @Description 绿色
//  @Param msg 输出文本
//  @Param arg 扩展输出参数
//  @Return string
func Greenf(msg string, arg interface{}) string {
	return fmt.Sprintf("\x1b[32m%s\x1b[0m %+v\n", msg, arg)
}

// Magenta
//  @Description 紫红色
//  @Param msg
//  @Return string
func Magenta(msg string) string {
	return SetColor(msg, 0, 0, TextMagenta)
}

// Cyan
//  @Description 青蓝色
//  @Param msg
//  @Return string
func Cyan(msg string) string {
	return SetColor(msg, 0, 0, TextCyan)
}

// White
//  @Description 白色
//  @Param msg
//  @Return string
func White(msg string) string {
	return SetColor(msg, 0, 0, TextGreen)
}

// SetColor
//  @Description 颜色设置 其中0x1B是标记
//  @Param msg 输出文本
//  @Param mode 配置模式 详见上方注释
//  @Param backColor 背景色 详见上方注释
//  @Param frontColor 前景色 详见上方注释
//  @Return string
func SetColor(msg string, mode, backColor, frontColor int) string {
	return fmt.Sprintf("%c[%d;%d;%dm%s%c[0m", 0x1B, mode, backColor, frontColor, msg, 0x1B)
}

// Info
//  @Description
//  @Param msg
//  @Return string
func Info(msg string) string {
	return SetColor(msg, 0, 0, TextCyan)
}

// Warn
//  @Description 警告样式
//  @Param msg
//  @Return string
func Warn(msg string) string {
	return SetColor(msg, 0, 0, TextYellow)
}

// Error
//  @Description 错误样式
//  @Param msg
//  @Return string
func Error(msg string) string {
	return SetColor(msg, 2, BackRed, TextBlack)
}

// Danger
//  @Description 危险样式
//  @Param msg
//  @Return string
func Danger(msg string) string {
	return SetColor(msg, 0, 0, TextRed)
}

// Debug
//  @Description debug样式
//  @Param msg
//  @Return string
func Debug(msg string) string {
	return SetColor(msg, 0, 0, TextCyan)
}

// Success
//  @Description 成功样式
//  @Param msg
//  @Return string
func Success(msg string) string {
	return SetColor(msg, 0, 0, TextGreen)
}

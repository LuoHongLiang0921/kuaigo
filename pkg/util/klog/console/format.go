package console

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"go.uber.org/zap/zapcore"
)

type Format string

const (
	percent = "%"
	left    = "{"
	endTag  = "}"
)
const (
	// TypeKeyword 关键字 % { }
	TypeKeyword = 1
	// TypeID 标识符
	TypeID = 2
	// TypeLiteral 字面量，支持数字、字母
	TypeLiteral = 3
)

const (
	// StateInit 初始状态
	StateInit = iota
	// StateStartTag1 %
	StateStartTag1
	// StateStartTag %{
	StateStartTag
	// StateID %{id}
	StateID
	StateLiteral
	// StateEndTag }
	StateEndTag
)

var ()

// Token token
type Token struct {
	// Type token 类型，关键字、标识符、字面量
	Type int
	// Text token 值
	Text string
}

// Formatter 格式化者
type Formatter struct {
	Format
	// 临时变量
	tokenText    strings.Builder // 临时保存token 的文本
	tokens       []Token         //保存解析出来的Token
	currentToken Token           //当前正在解析的Token
	cacheTokens  []Token
}

// NewFormatter
// 	@Description
//	@Param format
// 	@Return *Formatter
func NewFormatter(format Format) *Formatter {
	return &Formatter{Format: format}
}

// initToken
func (f *Formatter) initToken(char string) int {
	if f.tokenText.String() != "" {
		f.currentToken.Text = f.tokenText.String()
		f.tokens = append(f.tokens, f.currentToken)
		var newTokenText strings.Builder
		f.tokenText = newTokenText
	}

	newState := StateInit
	if char == percent {
		newState = StateStartTag1
		f.currentToken.Type = TypeKeyword
		f.tokenText.WriteString(char)
	} else if isLiteral(char) {
		newState = StateLiteral
		f.currentToken.Type = TypeLiteral
		f.tokenText.WriteString(char)
	} else if char == endTag {
		newState = StateEndTag
		f.currentToken.Type = TypeKeyword
		f.tokenText.WriteString(char)
	} else {
		newState = StateInit
	}
	return newState
}

func isAlpha(char string) bool {
	return char >= "a" && char <= "z" || char >= "A" && char <= "Z"
}

func isDigit(char string) bool {
	return char >= "0" && char <= "9"
}

func isBlank(char string) bool {
	return char == " " || char == "\t" || char == "\n"
}

func isLiteral(char string) bool {
	return isAlpha(char) || isDigit(char) || char == "=" || isBlank(char)
}

// Tokenize
// 	@Description 解析表达式中所有token
// 	@Receiver f 表达式
// 	@Return []Token token 列表
func (f *Formatter) Tokenize() []Token {
	state := StateInit
	reader := strings.NewReader(string(f.Format))
	for {
		ch, _, err := reader.ReadRune()
		if err == io.EOF {
			break
		}
		char := string(ch)
		switch state {
		case StateInit:
			state = f.initToken(char)
		case StateStartTag1:
			if char == left {
				f.currentToken.Type = TypeKeyword
				f.tokenText.WriteString(char)
				f.currentToken.Text = f.tokenText.String()
				f.tokens = append(f.tokens, f.currentToken)
				var newTokenText strings.Builder
				f.tokenText = newTokenText
				state = StateID
			}
		//	此状态目前看着没有用,待删除
		case StateStartTag:
			if isLiteral(char) {
				f.tokenText.WriteString(char)
				f.currentToken.Type = TypeID
				state = StateID
			} else {
				state = f.initToken(char)
			}
		case StateID:
			if char == endTag {
				state = f.initToken(char)
			} else if isLiteral(char) {
				f.tokenText.WriteString(char)
				f.currentToken.Type = TypeID
			}
		case StateLiteral:
			if isLiteral(char) {
				f.tokenText.WriteString(char)
			} else {
				state = f.initToken(char)
			}
		case StateEndTag:
			if isLiteral(char) {
				state = f.initToken(char)
			}
		}
	}
	return f.tokens
}

// String
// 	@Description 返回表达式的值，内置 + fields 显示，
// 	@Receiver f 格式串 msg 值
//	@Param fields 字段值
// 	@Return string
func (f *Formatter) String(fields ...zapcore.Field) string {
	//todo: tokens 缓存
	tokens := f.Tokenize()
	var (
		out strings.Builder
		in  = make(map[string]bool, len(fields))
	)

	// tokens
	for _, v := range tokens {
		switch v.Type {
		case TypeLiteral:
			out.WriteString(v.Text)
		case TypeID:
			id := v.Text
			idVal := `""`
			for _, field := range fields {
				if field.Key == id {
					in[id] = true
					switch field.Type {
					case zapcore.Int64Type:
						idVal = strconv.FormatInt(field.Integer, 10)
					case zapcore.StringType:
						idVal = field.String
						if idVal == "" {
							idVal = `""`
						}
					default:
						fmt.Println(field.Interface)
					}
				}
			}
			out.WriteString(idVal)
		}
	}
	// 剩余的fields
	for _, f := range fields {
		key := f.Key
		if !in[key] {
			idVal := `""`
			switch f.Type {
			case zapcore.Int64Type:
				idVal = strconv.FormatInt(f.Integer, 10)
			case zapcore.StringType:
				idVal = f.String
				if idVal == "" {
					idVal = `""`
				}
			default:
				fmt.Println(f.Interface)
			}
			out.WriteString(" " + key + "=" + idVal)
		}
	}
	return out.String()
}

// @Description

package kvalidator

import (
	"encoding/json"
	"net"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	Email   string = "^(((([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+(\\.([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+)*)|((\\x22)((((\\x20|\\x09)*(\\x0d\\x0a))?(\\x20|\\x09)+)?(([\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x7f]|\\x21|[\\x23-\\x5b]|[\\x5d-\\x7e]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(\\([\\x01-\\x09\\x0b\\x0c\\x0d-\\x7f]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}]))))*(((\\x20|\\x09)*(\\x0d\\x0a))?(\\x20|\\x09)+)?(\\x22)))@((([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])([a-zA-Z]|\\d|-|\\.|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.)+(([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])([a-zA-Z]|\\d|-|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.?$"
	Int     string = "^(?:[-+]?(?:0|[1-9][0-9]*))$"
	Float   string = "^(?:[-+]?(?:[0-9]+))?(?:\\.[0-9]*)?(?:[eE][\\+\\-]?(?:[0-9]+))?$"
	Base64  string = "^(?:[A-Za-z0-9+\\/]{4})*(?:[A-Za-z0-9+\\/]{2}==|[A-Za-z0-9+\\/]{3}=|[A-Za-z0-9+\\/]{4})$"
	Numeric string = "^[0-9]+$"
	Mobile  string = "^((13[0-9])|(14[5,7,9])|(15[0-3,5-9])|(17[0,1,2,3,5-8])|(18[0-9])|166|198|199|(147))\\d{8}$"
	RGBColor  string = "^rgb\\(\\s*(0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])\\s*,\\s*(0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])\\s*,\\s*(0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])\\s*\\)$"
	FullURL      = `^(?:ftp|tcp|udp|wss?|https?):\/\/[\w\.\/#=?&]+$`
	URLSchema    = `((ftp|tcp|udp|wss?|https?):\/\/)`
)

var (
	userRegexp    = regexp.MustCompile("^[a-zA-Z0-9!#$%&'*+/=?^_`{|}~.-]+$")
	hostRegexp    = regexp.MustCompile("^[^\\s]+\\.[^\\s]+$")
	userDotRegexp = regexp.MustCompile("(^[.]{1})|([.]{1}$)|([.]{2,})")
	rxEmail       = regexp.MustCompile(Email)
	rxBase64      = regexp.MustCompile(Base64)
	rxFloat       = regexp.MustCompile(Float)
	rxInt         = regexp.MustCompile(Int)
	rxNumeric     = regexp.MustCompile(Numeric)
	rxMobile      = regexp.MustCompile(Mobile)
	rxASCII     = regexp.MustCompile("^[\x00-\x7F]+$")
	rxRGBColor  = regexp.MustCompile(RGBColor)
	rxHexColor  = regexp.MustCompile("^#?([0-9a-fA-F]{3}|[0-9a-fA-F]{6})$")
	rxFullURL   = regexp.MustCompile(FullURL)
	rxURLSchema = regexp.MustCompile(URLSchema)
)

// IsEmail
//  @Description  验证email
//  @Param str
//  @Return bool
func IsEmail(str string) bool {
	_, err := regexp.MatchString(`^([\w\.\_\-]{2,10})@(\w{1,}).([a-z]{2,4})$`, str)
	return err != nil
}

// HasURLSchema string.
func HasURLSchema(s string) bool {
	return s != "" && rxURLSchema.MatchString(s)
}

// IsFullURL
//  @Description  验证是否为完整url
//  @Param s
//  @Return bool
func IsFullURL(s string) bool {
	return s != "" && rxFullURL.MatchString(s)
}

// IsURL
//  @Description   验证是否为url
//  @Param s
//  @Return bool
func IsURL(s string) bool {
	if s == "" {
		return false
	}

	_, err := url.Parse(s)
	return err == nil
}

// IsBase64
//  @Description  验证是否为base64
//  @Param str
//  @Return bool
func IsBase64(str string) bool {
	return rxBase64.MatchString(str)
}

// IsFloat
//  @Description  验证是否为float
//  @Param str
//  @Return bool
func IsFloat(str string) bool {
	return str != "" && rxFloat.MatchString(str)
}

// IsIP
//  @Description  验证是否为IPV4、IPV6
//  @Param str
//  @Return bool
func IsIP(str string) bool {
	return net.ParseIP(str) != nil
}

// IsIPv4
//  @Description  验证是否为IPV4
//  @Param str
//  @Return bool
func IsIPv4(str string) bool {
	ip := net.ParseIP(str)
	return ip != nil && strings.Contains(str, ".")
}

// IsIPv6
//  @Description  验证是否为IPV6
//  @Param str
//  @Return bool
func IsIPv6(str string) bool {
	ip := net.ParseIP(str)
	return ip != nil && strings.Contains(str, ":")
}

// IsInt
//  @Description  验证是否为整形
//  @Param str
//  @Return bool
func IsInt(str string) bool {
	if IsNull(str) {
		return true
	}
	return rxInt.MatchString(str)
}

// IsJSON
//  @Description  验证是否为json
//  @Param str
//  @Return bool
func IsJSON(str string) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(str), &js) == nil
}

// IsNull
//  @Description  验证是否为null
//  @Param str
//  @Return bool
func IsNull(str string) bool {
	return len(str) == 0
}

// IsNumeric
//  @Description  验证是否为数字
//  @Param str
//  @Return bool
func IsNumeric(str string) bool {
	if IsNull(str) {
		return true
	}
	return rxNumeric.MatchString(str)
}

// IsTime
//  @Description  验证是否为时间
//  @Param str
//  @Param format
//  @Return bool
func IsTime(str string, format string) bool {
	_, err := time.Parse(format, str)
	return err == nil
}

// IsUnixTime
//  @Description  验证是否为UnixTime
//  @Param str
//  @Return bool
func IsUnixTime(str string) bool {
	if _, err := strconv.Atoi(str); err == nil {
		return true
	}
	return false
}

// IsMAC checks if a string is valid MAC address.
// Possible MAC formats:
// 01:23:45:67:89:ab
// 01:23:45:67:89:ab:cd:ef
// 01-23-45-67-89-ab
// 01-23-45-67-89-ab-cd-ef
// 0123.4567.89ab
// 0123.4567.89ab.cdef
func IsMAC(str string) bool {
	_, err := net.ParseMAC(str)
	return err == nil
}

// IsMobileNumber
//  @Description  验证是否为手机号
//  @Param str
//  @Return bool
func IsMobileNumber(str string) bool {
	return str != "" && rxMobile.MatchString(str)
}

// IsHexColor
//  @Description  验证是否为十六进制颜色
//  @Param s
//  @Return bool
func IsHexColor(s string) bool {
	return s != "" && rxHexColor.MatchString(s)
}

// IsRGBColor
//  @Description  验证是否为为RGB颜色
//  @Param s
//  @Return bool
func IsRGBColor(str string) bool {
	return str != "" && rxRGBColor.MatchString(str)
}

// IsASCII
//  @Description  验证是否为ASCII
//  @Param s
//  @Return bool
func IsASCII(str string) bool {
	return str != "" && rxASCII.MatchString(str)
}

// StartsWith
//  @Description  是否以字符sub开头
//  @Param str
//  @Param sub
//  @Return bool
func StartsWith(str, sub string) bool {
	if str == "" {
		return false
	}

	return strings.HasPrefix(str, sub)
}

// EndsWith
//  @Description  是否以字符sub结尾
//  @Param str
//  @Param sub
//  @Return bool
func EndsWith(str, sub string) bool {
	if str == "" {
		return false
	}

	return strings.HasSuffix(str, sub)
}

// StringContains
//  @Description  是否包含字符sub
//  @Param str
//  @Param sub
//  @Return bool
func StringContains(str, sub string) bool {
	if str == "" {
		return false
	}
	return strings.Contains(str, sub)
}
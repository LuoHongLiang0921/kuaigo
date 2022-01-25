// @Description

package knet

import (
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/kcast"
	"net/url"
	"time"

)

// URL wrap url.URL.
type URL struct {
	Scheme     string
	Opaque     string        // encoded opaque data
	User       *url.Userinfo // username and password information
	Host       string        // host or host:port
	Path       string        // path (relative paths may omit leading slash)
	RawPath    string        // encoded path hint (see EscapedPath method)
	ForceQuery bool          // append a query ('?') even if RawQuery is empty
	RawQuery   string        // encoded query values, without '?'
	Fragment   string        // fragment for references, without '#'
	HostName   string
	Port       string
	params     url.Values
}

// ParseURL
//  @Description  解析URL字符串
//  @Param raw
//  @Return *URL
//  @Return error
func ParseURL(raw string) (*URL, error) {
	u, e := url.Parse(raw)
	if e != nil {
		return nil, e
	}

	return &URL{
		Scheme:     u.Scheme,
		Opaque:     u.Opaque,
		User:       u.User,
		Host:       u.Host,
		Path:       u.Path,
		RawPath:    u.RawPath,
		ForceQuery: u.ForceQuery,
		RawQuery:   u.RawQuery,
		Fragment:   u.Fragment,
		HostName:   u.Hostname(),
		Port:       u.Port(),
		params:     u.Query(),
	}, nil
}

// Password
//  @Description  从URL中获取密码
//  @Receiver u
//  @Return string
//  @Return bool
func (u *URL) Password() (string, bool) {
	if u.User != nil {
		return u.User.Password()
	}
	return "", false
}

// Username
//  @Description  从URL中获取用户名
//  @Receiver u
//  @Return string
func (u *URL) Username() string {
	return u.User.Username()
}

// QueryInt
//  @Description  从url中获取参数值 并转为INT
//  @Receiver u
//  @Param field
//  @Param expect
//  @Return ret
func (u *URL) QueryInt(field string, expect int) (ret int) {
	ret, err := kcast.ToIntE(u.Query().Get(field))
	if err != nil {
		return expect
	}

	return ret
}

// QueryInt64
//  @Description  从url中获取参数值 并转为INT64
//  @Receiver u
//  @Param field
//  @Param expect
//  @Return ret
func (u *URL) QueryInt64(field string, expect int64) (ret int64) {
	ret, err := kcast.ToInt64E(u.Query().Get(field))
	if err != nil {
		return expect
	}

	return ret
}

// QueryString
//  @Description  从url中获取参数值 并转为STRING
//  @Receiver u
//  @Param field
//  @Param expect
//  @Return ret
func (u *URL) QueryString(field string, expect string) (ret string) {
	ret = expect
	if mi := u.Query().Get(field); mi != "" {
		ret = mi
	}

	return
}

// QueryDuration
//  @Description  returns provided field's value in duration type.
//  @Receiver u
//  @Param field
//  @Param expect
//  @Return ret
func (u *URL) QueryDuration(field string, expect time.Duration) (ret time.Duration) {
	ret, err := kcast.ToDurationE(u.Query().Get(field))
	if err != nil {
		return expect
	}

	return ret
}

// QueryBool QueryString
//  @Description  从url中获取参数值 并转为BOOL
//  @Receiver u
//  @Param field
//  @Param expect
//  @Return ret
func (u *URL) QueryBool(field string, expect bool) (ret bool) {
	ret, err := kcast.ToBoolE(u.Query().Get(field))
	if err != nil {
		return expect
	}
	return ret
}

// Query
//  @Description  解析url并返回value
//  @Receiver u
//  @Return url.Values
func (u *URL) Query() url.Values {
	v, _ := url.ParseQuery(u.RawQuery)
	return v
}

package kentity

// ReqPassportOpenLogin 请求passport三方登录接口参数实体
type ReqPassportOpenLogin struct {
	OpenId     string `json:"openId"`
	Type       int    `json:"type"`
	AppId      int    `json:"appId"`
	DeviceType int    `json:"deviceType,omitempty"` // 登录设备类型：默认0  1：油果APP  2油果WEB端  3油果小程序
	ExpireTime int    `json:"expireTime"`
	ReqSrc     string `json:"reqSrc"`
	RegIp      string `json:"regIp"`
	Platform   bool   `json:"platform,omitempty"`
}

type RspPassportLoginData struct {
	PassportInfo PassportInfo `json:"passportInfo"`
	Mapping      interface{}  `json:"mapping"`
	Register     bool         `json:"register"` //标示是否为注册操作，只有注册才会出现，返回值为true
}

// PassportInfo 凭证详情
type PassportInfo struct {
	Birthday   string `json:"birthday"`
	Summary    string `json:"summary"`
	Area       string `json:"area"`
	AddTime    int64  `json:"addTime"`
	NickName   string `json:"nickName"`
	Sex        int    `json:"sex"`
	Icon       string `json:"icon"`
	RegIP      string `json:"regIp"`
	Token      string `json:"token"`
	UID        int64  `json:"uid"`
	ModifyTime int64  `json:"modifyTime"`
	Phone      string `json:"phone"`
	Background string `json:"background"`
	AreaType   int    `json:"areaType"`
	JwtToken   string `json:"jwtToken"`
	Account    string `json:"account"`
	Status     int    `json:"status"`
}

// RspPassportLogin 用户凭证服务返回数据实体
type RspPassportLogin struct {
	Code int                   `json:"code"`
	Msg  string                `json:"msg"`
	Data *RspPassportLoginData `json:"data"`
}

package kentity

import (
	"encoding/json"
	"github.com/LuoHongLiang0921/kuaigo/kpkg/constant"
	"github.com/LuoHongLiang0921/kuaigo/kpkg/core/server/kgin"
	"github.com/LuoHongLiang0921/kuaigo/kpkg/ecode"
	"github.com/LuoHongLiang0921/kuaigo/kpkg/errs"
	"github.com/LuoHongLiang0921/kuaigo/kpkg/klog"
	"strconv"
	"time"
)

// TabbyHeader
type TabbyHeader struct {
	//UA SNSNDK-HTTP/1.0.0 (iPhone; iOS 12.2; Scale/2.00)
	UA string `header:"User-Agent" json:"ua" binding:"required"`
	//RV 请求接口的版本号(request version)
	RV string `header:"rv" json:"rv" binding:"required"`
	//RT	客户端发起请求时客户端时间(request time),单位毫秒
	RT int64 `header:"rt" json:"rt" binding:"required"`
	//PK 请求app的包名标识(package)
	PK string `header:"pk" json:"pk" binding:"required"`
	//SI 请求会话唯KEY密匙体
	SI string `header:"si" json:"si" binding:"required"`
	//OV 兼容Android5.0以下需要base64两次的问题
	OV int `header:"ov" json:"ov"`
	//CS 可以不传，以供不同推广渠道用
	CS string `header:"cs" json:"cs"`
	//P 加密(json_encode(公共数据体数组)),每次请求必须传，客户端公共信息
	P string `header:"p" json:"p" binding:"required"`
	//S 加密(json_encode(敏感数据组)),手机号、密码、身份证号码等信息需加密传输
	S string `header:"s" json:"s"`

	ST    int64  `header:"st" json:"st"`    //服务端端接收到请求时服务端时间
	AppId int    `header:"appId" json:"ai"` //appId生成
	IsApp bool   `json:"isApp"`             //来源是否是app
	TI    string `header:"ti" json:"ti"`    //请求跟踪标识
	SK    string `header:"sk" json:"sk"`    //请求签名校验key值

	// PParams ..
	PParams PParams `json:"-"`
}

//PParams Header 中 的 p参数
type PParams struct {
	//DID 设备唯一标识
	DID string `header:"did" json:"did" binding:"required"`
	//DName 用户设置的设备名称,android部分不能设置
	DName string `header:"dName" json:"dname" binding:"required"`
	//SL 客户端系统当前语言(system language)
	SL string `header:"sl" json:"sl" binding:"required"`
	//SV 客户端系统版本(system version)
	SV string `header:"sv" json:"sv" binding:"required"`
	//OS 设备系统类型(ios/android)
	OS string `header:"os" json:"os" binding:"required"`
	//AV APP当前版本号(application version)
	AV string `header:"av" json:"av"`
	//BV APPbuild版本号，按Build时间方式生成(build version)
	BV string `header:"bv" json:"bv" binding:"required"`
	//DM 客户端设备型号(device model)
	DM string `header:"dm" json:"dm" binding:"required"`
	//DMID 客户端网卡MAC地址(device mac id)
	DMID string `header:"dmid" json:"dmid"`
	//DAID 客户端广告跟踪ID(device advertisement id)
	DAID string `header:"daid" json:"daid"`
	//Network 0没有网络,1其他,2 2G,3 3G,4 4G,5 wifi,6 蜂窝网络
	Network int `header:"network" json:"network" binding:"required"`
	//DW 客户端屏幕宽(device width)
	DW int `header:"dw" json:"dw" binding:"required"`
	//DH 客户端屏幕高(device height)
	DH int `header:"dh" json:"dh" binding:"required"`
	//Lon 客户端发请求时所在位置经度
	Lon float64 `header:"lon" json:"lon"`
	//Lag 客户端发请求时所在位置纬度
	Lag float64 `header:"lag" json:"lag"`
	//Channel 客户端安装渠道(install channel)
	Channel string `header:"channel" json:"channel" binding:"required"`
	//IMEI android的imei号
	IMEI string `header:"imei" json:"imei"`
	//SdkID 项目分配的固定值
	SdkID string `header:"sdkId" json:"sdkid" binding:"required"`
	//AK 用户系统成功后授权的访问token(AccessToken)
	AK string `header:"ak" json:"ak"`

	FPP              string `header:"fpp" json:"fpp"`                           //请求来源的父级页面
	FP               string `header:"fp" json:"fp"`                             //请求来源的页面
	UserID           string `header:"userId" json:"userId"`                     //当前登陆用户Id
	UserTokenIsValid bool   `header:"userTokenIsValid" json:"userTokenIsValid"` // 当前用户的Token是否有效

	// 2021-02-03 新增 yanghanwei
	// AndId Android设备的Android ID
	AndId string `header:"andId" json:"andId" binding:"required"`
	// OaId 设备标识符,Android系统10版本以上必传
	OaId string `header:"oaId" json:"oaId" binding:"required"`
	// Hm 手机品牌
	Hm string `header:"hm" json:"hm" binding:"required"`
	// Ht 手机型号
	Ht string `header:"ht" json:"ht" binding:"required"`
	// 手机sim卡标识，0 未知 1 中国移动 2 中国联通 3 中国电信
	Carrier int `header:"carrier" json:"carrier" binding:"required"`
	// Cpu
	Cpu string `header:"cpu" json:"cpu" binding:"required"`
	// CpuId
	CpuId string `header:"cpuid" json:"cpuid" binding:"required"`
	// Dpi 手机屏幕dpi
	Dpi int `header:"dpi" json:"dpi" binding:"required"`
	// abTest 参数
	AB string `header:"ab" json:"ab"`
}

// AccessLogP 入口日志p参数
type AccessLogP struct {
	//AV APP当前版本号(application version)
	AV string `json:"av"`
	//DID 设备唯一标识
	DID string `json:"did"`
	//OS 设备系统类型(ios/android)
	OS string `json:"os"`
	//SV 客户端系统版本(system version)
	SV string `json:"sv"`
	//Channel 客户端安装渠道
	Channel string `json:"channel"`
}

// GetAccessLog 获取 入口
// @Description:
// @receiver p p参数
func (p *PParams) GetAccessLog() *AccessLogP {
	var newP AccessLogP
	newP.AV = p.AV
	newP.OS = p.OS
	newP.DID = p.DID
	newP.SV = p.SV
	newP.Channel = p.Channel
	return &newP
}

// GetServiceSource 获取来源
// @Description: 根据p 参数os 设置不同的
// @receiver p参数
// @return string
func (p *PParams) GetServiceSource() string {
	switch p.OS {
	case constant.OsAndroid:
		fallthrough
	case constant.OsIos:
		return klog.ServiceSourceApp
	case constant.OsH5:
		return klog.ServiceSourceH5
	case constant.OsWeb:
		return klog.ServiceSourceWeb
	default:
		return p.OS
	}
}

//后台接口header
type AdminHeader struct {
	//后台用户id
	AdminId string `header:"aid" json:"aid" binding:"required"`
	//后台用户登录名
	AdminLoginName string `header:"aln" json:"aln" binding:"required"`
	//后台用户真实名字
	AdminRealName string `header:"arn" json:"arn" binding:"required"`
	//cas 后台应用唯一code标示
	AdminAppCode string `header:"aac" json:"aac" binding:"required"`
	//appId
	AppId int `header:"appId" json:"ai"  binding:"required"`
	//请求跟踪标识
	TI string `header:"ti" json:"ti"`
}

// BindHeader
func BindHeader(ctx *kgin.TContext) *TabbyHeader {
	//todo: c.ShouldBindHeader
	pk := ctx.GetHeader(constant.HeaderFieldPk)
	return &TabbyHeader{
		RV:    ctx.GetHeader(constant.HeaderFieldRv),
		RT:    getRt(ctx),
		ST:    getSt(ctx),
		PK:    pk,
		AppId: getAppId(ctx),
		SI:    ctx.GetHeader(constant.HeaderFieldSi),
		IsApp: false,
		UA:    ctx.GetHeader(constant.HeaderFieldUserAgent),
		TI:    ctx.GetHeader(constant.HeaderFieldTi),
		SK:    ctx.GetHeader(constant.HeaderFieldSk),
		OV:    getOv(ctx),
		CS:    ctx.GetHeader(constant.HeaderFieldCs),
	}
}

// GetPParam 获取PParam
func GetPParam(c *kgin.TContext) (*PParams, error) {
	p := c.GetHeader(constant.KeyP)
	if p == "" {
		return nil, errs.NewCustomError(ecode.CodeHeaderPParamError, "p参数不合法")
	}
	var pp PParams
	if err := json.Unmarshal([]byte(p), &pp); err != nil {
		return nil, errs.NewCustomError(ecode.CodeHeaderPParamError, "解析p参数到PParams失败")
	}
	return &pp, nil
}

// 绑定后台接口header
func BindAdminHeader(ctx *kgin.TContext) *AdminHeader {
	return &AdminHeader{
		AdminId:        ctx.GetHeader(constant.AdminHeaderFieldAid),
		AdminLoginName: ctx.GetHeader(constant.AdminHeaderFieldAln),
		AdminAppCode:   ctx.GetHeader(constant.AdminHeaderFieldAac),
		AppId:          getAppId(ctx),
		TI:             ctx.GetHeader(constant.HeaderFieldTi),
	}
}

// getAppId
func getAppId(ctx *kgin.TContext) int {
	appIdStr := ctx.GetHeader(constant.HeaderFieldAi)
	appId, err := strconv.Atoi(appIdStr)
	if err != nil {
		return 0
	}
	return appId
}

// getRt
func getRt(ctx *kgin.TContext) int64 {
	rtStr := ctx.GetHeader(constant.HeaderFieldRt)
	rt, err := strconv.ParseInt(rtStr, 10, 64)
	if err == nil {
		return rt
	}
	return 0
}

// getSt
func getSt(ctx *kgin.TContext) int64 {
	stStr := ctx.GetHeader(constant.HeaderFieldSt)
	stTmp, err := strconv.ParseInt(stStr, 10, 64)
	if err == nil {
		if stTmp > 0 {
			return stTmp
		}
		return time.Now().UnixNano() / 1e6
	}
	return 0
}

// getOv
func getOv(ctx *kgin.TContext) int {
	ovStr := ctx.GetHeader(constant.HeaderFieldOv)
	ov, err := strconv.Atoi(ovStr)
	if err == nil {
		return ov
	}
	return 0
}

// SetTi
func SetTi(i int, header *TabbyHeader) {
	if i > 0 && header.TI != "" {
		// String.format("%s-%d", this.ti, i);
		header.TI = header.TI + "-" + strconv.Itoa(i)
	}
}

//// SetChannel 设置渠道
//func SetChannel(channel string, public *PParams) {
//	if channel != "" {
//		channel = strings.Trim(channel, "")
//		channel = strings.Replace(channel, "\n", "", -1)
//		channel = strings.Replace(channel, "\t", "", -1)
//	}
//	public.Channel = channel
//}

// SetDname 设置浏览器名称
func SetDname(dname string, public *PParams) {
	// todo
}

// SetAppId 设置AppId
func SetAppId(ctx *kgin.TContext, appId string) {
	ctx.Set(constant.KeyAppId, appId)
}

// SetHeader
func SetHeader(ctx *kgin.TContext, header *TabbyHeader) {
	ctx.Set(constant.KeyHeader, header)
}

// SetAdminHeader
func SetAdminHeader(ctx *kgin.TContext, adminHeader *AdminHeader) {
	ctx.Set(constant.KeyAdminHeader, adminHeader)
}

// SetPParam
func SetPParam(ctx *kgin.TContext, public *PParams) {
	ctx.Set(constant.KeyPParam, public)
}

// SetSecret
func SetSecret(ctx *kgin.TContext, sMap map[string]string) {
	ctx.Set(constant.KeySecret, sMap)
}

//// SetCurrUser
//func SetCurrUser(ctx *kgin.TContext, claims *xjwt.MyClaims) {
//	ctx.Set(constant.KeyCurrUser, claims)
//}

// Clear
func Clear(ctx *kgin.TContext) {
	ctx.Set(constant.KeySecret, nil)
	ctx.Set(constant.KeyPParam, nil)
	ctx.Set(constant.KeyHeader, nil)
	ctx.Set(constant.KeyCurrUser, nil)
}

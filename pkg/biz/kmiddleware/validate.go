package kmiddleware

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"github.com/LuoHongLiang0921/kuaigo/pkg/biz/kentity"
	"github.com/LuoHongLiang0921/kuaigo/pkg/biz/tcontext"
	"github.com/LuoHongLiang0921/kuaigo/pkg/constant"
	"github.com/LuoHongLiang0921/kuaigo/pkg/core/server/ginserver"
	"github.com/LuoHongLiang0921/kuaigo/pkg/ecode"
	"github.com/LuoHongLiang0921/kuaigo/pkg/errs"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/klog"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// RequiredAppValidLegal APP header公参统一校验
//  @Description:
//  @Receiver m
//  @Param c
func (m *Middleware) RequiredAppValidLegal() []ginserver.HandlerFunc {
	return []ginserver.HandlerFunc{
		m.ExtractWebAppHeader,
		m.AppHeaderValidate,
		m.AppSignValidate,
	}
}

// RequiredWebValidLegal web端 header公参统一校验
//  @Description:
//  @Receiver m
//  @Param c
func (m *Middleware) RequiredWebValidLegal() []ginserver.HandlerFunc {
	return []ginserver.HandlerFunc{
		m.ExtractWebAppHeader,
		m.AppHeaderValidate,
		m.AppSignValidate,
	}
}

// RequiredInnerValidLegal 内部服务 header公参统一校验
//  @Description:
//  @Receiver m
func (m *Middleware) RequiredInnerValidLegal() []ginserver.HandlerFunc {
	return []ginserver.HandlerFunc{
		m.ExtractInnerServiceHeader,
	}
}

// RequiredAdminValidLegal 内部服务 header公参统一校验
//  @Description:
//  @Receiver m
func (m *Middleware) RequiredAdminValidLegal() []ginserver.HandlerFunc {
	return []ginserver.HandlerFunc{
		m.ExtractAdminHeader,
	}
}

// RequiredOpenValidLegal 管理系统 header公参统一校验
//  @Description:
//  @Receiver m
func (m *Middleware) RequiredOpenValidLegal() []ginserver.HandlerFunc {
	return []ginserver.HandlerFunc{
		m.ExtractOuterServiceHeader,
	}
}

// ExtractWebAppHeader 抽取header到xlog.Common 中 供api 层（app）使用
//  @Description:
//  @Receiver m
//  @Param c
func (m *Middleware) ExtractWebAppHeader(c *ginserver.TContext) {
	klog.Debug("抽取header到xlog.Common 中starting")
	beg := time.Now()
	tabbyHeader := kentity.BindHeader(c)
	kentity.SetHeader(c, tabbyHeader)
	var com klog.Common
	com.AppId = tabbyHeader.AppId
	com.RequestIp = c.ClientIP()
	com.RequestUri = c.Request.RequestURI
	com.TraceId = tabbyHeader.TI
	com.ProcessCode = klog.ProcessCodeRequest
	pp, err := kentity.GetPParam(c)
	if err != nil {
		c.Set(constant.CommonKey, com)
		klog.AccessLogger.WithContext(klog.WithCommonLog(c, com)).Info("before access log,p invalid")
		m.abortWithErrorJSON(c, err)
		return
	}
	kentity.SetPParam(c, pp)
	com.ServiceSource = pp.GetServiceSource()
	c.Set(constant.CommonKey, com)
	klog.AccessLogger.WithContext(klog.WithCommonLog(c, com)).Info("before access log")
	c.Next()
	costTime := time.Since(beg) / 1e6
	com.ProcessCode = klog.ProcessCodeResponse
	com.CostTime = int64(costTime)
	com.Code = getRespCodeFromContext(c)
	com.UID = tabbyHeader.PParams.UserID
	com.P = tabbyHeader.PParams.GetAccessLog()
	klog.AccessLogger.WithContext(klog.WithCommonLog(c, com)).Info("after access log")
}

// ExtractInnerServiceHeader
//  @Description   内部服务调用中间件  serviceSource为 service
//  @Receiver m 中间件
//  @Param c 请求上下文
func (m *Middleware) ExtractInnerServiceHeader(c *ginserver.TContext) {
	beg := time.Now()
	tabbyHeader := kentity.BindHeader(c)
	kentity.SetHeader(c, tabbyHeader)
	var com klog.Common
	com.AppId = tabbyHeader.AppId
	com.RequestIp = c.ClientIP()
	com.RequestUri = c.Request.RequestURI
	com.TraceId = tabbyHeader.TI
	com.ServiceSource = klog.ServiceSourceService
	c.Set(constant.CommonKey, com)
	com.ProcessCode = klog.ProcessCodeRequest
	//com.P = nil
	klog.AccessLogger.WithContext(klog.WithCommonLog(c, com)).Info("before access log")
	c.Next()
	costTime := time.Since(beg) / 1e6
	com.ProcessCode = klog.ProcessCodeResponse
	com.CostTime = int64(costTime)
	com.Code = getRespCodeFromContext(c)
	//com.P = nil
	klog.AccessLogger.WithContext(klog.WithCommonLog(c, com)).Info("after access log")
}

// ExtractOuterServiceHeader 公网第三方、或H5无P验签调用中间件
// @Description  serviceSource为 service
// @receiver m 中间件
// @Param c 请求上下文
func (m *Middleware) ExtractOuterServiceHeader(c *ginserver.TContext) {
	beg := time.Now()
	var com klog.Common
	com.AppId = c.GetInt("OuterAppId")
	com.RequestIp = c.ClientIP()
	com.RequestUri = c.Request.RequestURI
	com.ServiceSource = klog.ServiceSourceService
	c.Set(constant.CommonKey, com)
	com.ProcessCode = klog.ProcessCodeRequest
	klog.AccessLogger.WithContext(klog.WithCommonLog(c, com)).Info("before access log")
	c.Next()
	costTime := time.Since(beg) / 1e6
	com.ProcessCode = klog.ProcessCodeResponse
	com.CostTime = int64(costTime)
	com.Code = getRespCodeFromContext(c)
	klog.AccessLogger.WithContext(klog.WithCommonLog(c, com)).Info("after access log")
}

// ExtractAdminHeader 后台中间件
//  @Description
//  @Receiver m
//  @Param c
func (m *Middleware) ExtractAdminHeader(c *ginserver.TContext) {
	beg := time.Now()
	adminHeader := kentity.BindAdminHeader(c)
	kentity.SetAdminHeader(c, adminHeader)
	var com klog.Common
	com.AppId = adminHeader.AppId
	com.RequestIp = c.ClientIP()
	com.RequestUri = c.Request.RequestURI
	com.TraceId = adminHeader.TI
	com.UID = adminHeader.AdminId
	com.ServiceSource = klog.ServiceSourceAdmin
	c.Set(constant.CommonKey, com)
	com.ProcessCode = klog.ProcessCodeRequest
	//com.P = ""
	klog.AccessLogger.WithContext(klog.WithCommonLog(c, com)).Info("before access log")
	c.Next()
	costTime := time.Since(beg) / 1e6
	com.ProcessCode = klog.ProcessCodeResponse
	com.CostTime = int64(costTime)
	com.Code = getRespCodeFromContext(c)
	//com.P = ""
	klog.AccessLogger.WithContext(klog.WithCommonLog(c, com)).Info("after access log")
}

//Deprecated: 请使用 组合中间件 RequiredAppValidLegal 替代
//  @Description   解密 p s 参数
//  @Receiver m
//  @Param c 请求上下文
func (m *Middleware) App(c *ginserver.TContext) {
	m.AppHeaderValidate(c)
}

// AppHeaderValidate
//  @Description: 解析APP header中必传参数p s
//  @Receiver m
//  @Param c
func (m *Middleware) AppHeaderValidate(c *ginserver.TContext) {
	// 解析"p"参数
	err := m.parsePParam(c)
	if err != nil {
		m.abortWithErrorJSON(c, err)
	}
	err = m.parseSParam(c)
	if err != nil {
		m.abortWithErrorJSON(c, err)
	}
	m.AppSignValidate(c)
}

// AppSignValidate
//  @Description  app 验签
//  @Receiver m
//  @Param c
func (m *Middleware) AppSignValidate(c *ginserver.TContext) {
	klog.Debug("app 验签 starting")
	noSign := c.GetHeader("noSign")

	if noSign == "1" {
		if m.AllowCloseSignValidate {
			klog.Warn("签名验签已关闭，线上环境注意打开")
			return
		}
	}
	// 从上下文获取信息
	snsNdkHeader := tcontext.GetHeader(c)
	snsNdkPublic := tcontext.GetPParam(c)
	if snsNdkHeader == nil || snsNdkPublic == nil {
		klog.Warn("公共参数TabbyHeader或PParams不合法，非法请求")
		m.abortWithErrorJSON(c, errs.NewCustomError(ecode.CodeInternalServerError, "公共参数TabbyHeader或PParams不合法，非法请求"))
		return
	}
	// header部分签名
	p := c.GetHeader(constant.KeyP)
	s := c.GetHeader(constant.KeyS)

	// 获取加密header
	headerMd5 := getMd5AppHeader(p, s, snsNdkPublic, snsNdkHeader)
	// 获取加密body
	formBodyMd5, err := getMd5AppBody(c, snsNdkPublic, snsNdkHeader)
	if err != nil {
		m.abortWithErrorJSON(c, err)
		return
	}
	// 封装最终加密的值
	var sb bytes.Buffer
	sb.WriteString(headerMd5)
	sb.WriteString(formBodyMd5)
	sb.WriteString(snsNdkHeader.SI)

	sk := snsNdkHeader.SK
	//signature := Md5.md5Encode(sb.toString())
	signature := fmt.Sprintf("%x", md5.Sum(sb.Bytes()))
	klog.Debug("最后加密前：" + sb.String() + ", 最后加密后:" + signature)

	// 开始验签
	if !strings.EqualFold(signature, snsNdkHeader.SK) {
		klog.Debug("签名验证失败,sk:" + sk + ",signature:" + signature + ",snsNdkHeader:" + signature + ",p:" + p + ",public{}：" + sb.String() + ", 最后加密后:" + signature)
		m.abortWithErrorJSON(c, errs.NewCustomError(ecode.CodeInValidParameterSignature, "签名验证失败"))
	}
}

// AdminSignValidate
//  @Description  后台验证
//  @Receiver m
//  @Param c
func (m *Middleware) AdminSignValidate(c *ginserver.TContext) {
	adminId := c.GetHeader(constant.HeaderFieldAdminId)
	adminAppId := c.GetHeader(constant.HeaderFieldAdminId)
	t := c.GetHeader(constant.HeaderFieldT)
	sk := c.GetHeader(constant.HeaderFieldSk)
	secret := m.adminJwt.JwtKey

	pattern := "\\d+" //反斜杠要转义
	flag, err := regexp.MatchString(pattern, t)
	if err != nil || !flag {
		klog.Debug("签名验证失败,时间参数错误,adminId:" + adminId + ",adminAppId:" + adminAppId + ",t:" + t + ",sk:" + sk)
		m.abortWithErrorJSON(c, errs.NewCustomError(ecode.CodeInvalidParameterTimestamp, "签名验证失败,时间参数错误"))
		return
	}
	//判断时间是否超过3分钟
	timeIfExpire(t, adminId, adminAppId, sk)
	// 加密header
	headerMd5 := getMd5AdminHeader(adminId, adminAppId, t)
	// 加密body
	formBodyMd5, err := getMd5AdminBody(c)
	if err != nil {
		m.abortWithErrorJSON(c, err)
		return
	}
	var sb bytes.Buffer
	sb.WriteString(headerMd5)
	sb.WriteString(formBodyMd5)
	sb.WriteString(secret)
	signature := fmt.Sprintf("%x", md5.Sum(sb.Bytes()))
	klog.Debug("最后加密前：" + sb.String() + ", 最后加密后:" + signature)

	// 开始验签
	if !strings.EqualFold(signature, sk) {
		klog.Debug("签名验证失败,adminId:" + adminId + ",adminAppId:" + adminAppId + ",secret:" + secret + ",t:" + t + ",signature:" + signature + ",sk:" + sk)
		m.abortWithErrorJSON(c, errs.NewCustomError(ecode.CodeInValidParameterSignature, "签名验证失败"))
		//todo 暂时移除异常抛出
	}
}

// timeIfExpire
//  @Description  判断时间是否超过3分钟
//  @Param t
//  @Param adminId
//  @Param adminAppId
//  @Param sk
func timeIfExpire(t string, adminId string, adminAppId string, sk string) {
	l, _ := strconv.ParseInt(t, 10, 64)
	currentTimeMillis := time.Now().UnixNano() / 1e6
	//时间超过3分钟
	if absInt(currentTimeMillis-l) > constant.SignExpireTime {
		// 打印日志,不抛出异常
		klog.Debug("签名验证失败,签名超过3分钟,adminId:" + adminId + ",adminAppId:" + adminAppId + ",t:" + t + ",sk:" + sk + ",current:" + strconv.FormatInt(currentTimeMillis, 10))
	}
}

// getMd5AppHeader
//  @Description  获取加密header
//  @Param p
//  @Param s
//  @Param snsNdkPublic
//  @Param snsNdkHeader
//  @Return string
func getMd5AppHeader(p string, s string, snsNdkPublic *kentity.PParams, snsNdkHeader *kentity.KuaigoHeader) string {
	var sbHeader bytes.Buffer
	//android 小于21的版本，考虑表情bug，客户端做了一次url encode,但kong网关已反解，所以验签时需要重编码
	if strings.EqualFold(snsNdkPublic.OS, constant.OsAndroid) && snsNdkHeader.OV < constant.OvMin && snsNdkHeader.OV > 0 {
		if p != "" {
			p = url.QueryEscape(p)
		}
		if s != "" {
			s = url.QueryEscape(s)
		}
	}
	sbHeader.WriteString(snsNdkHeader.RV)
	sbHeader.WriteString(strconv.FormatInt(snsNdkHeader.RT, 10))
	sbHeader.WriteString(snsNdkHeader.PK)
	sbHeader.WriteString(snsNdkHeader.TI)
	sbHeader.WriteString(p)
	if s != "" {
		sbHeader.WriteString(s)
	}
	headerMd5 := fmt.Sprintf("%x", md5.Sum(sbHeader.Bytes()))
	klog.Debug("header加密前：" + sbHeader.String() + ", header加密后:" + headerMd5)
	return headerMd5
}

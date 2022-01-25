package kmiddleware

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/LuoHongLiang0921/kuaigo/pkg/biz/kentity"
	"github.com/LuoHongLiang0921/kuaigo/pkg/constant"
	"github.com/LuoHongLiang0921/kuaigo/pkg/core/server/ginserver"
	"github.com/LuoHongLiang0921/kuaigo/pkg/ecode"
	"github.com/LuoHongLiang0921/kuaigo/pkg/errs"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/kjwt"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/klog"
	"net/url"
	"sort"
	"strings"
)

// isDebug 是否为debug模式
// 	@Description   是否为debug模式
// 	@receiver m Middleware
//	@Param c xgin.TContext
// 	@return bool
func (m *Middleware) isDebug(c *ginserver.TContext) bool {
	debug := c.GetHeader("debug")
	return debug == "1"
}

// mustParseDIDAndAK
//  @Description
//  @Receiver m
//  @Param c
//  @Param pp
//  @Return error
func (m *Middleware) mustParseDIDAndAK(c *ginserver.TContext, pp *kentity.PParams) error {
	if pp.DID != "" {
		if pp.AK != "" {
			if m.jwt.JwtKey == "" {
				return errs.NewCustomError(ecode.CodeAuthTokenTimeout, "Jwt中间件未开启")
			}
			myClaims, err := kjwt.NewClient(m.jwt.JwtKey, m.jwt.Expire).ParseToken(pp.AK)
			if err != nil {
				if strings.Contains(err.Error(), "expired") {
					return errs.NewCustomError(ecode.CodeAuthTokenTimeout, "token已过期")
				}
				return errs.NewCustomError(ecode.CodeInvalidAuthToken, "解密token失败")
			}
			// TODO：单设备登录需要使用redis
			if myClaims != nil {
				if myClaims.UserID != "" {
					pp.UserID = myClaims.UserID
					pp.UserTokenIsValid = true
				}
				kentity.SetCurrUser(c, myClaims)
			}
		}
		//xentity.SetChannel(pp.Channel, pp)
		channel := pp.Channel
		if channel != "" {
			channel = strings.Trim(channel, "")
			channel = strings.Replace(channel, "\n", "", -1)
			channel = strings.Replace(channel, "\t", "", -1)
		}
		pp.Channel = channel
	}
	return nil
}

// parseP
//  @Description  解析p参数 debug 模式，不走空网关，p参数
//  @Receiver m
//  @Param c
//  @Return error
func (m *Middleware) parsePParam(c *ginserver.TContext) error {
	//p := c.GetHeader(constant.KeyP)
	//xlog.Debug(context.TODO(), "解密完p:"+p)
	//if p == "" {
	//	return errs.NewCustomError(ecode.CodeHeaderPParamError, "p参数不合法")
	//}
	//// p 转换为 PParams
	//var pp xentity.PParams
	//if err := json.Unmarshal([]byte(p), &pp); err != nil {
	//	return errs.NewCustomError(ecode.CodeHeaderPParamError, "解析p参数到PParams失败")
	//}
	pp, err := kentity.GetPParam(c)
	if err != nil {
		return errs.NewCustomError(ecode.CodeHeaderPParamError, "解析p参数到PParams失败")
	}
	m.parseDIDAndAK(c, pp)
	if v, ok := c.Get(constant.KeyHeader); ok {
		if tabbyHeader, ok := v.(*kentity.KuaigoHeader); ok {
			tabbyHeader.PParams = *pp
			c.Set(constant.KeyHeader, tabbyHeader)
		}
	}
	kentity.SetPParam(c, pp)
	return nil
}

// getMd5AppBody
//  @Description  获取加密body
//  @Param c
//  @Param p
//  @Param s
//  @Param snsNdkPublic
//  @Param snsNdkHeader
//  @Return string
//  @Return error
func getMd5AppBody(c *ginserver.TContext, snsNdkPublic *kentity.PParams, snsNdkHeader *kentity.KuaigoHeader) (string, error) {
	var sbBody bytes.Buffer
	// 获取form表单并排序
	listKey, err := getFormList(c)
	if err != nil {
		return "", err
	}
	// 排序
	sort.Strings(listKey)
	for _, key := range listKey {
		value, flag := c.GetPostForm(key)
		if flag {
			trimValue := strings.Trim(value, "")
			if value != "" && len(trimValue) > 0 {
				sbBody.WriteString(trimValue)
			}
		}
	}
	formBody := sbBody.String()
	if strings.EqualFold(snsNdkPublic.OS, constant.OsAndroid) && snsNdkHeader.OV < constant.OvMin && snsNdkHeader.OV > 0 {
		formBody = url.QueryEscape(formBody)
	}
	//formBodyMd5 := md5.md5Encode(formBody);
	formBodyMd5 := fmt.Sprintf("%x", md5.Sum([]byte(formBody)))
	klog.Debug("body加密前：" + formBody + ", body加密后:" + formBodyMd5)
	return formBodyMd5, nil
}

// getMd5AdminHeader
//  @Description  获取加密header
//  @Param adminId
//  @Param adminAppId
//  @Param t
//  @Return string
func getMd5AdminHeader(adminId string, adminAppId string, t string) string {
	var sbHeader bytes.Buffer
	sbHeader.WriteString(adminId)
	sbHeader.WriteString(adminAppId)
	sbHeader.WriteString(t)
	header := sbHeader.String()
	return fmt.Sprintf("%x", md5.Sum([]byte(header)))
}

// getMd5AdminBody
//  @Description  获取加密body
//  @Param c
//  @Return string
//  @Return error
func getMd5AdminBody(c *ginserver.TContext) (string, error) {
	var sbBody bytes.Buffer
	// 获取form表单并排序
	listKey, err := getFormList(c)
	if err != nil {
		return "", err
	}
	// 排序
	sort.Strings(listKey)
	for _, key := range listKey {
		value, flag := c.GetPostForm(key)
		if flag {
			trimValue := strings.Trim(value, "")
			if value != "" && len(trimValue) > 0 {
				sbBody.WriteString(trimValue)
			}
		}
	}
	// 此处java没有使用body参数
	formBody := ""
	return fmt.Sprintf("%x", md5.Sum([]byte(formBody))), nil
}

// absInt
//  @Description  取绝对值
//  @Param n
//  @Return int64
func absInt(n int64) int64 {
	if n < 0 {
		return -n
	}
	return n
}

// parseDIDAndAK
//  @Description  解析did、ak
//  @Receiver m
//  @Param c
//  @Param pp
func (m *Middleware) parseDIDAndAK(c *ginserver.TContext, pp *kentity.PParams) {
	if pp.DID != "" {
		if pp.AK != "" && m.jwt != nil {
			myClaims, err := kjwt.NewClient(m.jwt.JwtKey, m.jwt.Expire).ParseToken(pp.AK)
			if err != nil {
				myClaims = new(kjwt.MyClaims)
			}
			myClaims.SetError(err)
			// TODO：单设备登录需要使用redis
			if myClaims.UserID != "" {
				pp.UserID = myClaims.UserID
				pp.UserTokenIsValid = true
			}
			kentity.SetCurrUser(c, myClaims)
		}
		//xentity.SetChannel(pp.Channel, pp)
		channel := pp.Channel
		if channel != "" {
			channel = strings.Trim(channel, "")
			channel = strings.Replace(channel, "\n", "", -1)
			channel = strings.Replace(channel, "\t", "", -1)
		}
		pp.Channel = channel
	}
}

// parseS
//  @Description  解析 s 参数
//  @Receiver m
//  @Param c
//  @Return error
func (m *Middleware) parseSParam(c *ginserver.TContext) error {
	// 解析 "s"参数(解析S加密参数到map中由会话线程中保存)
	s := c.GetHeader(constant.KeyS)
	klog.Debug("解密完s:" + s)
	if s == "" {
		return nil
	}
	var sMap map[string]string
	if err := json.Unmarshal([]byte(s), &sMap); err != nil {
		return errs.NewCustomError(ecode.CodeHeaderSParamError, "解析s参数失败")
	}
	if sMap != nil {
		kentity.SetSecret(c, sMap)
	}
	return nil
}

// getFormList
//  @Description  根据contentType调整获取form 参数方法
//  @Param c
//  @Return []string
//  @Return error
func getFormList(c *ginserver.TContext) ([]string, error) {
	form := constant.KeyFormData
	formUrlencoded := constant.KeyFormUrlencoded
	contentType := c.ContentType()
	listKey := make([]string, 0)
	if strings.Contains(contentType, form) {
		form, err := c.MultipartForm()
		if err != nil {
			return nil, errs.NewCustomError(500, "获取form参数失败")
		}
		if form != nil {
			if len(form.Value) > 0 {
				for key := range form.Value {
					listKey = append(listKey, key)
				}
			}
		}
		return listKey, nil
	} else if strings.Contains(contentType, formUrlencoded) {
		err := c.Request.ParseForm()
		if err != nil {
			return nil, errs.NewCustomError(500, "获取form参数失败")
		}
		formMap := c.Request.PostForm
		if len(formMap) > 0 {
			for key := range formMap {
				listKey = append(listKey, key)
			}
		}
	}
	return listKey, nil
}

// getRespCodeFromContext
//  @Description
//  @Param c
//  @Return int
func getRespCodeFromContext(c *ginserver.TContext) int {
	if code, is := c.Get(constant.RespCode); is {
		return code.(int)
	}
	return 0
}

// abortWithErrorJSON
//  @Description
//  @Receiver m
//  @Param c
//  @Param err
func (m *Middleware) abortWithErrorJSON(c *ginserver.TContext, err error) {
	if _, ok := err.(errs.Error); ok {
		kentity.AbortWithErrorJSON(c, err)
		return
	}
	kentity.AbortWithJSONResponse(c, &kentity.Response{
		Code: ecode.CodeInternalServerError,
		Msg:  err.Error(),
	})
}

// getDeviceType 获取 deviceType
// 	@Description: 获取 deviceType
//	@Param appId
// 	@return int
func getDeviceType(appId int) int {
	if appId == 101 {
		return 7
	}
	if appId == 104 {
		return 4
	}
	if appId == 110 {
		return 1
	}
	return 0
}

// getPlatForm 获取 platForm
// 	@Description 获取 platForm
//	@Param appId
// 	@return bool
func getPlatForm(appId int) bool {
	// 分平台 波波，秒拍为false,油果其他为true
	followPlatformMp := map[int]bool{
		101: false,
		104: false,
		110: true,
	}
	return followPlatformMp[appId]
}

// getClientPlatform 获取platform
// 	@Description  获取platform
//	@Param plat
// 	@return string
func getClientPlatform(plat string) string {
	switch plat {
	case "ANDROID":
		return "1"
	case "IOS":
		return "2"
	case "BB-CREATOR":
		return "1"
	default:
		return ""
	}
}

// 0没有网络,1其他,2 2G,3 3G,4 4G,5 wifi,6 蜂窝网络
//
//网络类型：
//-1: 无网络 (或无法识别)
//0: 无网络 (或无法识别)
//1: WIFI
//2: 【2G】GPRS / win8 2G
//3: 【2.75G】EDGE
//4: 【3G】UMTS / IOS 3G (IOS客户端仅能识别是否3G) / win8 3G
//5: 【3.5G】HSDPA
//6: 【3.75G】HSUPA
//7: 【3.5G】HSPA
//8: 【2G】CDMA
//9: 【3G】EVDO_0(电信)
//10: 【3G】EVDO_A(电信)
//11: 【2.5G】1xRTT(电信2.5G)
//12: 【3G】HSPAP
//13: Ethernet (有线网)
//14: 【4G】LTE
//15: 【3G】EHRPD
func getNid(network int) string {
	switch network {
	case 0:
		return "0"
	case 1:
		return ""
	case 2:
		return "2"
	case 3:
		return "4"
	case 4:
		return "14"
	case 5:
		return "1"
	case 6:
		return "15"
	default:
		return ""
	}
}

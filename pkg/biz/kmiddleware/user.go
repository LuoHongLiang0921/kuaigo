package kmiddleware

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/LuoHongLiang0921/kuaigo/pkg/biz/kentity"
	"github.com/LuoHongLiang0921/kuaigo/pkg/biz/tcontext"
	"github.com/LuoHongLiang0921/kuaigo/pkg/constant"
	"github.com/LuoHongLiang0921/kuaigo/pkg/core/net/khttp"
	"github.com/LuoHongLiang0921/kuaigo/pkg/core/server/ginserver"
	"github.com/LuoHongLiang0921/kuaigo/pkg/core/storage/redis"
	"github.com/LuoHongLiang0921/kuaigo/pkg/ecode"
	"github.com/LuoHongLiang0921/kuaigo/pkg/errs"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/kjwt"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/klog"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// 登录失败错误码
var LoginFailed = 102002

type GuestUserModel interface {
	// UserOpenLogin 三方登录
	UserGuestLogin(c *ginserver.TContext, ctx context.Context, req kentity.ReqPassportOpenLogin) (*kjwt.MyClaims, error)
}

type guestUserModel struct {
	cache *redis.Redis
	addr  string
}

// NewGuestUserModel init实例
// 	@Description:  init实例
//	@Param cache redis.Redis
//	@Param addr
// 	@return XGuestUserModel
func NewGuestUserModel(cache *redis.Redis, addr string) GuestUserModel {
	return &guestUserModel{cache: cache, addr: addr}
}

// CheckRealLogin 用户登录校验
// 	@Description   用户登录校验
// 	@receiver m Middleware
//	@Param c ginserver.TContext
func (m *Middleware) CheckRealLogin(c *ginserver.TContext) {
	myClaims := tcontext.GetMyClaims(c)
	if myClaims == nil || myClaims.LoginRole == 2 {
		m.abortWithErrorJSON(c, errs.NewCustomError(ecode.CodeUserNotLogin, "用户未登录"))
		return
	}
}

// UserGuestLogin 游客登录
// 	@Description  游客登录
// 	@receiver um xguestUserModel
//	@Param c ginserver.TContext
//	@Param ctx context.Context
//	@Param req xentity.ReqPassportOpenLogin
// 	@return *xjwt.MyClaims
// 	@return error
func (um *guestUserModel) UserGuestLogin(c *ginserver.TContext, ctx context.Context, req kentity.ReqPassportOpenLogin) (*kjwt.MyClaims, error) {
	guestUserKey := fmt.Sprintf("%s%v:%s", constant.GuestUserKeyPre, req.AppId, req.OpenId)
	claimsStr := um.cache.Get(guestUserKey)
	// 1.判断 token 中解析出来的 用户信息，如果有可以直接返回，没有的话校验redis中是否存在
	myClaims := tcontext.GetMyClaims(c)
	if myClaims != nil {
		return myClaims, nil
	}
	// 2.判断缓存中是否存在
	claims := &kjwt.MyClaims{}
	if claimsStr != "" {
		err := json.Unmarshal([]byte(claimsStr), claims)
		if err == nil {
			return claims, nil
		}
	}
	// 3. 调用游客登录接口生成
	//reqHeader := http.Header{}
	//reqHeader.Set("Content-Type", "application/json")
	URL := um.addr + "/v1/login.json"
	pRsp := &kentity.RspPassportLogin{}
	resp, err := khttp.PostJson(ctx, URL, req)
	logFmt := fmt.Sprintf("url:%s, req:%+v, recResp:%+v", URL, req, pRsp)
	if err != nil {
		logFmt = logFmt + fmt.Sprintf(", err:%s", err.Error())
		msg := pRsp.Msg
		if pRsp.Code == 102 {
			msg = "凭证过期"
		}
		return claims, errs.NewCustomError(LoginFailed, msg)
	}
	resp.Json(pRsp)
	if pRsp.Code != 0 {
		klog.WithContext(ctx).Error(logFmt)
		return nil, errs.NewCustomError(LoginFailed, pRsp.Msg)
	}
	if pRsp.Data.PassportInfo.Status == 2 {
		return nil, errs.NewCustomError(LoginFailed, "用户已被封禁")
	}
	claims.UserID = fmt.Sprintf("%v", pRsp.Data.PassportInfo.UID)
	claims.AppID = req.AppId
	claims.DID = req.OpenId
	claims.PlatToken = pRsp.Data.PassportInfo.Token
	// 0,1： 真实登录， 2：游客登录（0 为了兼容老版本用户）
	claims.LoginRole = 2
	claims.StandardClaims = jwt.StandardClaims{}
	// 登录成功之后从新设置 缓存信息
	claimsByte, _ := json.Marshal(claims)
	flag := um.cache.Set(guestUserKey, claimsByte, time.Duration(constant.GuestLoginExpireTime)*time.Minute)
	if !flag {
		klog.WithContext(ctx).Error("游客登录设置redis缓存失败")
	}
	return claims, nil
}

// CheckGuestLogin 校验游客登录
// 	@Description  校验游客登录
// 	@receiver m Middleware
//	@Param c ginserver.TContext
func (m *Middleware) CheckGuestLogin(c *ginserver.TContext) {
	ctx := tcontext.WithContext(c)
	appId := tcontext.GetAppID(ctx)
	tParam := tcontext.GetPParam(c)
	myClaims := tcontext.GetMyClaims(c)
	if myClaims != nil {
		return
	}
	// http://wiki.yixiahd.com/pages/viewpage.action?pageId=2229017
	rt, err := m.GuestModel.UserGuestLogin(c, ctx, kentity.ReqPassportOpenLogin{
		OpenId:     tParam.DID,
		Type:       7,
		AppId:      appId,
		DeviceType: getDeviceType(appId),
		ExpireTime: constant.GuestLoginExpireTime * 60,
		ReqSrc:     m.ABSource,
		RegIp:      c.ClientIP(),
		Platform:   getPlatForm(appId),
	})
	if err == nil {
		kentity.SetCurrUser(c, rt)
	}
}

// UserValidate
//  @Description  用户合法校验
// 		- user_id >0 合法，否则不合法
//  @Receiver m
//  @Param c
func (m *Middleware) UserValidate(c *ginserver.TContext) {
	klog.Debug("用户合法验证 starting")
	snsNdkPublic := tcontext.GetPParam(c)

	myClaims := tcontext.GetMyClaims(c)
	if myClaims == nil {
		//c.AbortWithStatusJSON(http.StatusOK, &xentity.Response{Code: ecode.CodeUserNotLogin, Msg: "用户未登录", Data: nil})
		m.abortWithErrorJSON(c, errs.NewCustomError(ecode.CodeUserNotLogin, "用户未登录"))
		return
	}
	err := myClaims.GetError()
	if err != nil {
		if strings.Contains(err.Error(), "expired") {
			m.abortWithErrorJSON(c, errs.NewCustomError(ecode.CodeAuthTokenTimeout, "token已过期"))
			return
		}
		m.abortWithErrorJSON(c, errs.NewCustomError(ecode.CodeInvalidAuthToken, "解密token失败"))
		return
	}
	if snsNdkPublic == nil {
		m.abortWithErrorJSON(c, errs.NewCustomError(401, "权限校验失败"))
		return
	} else if !snsNdkPublic.UserTokenIsValid {
		m.abortWithErrorJSON(c, errs.NewCustomError(ecode.CodeUserNotLogin, "用户未登录"))
		return
	} else if snsNdkPublic.UserID == "" {
		m.abortWithErrorJSON(c, errs.NewCustomError(ecode.CodeUserNotLogin, "用户未登录"))
		return
	}
}

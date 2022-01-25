package tcontext

import (
	"context"
	"github.com/LuoHongLiang0921/kuaigo/pkg/biz/kentity"
	"github.com/LuoHongLiang0921/kuaigo/pkg/constant"
	"github.com/LuoHongLiang0921/kuaigo/pkg/core/server/kgin"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/kjwt"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/klog"
)

type headerName string

const (
	commonCtx headerName = "common"
)

type headerContextKey struct{}

type adminHeaderContextKey struct{}

type userContextKey struct{}

type KuaigoContext struct{}

// GetTabbyContext
// 	@Description:
//	@Param ctx
// 	@return *TabbyContext
// 	@return error
func GetTabbyContext(ctx *kgin.TContext) (*KuaigoContext, error) {
	return nil, nil
}

// GetHeader
//  @Description  获取用户请求头x
//  @Param c
//  @Return *xentity.TabbyHeader
func GetHeader(c *kgin.TContext) *kentity.KuaigoHeader {
	vo, flag := c.Get(constant.KeyHeader)
	if flag {
		snsNdkHeader, ok := vo.(*kentity.KuaigoHeader)
		if ok {
			return snsNdkHeader
		}
	}
	return nil
}

// GetPParam
//  @Description  获取p参数中的信息
//  @Param ctx
//  @Return *xentity.PParams
func GetPParam(ctx *kgin.TContext) *kentity.PParams {
	vo, flag := ctx.Get(constant.KeyPParam)
	if flag {
		snsNdkPublic, ok := vo.(*kentity.PParams)
		if ok {
			return snsNdkPublic
		}
	}
	return nil
}

// WithContext 从中间件中获取 common 日志
// 	@Description: 从中间件中获取 common 日志
//	@Param c
// 	@return context.Context
func WithContext(c *kgin.TContext) context.Context {
	ctx := context.Background()
	com, ok := c.Get(constant.CommonKey)
	if ok {
		ctx = klog.WithCommonLog(ctx, com)
	}

	header, ok := c.Get(constant.KeyHeader)
	if ok {
		ctx = context.WithValue(ctx, headerContextKey{}, header)
	}

	v, ok := c.Get(constant.KeyCurrUser)
	if ok {
		ctx = context.WithValue(ctx, userContextKey{}, v)
	}

	adminHeader, ok := c.Get(constant.KeyAdminHeader)
	if ok {
		ctx = context.WithValue(ctx, adminHeaderContextKey{}, adminHeader)
	}
	ctx = klog.RunningLoggerContext(ctx)
	return ctx
}

// GetCurrUser
//  @Description  获取当前登录用户信息
//  @Param c
//  @Return string
//  @Return int
//  @Return string
func GetCurrUser(c *kgin.TContext) (string, int, string) {
	vo, flag := c.Get(constant.KeyCurrUser)
	if flag {
		currUser, ok := vo.(*kjwt.MyClaims)
		if ok {
			return currUser.UserID, currUser.AppID, currUser.DID
		}
	}
	return "", 0, ""
}

// GetReqHeader
//  @Description  获取请求头header
//  @Param ctx
//  @Return *xentity.TabbyHeader
func GetReqHeader(ctx context.Context) *kentity.KuaigoHeader {
	v, ok := ctx.Value(headerContextKey{}).(*kentity.KuaigoHeader)
	if ok {
		return v
	}
	return &kentity.KuaigoHeader{}
}

// GetAdminHeader
//  @Description  获取后台接口请求头header
//  @Param ctx
//  @Return *xentity.AdminHeader
func GetAdminHeader(ctx context.Context) *kentity.AdminHeader {
	v, ok := ctx.Value(adminHeaderContextKey{}).(*kentity.AdminHeader)
	if ok {
		return v
	}
	return &kentity.AdminHeader{}
}

// GetAppID
//  @Description  获取AppID
//  @Param ctx
//  @Return int
func GetAppID(ctx context.Context) int {
	v, ok := ctx.Value(headerContextKey{}).(*kentity.KuaigoHeader)
	if ok {
		return v.AppId
	}
	return 0
}

// GetUserID
//  @Description  获取用户id
//  @Param ctx
//  @Return string
func GetUserID(ctx context.Context) string {
	v, ok := ctx.Value(userContextKey{}).(*kjwt.MyClaims)
	if ok {
		return v.UserID
	}
	return ""
}

// GetMyClaims
//  @Description  获取jwt claims
//  @Param ctx
//  @Return *xjwt.MyClaims
func GetMyClaims(ctx *kgin.TContext) *kjwt.MyClaims {
	if v, ok := ctx.Get(constant.KeyCurrUser); ok {
		mc, flag := v.(*kjwt.MyClaims)
		if flag {
			return mc
		}
	}
	return nil
}

// GetAbId 获取abId
// 	@Description: 获取abId
//	@Param ctx xgin.TContext
// 	@return string abId
func GetAbId(ctx *kgin.TContext) string {
	if v, ok := ctx.Get(constant.AbId); ok {
		mc, flag := v.(string)
		if flag {
			return mc
		}
	}
	return ""
}

// GetAbParam 获取abParam
// 	@Description:  获取abParam
//	@Param ctx xgin.TContext
// 	@return string abParam
func GetAbParam(ctx *kgin.TContext) string {
	if v, ok := ctx.Get(constant.AbParam); ok {
		mc, flag := v.(string)
		if flag {
			return mc
		}
	}
	return ""
}

// GetFUDid 获取fudid
// 	@Description GetFUDid 获取fudid
//	@Param ctx xgin.TContext
// 	@return string fudid
func GetFUDid(ctx *kgin.TContext) string {
	if fDid, err := ctx.Cookie(constant.FUDid); err == nil {
		return fDid
	}
	return ""
}

// CheckGuest 校验是否是游客
// 	@Description: 校验是否是游客 true是 false否
//	@Param ctx xgin.TContext
// 	@Return bool
func CheckGuest(ctx *kgin.TContext) bool {
	myClaims := GetMyClaims(ctx)
	if myClaims == nil || myClaims.LoginRole == 2 {
		return true
	}
	return false
}

// CheckRealUser 校验是否是真实用户
// 	@Description: 校验是否是游客 true是 false否
//	@Param ctx xgin.TContext
// 	@return bool
func CheckRealUser(ctx *kgin.TContext) bool {
	myClaims := GetMyClaims(ctx)
	if myClaims != nil && (myClaims.LoginRole == 0 || myClaims.LoginRole == 1) {
		return true
	}
	return false
}

package kcontext

import (
	"context"
	"github.com/LuoHongLiang0921/kuaigo/kpkg/biz/kentity"
	"github.com/LuoHongLiang0921/kuaigo/kpkg/constant"
	"github.com/LuoHongLiang0921/kuaigo/kpkg/core/server/kgin"
	"github.com/LuoHongLiang0921/kuaigo/kpkg/klog"
)

type headerName string

const (
	commonCtx headerName = "common"
)

type headerContextKey struct{}

type adminHeaderContextKey struct{}

type userContextKey struct{}

type KuaigoContext struct {
}

// GetTabbyContext
// 	@Description:
//	@param ctx
// 	@return *TabbyContext
// 	@return error
func GetTabbyContext(ctx *kgin.TContext) (*KuaigoContext, error) {
	return nil, nil
}

// GetHeader 获取用户请求头
// 	@Description: 获取用户请求头
//	@param c xgin.TContext
// 	@return *xentity.TabbyHeader
func GetHeader(c *kgin.TContext) *kentity.TabbyHeader {
	vo, flag := c.Get(constant.KeyHeader)
	if flag {
		snsNdkHeader, ok := vo.(*kentity.TabbyHeader)
		if ok {
			return snsNdkHeader
		}
	}
	return nil
}

// GetPParam 获取p参数中的信息
// 	@Description: 获取p参数中的信息
//	@param ctx
// 	@return *xentity.PParams p参数实体
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
//	@param c
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

	return ctx
}

// GetCurrUser 获取当前登录用户信息
// 	@Description: 获取当前登录用户信息
//	@param c
// 	@return string userId
// 	@return int appId
// 	@return string did
//func GetCurrUser(c *kgin.TContext) (string, int, string) {
//	vo, flag := c.Get(constant.KeyCurrUser)
//	if flag {
//		currUser, ok := vo.(*xjwt.MyClaims)
//		if ok {
//			return currUser.UserID, currUser.AppID, currUser.DID
//		}
//	}
//	return "", 0, ""
//}

//GetReqHeader 获取请求头header
// 	@Description: 获取请求头header
//	@param ctx
// 	@return *xentity.TabbyHeader
func GetReqHeader(ctx context.Context) *kentity.TabbyHeader {
	v, ok := ctx.Value(headerContextKey{}).(*kentity.TabbyHeader)
	if ok {
		return v
	}
	return &kentity.TabbyHeader{}
}

//GetAdminHeader 获取后台接口请求头header
// 	@Description: 获取后台接口请求头header
//	@param ctx
// 	@return *xentity.AdminHeader
func GetAdminHeader(ctx context.Context) *kentity.AdminHeader {
	v, ok := ctx.Value(adminHeaderContextKey{}).(*kentity.AdminHeader)
	if ok {
		return v
	}
	return &kentity.AdminHeader{}
}

// GetAppID 获取AppID
// 	@Description: 获取AppID
//	@param ctx
// 	@return int appId
func GetAppID(ctx context.Context) int {
	v, ok := ctx.Value(headerContextKey{}).(*kentity.TabbyHeader)
	if ok {
		return v.AppId
	}
	return 0
}

// GetUserID 获取用户id
// 	@Description:  获取用户id
//	@param ctx
// 	@return string userId
//func GetUserID(ctx context.Context) string {
//	v, ok := ctx.Value(userContextKey{}).(*xjwt.MyClaims)
//	if ok {
//		return v.UserID
//	}
//	return ""
//}

// GetMyClaims 获取jwt claims
// 	@Description: 获取jwt claims
//	@param ctx xgin.TContext
// 	@return *xjwt.MyClaims
//func GetMyClaims(ctx *kgin.TContext) *xjwt.MyClaims {
//	if v, ok := ctx.Get(constant.KeyCurrUser); ok {
//		mc, flag := v.(*xjwt.MyClaims)
//		if flag {
//			return mc
//		}
//	}
//	return nil
//}

// GetAbId 获取abId
// 	@Description: 获取abId
//	@param ctx xgin.TContext
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
//	@param ctx xgin.TContext
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
// 	@Description:  // GetFUDid 获取fudid
//	@param ctx xgin.TContext
// 	@return string fudid
func GetFUDid(ctx *kgin.TContext) string {
	if fDid, err := ctx.Cookie(constant.FUDid); err == nil {
		return fDid
	}
	return ""
}

// CheckGuest 校验是否是游客
// 	@Description: 校验是否是游客 true是 false否
//	@param ctx xgin.TContext
// 	@return bool
//func CheckGuest(ctx *kgin.TContext) bool {
//	myClaims := GetMyClaims(ctx)
//	if myClaims == nil || myClaims.LoginRole == 2 {
//		return true
//	}
//	return false
//}

// CheckRealUser 校验是否是真实用户
// 	@Description: 校验是否是游客 true是 false否
//	@param ctx xgin.TContext
// 	@return bool
//func CheckRealUser(ctx *xgin.TContext) bool {
//	myClaims := GetMyClaims(ctx)
//	if myClaims != nil && (myClaims.LoginRole == 0 || myClaims.LoginRole == 1) {
//		return true
//	}
//	return false
//}

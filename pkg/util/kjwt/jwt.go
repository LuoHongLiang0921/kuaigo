// @Description json web token 工具类

package kjwt

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Client struct {
	// JwtKey
	JwtKey string
	// jwt expire
	Expire time.Duration
}

// NewClient
//  @Description  创建一个jwt Token
//  @Param jwtKey
//  @Param expire
//  @Return *Client
func NewClient(jwtKey string, expire time.Duration) *Client {
	return &Client{JwtKey: jwtKey, Expire: expire}
}

// MyClaims ...
// todo 需要解耦改为interface
type MyClaims struct {
	UserID    string `json:"userID"`
	AppID     int    `json:"appID"`
	DID       string `json:"did"`
	PlatToken string `json:"platToken"` // 兼容 bobo 旧服务， 新框架使用的是新的token
	LoginRole int    `json:"loginRole"` // 登录角色 0 游客， 1 真实用户
	error     error  `json:"-"`
	jwt.StandardClaims
}

//SetError
func (c *MyClaims) SetError(err error) {
	c.error = err
}

// GetError
func (c *MyClaims) GetError() error {
	return c.error
}

//CreateToken
// CreateToken 创建token
// 	@Description  创建token
// 	@receiver c
//	@Param claims
// 	@return string token 字符串
// 	@return error
func (c *Client) CreateToken(claims *MyClaims) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims.ExpiresAt = time.Now().Add(c.Expire).Unix()
	claims.IssuedAt = time.Now().Unix()
	claims.Issuer = claims.DID
	token.Claims = claims
	return token.SignedString([]byte(c.JwtKey))
}

// ParseToken
//  @Description  解析token
//  @Receiver c
//  @Param tokenStr token 字符串
//  @Return *MyClaims
//  @Return error
func (c *Client) ParseToken(tokenStr string) (*MyClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(c.JwtKey), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}

}

// RefreshToken
//  @Description  刷新token
//  @Receiver c
//  @Param tokenStr token 字符串
//  @Return string
func (c *Client) RefreshToken(tokenStr string) string {
	claims, err := c.ParseToken(tokenStr)
	if err != nil {
		return ""
	}
	tokenStr, err = c.CreateToken(claims)
	if err != nil {
		return ""
	}
	return tokenStr
}

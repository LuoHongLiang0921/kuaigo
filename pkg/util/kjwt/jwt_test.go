// @Description

package kjwt

import (
	"fmt"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func TestClient_CreateToken(t *testing.T) {
	client := NewClient("f165058924942e1bd19ce6ded681984c", 300*time.Hour)
	aaa, _ := client.CreateToken(&MyClaims{
		UserID:         "178319450342137920", //common._uId
		AppID:          104,                  //bobo 101
		DID:            "did",                //设备id，暂时不用
		PlatToken:      "",                   //common._token
		LoginRole:      2,                    // 0,1： 真实登录， 2：游客登录（0 为了兼容老版本用户）
		StandardClaims: jwt.StandardClaims{},
	})

	fmt.Println(aaa)

	bbb, err := client.ParseToken(aaa)
	fmt.Println(err)
	fmt.Println(bbb)
}

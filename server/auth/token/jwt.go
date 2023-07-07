//定义了生成JWT令牌的结构体JWTTokenGen和生成令牌的方法实现

package token

import (
	"crypto/rsa"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// JWTTokenGen 是用于生成 JWT 令牌的结构体
type JWTTokenGen struct {
	issuer     string   // 发行者
	nowFunc    func() time.Time // 当前时间函数
	privateKey *rsa.PrivateKey  // 私钥
}

//创建一个 JWTTokenGen 实例
func NewJWTTokenGen(issuer string, privateKey *rsa.PrivateKey) *JWTTokenGen {
	return &JWTTokenGen{
		issuer:     issuer,
		nowFunc:    time.Now,
		privateKey: privateKey,
	}
}

func (t *JWTTokenGen) GenerateToken(accountID string, expire time.Duration) (string, error) {
	nowSec := t.nowFunc().Unix()

	tkn := jwt.NewWithClaims(jwt.SigningMethodRS512, jwt.StandardClaims{
		Issuer:    t.issuer,
		IssuedAt:  nowSec,
		ExpiresAt: nowSec + int64(expire.Seconds()),
		Subject:   accountID,     // 主题（账户ID）
	})

	return tkn.SignedString(t.privateKey)
}

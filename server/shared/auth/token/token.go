package token

import (
	"crypto/rsa"
	"fmt"

	"github.com/dgrijalva/jwt-go"
)

type JWTTokenVerifier struct {
	PublicKey  *rsa.PublicKey
}

func (v *JWTTokenVerifier) Verify(token string) (string, error) {

	//解析 JWT 令牌。指定了令牌的格式为 jwt.StandardClaims，并传递一个回调函数，该函数返回公钥用于验证签名。
	t, err := jwt.ParseWithClaims(token, &jwt.StandardClaims{},
		func(*jwt.Token) (interface{}, error) {
			return v.PublicKey, nil
		})
	if err != nil {
		return "", fmt.Errorf("cannot parse token :%v", err)
	}

	//检查解析后的令牌是否有效
	if !t.Valid {
        return "", fmt.Errorf("token not valid : %v", err)
	}

	//将解析后的令牌的声明转换为 jwt.StandardClaims 类型
	clm, ok := t.Claims.(*jwt.StandardClaims)
	if !ok{
        return "", fmt.Errorf("token claims is not StandardClaims: %v",ok)
	}

	//检查声明是否有效
	if err = clm.Valid(); err !=nil{
		return "", fmt.Errorf("claim not valid : %v",err)
	}

	//返回令牌中的主题作为验证结果
	return clm.Subject, nil
}

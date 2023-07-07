package auth

import (
	"context"
	"fmt"
	"io"
	"os"
	"rentCar/server/shared/auth/token"
	"rentCar/server/shared/id"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	authorizationHeader = "authorization"
	bearerPrefix        = "bearer "
)

// Interceptor create grpc auth interceptor:创建grpc认证拦截器
func Interceptor(publicKeyFile string) (grpc.UnaryServerInterceptor, error) {
	//获取文件
	f, err := os.Open(publicKeyFile)
	if err != nil {
		return nil, fmt.Errorf("cannot open publickeyfile err : %v", err)
	}

	//读取文件
	b, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("cannot read publickeyfile err : %v", err)
	}

	//解析文件
	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(b)
	if err != nil {
		return nil, fmt.Errorf("cannot parse publickeyfile bytes err : %v", err)
	}

	i := &interceptor{
		verifier: &token.JWTTokenVerifier{
			PublicKey: pubKey,
		},
	}
	return i.HandleReq, nil
}

type tokenVerifier interface {
	Verify(token string) (string, error)
}

type interceptor struct {
	verifier tokenVerifier
}

// HandleReq是grpc认证拦截器的处理函数
func (i *interceptor) HandleReq(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	tkn, err := tokenFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "")
	}

	aid, err := i.verifier.Verify(tkn)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "token not valid : %v", err)
	}

	//将账户ID添加到上下文中，并继续处理请求
	return handler(ContextWithAccountID(ctx, id.AccountID(aid)), req)
}

// 从上下文中获取令牌
func tokenFromContext(c context.Context) (string, error) {
	// 从上下文中提取元数据
	m, ok := metadata.FromIncomingContext(c)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "")
	}

	tkn := ""
	//在授权头部中查找令牌
	for _, v := range m[authorizationHeader] {
		if strings.HasPrefix(v, bearerPrefix) {
			tkn = v[len(bearerPrefix):]
		}
	}

	if tkn == "" {
		return "", status.Error(codes.Unauthenticated, "")
	}

	return tkn, nil
}

type accountIDKey struct{}



// ContextWithAccountID 创建带有账户ID的上下文
func ContextWithAccountID(c context.Context, aid id.AccountID) context.Context {
	return context.WithValue(c, accountIDKey{}, aid)
}

// 从上下文中获取accountID
func AccountIDFromContext(c context.Context) (id.AccountID, error) {
	v := c.Value(accountIDKey{})
	aid, ok := v.(id.AccountID)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "")
	}
	return id.AccountID(aid), nil
}

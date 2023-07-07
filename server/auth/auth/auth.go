//定义了认证服务的结构体Service和相关方法的实现

package auth

import (
	"context"
	authpb "rentCar/server/auth/api/gen/v1"
	"rentCar/server/auth/dao"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Service 定义认证服务结构体
type Service struct {
	Logger         *zap.Logger
	OpenResolver   OpenResolver
	Mongo          *dao.Mongo
	TokenGenerator TokenGenerator
	TokenExpire    time.Duration
}

// OpenResolver 定义开放授权解析器接口
type OpenResolver interface {
	Resolve(code string) (string, error)
}

// TokenGenerator 定义令牌生成器接口
type TokenGenerator interface {
	GenerateToken(accountID string, expire time.Duration) (string, error)
}

// Login 实现登录功能
func (s *Service) Login(c context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	s.Logger.Info("received code", zap.String("code", req.Code))

	// 解析授权码获取 OpenID
	openID, err := s.OpenResolver.Resolve(req.Code)
	if err != nil {
		return nil, status.Errorf(codes.Unavailable, "cannot resolve openID err: %v", err)
	}

	// 解析 OpenID 获取账户ID
	accountID, err := s.Mongo.ResolveAccountID(c, openID)
	if err != nil {
		s.Logger.Error("cannot resolve account id ", zap.Error(err))
		return nil, status.Error(codes.Internal, "")
	}

	// 生成访问令牌
	tkn, err := s.TokenGenerator.GenerateToken(accountID.String(), s.TokenExpire)
	if err != nil {
		s.Logger.Error("cannot generate token", zap.Error(err))
		return nil, status.Error(codes.Internal, "")
	}

	return &authpb.LoginResponse{
		AccessToken: tkn,
		ExpiresIn:   int32(s.TokenExpire.Seconds()),
	}, nil
}

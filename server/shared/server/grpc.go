package server

import (
	"net"
	"rentCar/server/shared/auth"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type GRPCConfig struct {
	Name              string
	Addr              string
	AuthPublicKeyFile string
	Logger            *zap.Logger
	RegisterFunc      func(*grpc.Server)
}

func RunGRPCServer(c *GRPCConfig) error {
	nameField := zap.String("name", c.Name)
	lis, err := net.Listen("tcp", c.Addr)
	if err != nil {
		c.Logger.Fatal("net listen err :%v", nameField, zap.Error(err))
	}

	var opts []grpc.ServerOption
	if c.AuthPublicKeyFile != "" {
		in, err := auth.Interceptor(c.AuthPublicKeyFile)
		if err != nil {
			c.Logger.Fatal("cannot create auth interceptor", nameField, zap.Error(err))
		}
		opts = append(opts, grpc.UnaryInterceptor(in))
	}

	//创建grpc服务器实例
	s := grpc.NewServer(opts...)

	//注册认证服务
	c.RegisterFunc(s)

	c.Logger.Info("server started ", nameField, zap.String("addr", c.Addr))
	return s.Serve(lis)
}

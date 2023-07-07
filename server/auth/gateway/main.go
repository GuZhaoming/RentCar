package main

import (
	"context"
	"log"
	"net/http"
	authpb "rentCar/server/auth/api/gen/v1"
	rentalpb "rentCar/server/rental/api/gen/v1"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("cannot create logger :%v", err)
	}
	c := context.Background()
	c, cancel := context.WithCancel(c)
	defer cancel()

	mux := runtime.NewServeMux(runtime.WithMarshalerOption(
		runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				UseProtoNames:  true,
				UseEnumNumbers: true,
			},
		},
	))

	serverConfig := []struct {
		name         string
		addr         string
		registerFunc func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error)
	}{
		{
			name:         "auth",
			addr:         "localhost:8081",
			registerFunc: authpb.RegisterAuthServiceHandlerFromEndpoint,
		},
		{
			name:         "rental",
			addr:         "localhost:8082",
			registerFunc: rentalpb.RegisterTripServiceHandlerFromEndpoint,
		},
	}

	//将 gRPC 服务绑定到 HTTP/JSON 接口
	for _, s := range serverConfig {
		err := s.registerFunc(
			c,
			mux,
			s.addr,
			//创建不安全的传输凭证进行连接。
			[]grpc.DialOption{
				grpc.WithTransportCredentials(insecure.NewCredentials()),
			})

		if err != nil {
			log.Fatalf("cannot register gateway server: %s %v", s.name, err)
		}
	}
	addr := ":8080"
    logger.Sugar().Infof("grpc gateway started at : %s",addr)
	//启动 HTTP 服务器，将转换后的 HTTP/JSON 接口提供给客户端访问。
	logger.Sugar().Fatal(http.ListenAndServe(":8080", mux))
}

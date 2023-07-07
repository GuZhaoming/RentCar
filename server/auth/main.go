// 文件是程序的入口，负责初始化和启动gRPC服务器，并注册认证服务
package main

import (
	"context"
	"io"
	"log"
	"os"
	authpb "rentCar/server/auth/api/gen/v1"
	"rentCar/server/auth/auth"
	"rentCar/server/auth/dao"
	"rentCar/server/auth/token"
	"rentCar/server/auth/wechat"
	"rentCar/server/shared/server"
	"time"
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {

	//创建日志记录器
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("cannot create logger :%v", err)
	}

	c := context.Background()
	//连接mongoDB
	MongoClient, err := mongo.Connect(c, options.Client().ApplyURI("mongodb://localhost:27017/rentCar?readPreference=primary&ssl=false&directConnection=true"))
	if err != nil {
		logger.Fatal("cannot connect mongo :", zap.Error(err))
	}

	//打开私钥文件
	pkFile, err := os.Open("auth/private.key")
	if err != nil {
		logger.Fatal("cannot open private key", zap.Error(err))
	}

	//读取私钥文件的内容
	pkBytes, err := io.ReadAll(pkFile)
	if err != nil {
		logger.Fatal("cannot read private key ", zap.Error(err))
	}
    
	//解析私钥
	privKey, err := jwt.ParseRSAPrivateKeyFromPEM(pkBytes)
	if err != nil {
		logger.Fatal("cannot parse private key", zap.Error(err))
	}

	logger.Sugar().Fatal(
	server.RunGRPCServer(&server.GRPCConfig{
		Name: "auth",
		Addr: ":8081",
		Logger: logger,
		RegisterFunc: func(s *grpc.Server) {
			authpb.RegisterAuthServiceServer(s, &auth.Service{
				OpenResolver: &wechat.Service{
					AppID:     "wx85aad7821638453a",
					AppSecret: "d53ecabb8bc730b5d12697478ba7f068",
				},
				Logger:         logger,
				Mongo:          dao.NewMongo(MongoClient.Database("rentCar")),
				TokenExpire:    2 * time.Hour,
				TokenGenerator: token.NewJWTTokenGen("rentCar/auth", privKey),
			})
		},
	}))
}

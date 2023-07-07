// 文件是程序的入口，负责初始化和启动gRPC服务器，并注册认证服务
package main

import (
	"context"
	"log"
	rentalpb "rentCar/server/rental/api/gen/v1"
	"rentCar/server/rental/trip"
	"rentCar/server/rental/trip/client/car"
	"rentCar/server/rental/trip/client/poi"
	"rentCar/server/rental/trip/client/profile"
	"rentCar/server/rental/trip/dao"
	"rentCar/server/shared/server"

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
	
	logger.Sugar().Fatal(
		server.RunGRPCServer(&server.GRPCConfig{
			Name:              "rental",
			Addr:              ":8082",
			AuthPublicKeyFile: "shared/auth/public.key",
			Logger:            logger,
			RegisterFunc: func(s *grpc.Server) {
				rentalpb.RegisterTripServiceServer(s, &trip.Service{
					Logger: logger,
					CarManager: &car.Manager{},
					ProfileManager: &profile.Manager{},
					POIManager: &poi.Manager{},
					Mongo: dao.NewMongo(MongoClient.Database("rentCar")),
				})},
		}),
	)
}

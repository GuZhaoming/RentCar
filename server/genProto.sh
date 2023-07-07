function genProto {
    DOMAIN=$1
    PROTO_PATH=./${DOMAIN}/api
    GO_OUT_PATH=./${DOMAIN}/api/gen/v1

    # mkdir -p $GO_OUT_PATH

    # protoc --go_out=$GO_OUT_PATH --go_opt=paths=source_relative $PROTO_PATH/${DOMAIN}.proto
    # protoc --go-grpc_out=$GO_OUT_PATH --go-grpc_opt=paths=source_relative $PROTO_PATH/${DOMAIN}.proto
    # protoc --grpc-gateway_out=paths=source_relative,grpc_api_configuration=./${DOMAIN}/api/${DOMAIN}.yaml:./${DOMAIN}/api/gen/v1 ./${DOMAIN}/api/${DOMAIN}.proto
   
    # 移动生成的文件到目标文件夹
    # mv $GO_OUT_PATH/${DOMAIN}/api/* $GO_OUT_PATH/
    # rm -rf $GO_OUT_PATH/${DOMAIN}

    PBJS_BIN_PATH=../wx/miniprogram/node_modules/.bin
    PBJS_OUT_PATH=../wx/miniprogram/service/proto_gen/${DOMAIN}

    mkdir -p $PBJS_OUT_PATH

    $PBJS_BIN_PATH/pbjs $PROTO_PATH/${DOMAIN}.proto --es6 $PBJS_OUT_PATH/${DOMAIN}_pb.js 

    $PBJS_BIN_PATH/pbts $PROTO_PATH/${DOMAIN}.proto --ts $PBJS_OUT_PATH/${DOMAIN}_pb.ts
}

genProto auth
genProto rental  



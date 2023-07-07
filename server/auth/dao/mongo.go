package dao

import (
	"context"
	"fmt"
	"rentCar/server/shared/id"
	mgutil "rentCar/server/shared/mongo"
	"rentCar/server/shared/mongo/objid"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const openIDField = "open_id"

// Mongo 是一个 MongoDB 数据访问对象
type Mongo struct {
	car *mongo.Collection // MongoDB 集合对象
}

// NewMongo 创建一个新的 Mongo 实例
func NewMongo(db *mongo.Database) *Mongo {
	return &Mongo{
		car: db.Collection("account"), // 获取 "account" 集合
	}
}

// ResolveAccountID 根据 openID 解析对应的账户ID
func (m *Mongo) ResolveAccountID(c context.Context, openID string) (id.AccountID, error) {

	// 生成新的 ObjectID
	insertedID := mgutil.NewObjID()
	// 在 "account" 集合中查找并更新符合条件的文档
	res := m.car.FindOneAndUpdate(
		c,
		bson.M{openIDField: openID}, // 根据 openID 进行查询
		mgutil.SetOnInsert(bson.M{ // 如果找不到匹配的文档，则插入以下数据
			mgutil.IDFieldName: insertedID, // 设置 ID 字段为新的 ObjectID
			openIDField: openID,     // 设置 OpenID 字段为传入的 openID
		}),
		options.FindOneAndUpdate().
			SetUpsert(true).                  // 如果找不到匹配的文档，则插入新文档
			SetReturnDocument(options.After)) // 返回更新后的文档

	// 检查操作是否出错
	if err := res.Err(); err != nil {
		return "", fmt.Errorf("cannot FindOneAndUpdate: %v", err)
	}

	// 解码结果并存储在 row 变量中
	var row mgutil.IDField
	err := res.Decode(&row)
	if err != nil {
		return "", fmt.Errorf("cannot Decode result: %v", err)
	}

	// 返回解析后的账户ID
	//return row.ID.Hex(), nil
	return objid.ToAccountID(row.ID), nil
}

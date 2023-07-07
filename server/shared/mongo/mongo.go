package mgutil

import (
	"fmt"
	"rentCar/server/shared/mongo/objid"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	IDFieldName       = "_id"
	UpdateAtFieldName = "updatedat"
)

type IDField struct {
	ID primitive.ObjectID `bson:"_id"`
}

type UpdateAtField struct {
	UpdatedAt int64 `bson:"updatedat"`
}

var NewObjID = primitive.NewObjectID

func Set(v interface{}) bson.M {
	return bson.M{
		"$set": v,
	}
}

func SetOnInsert(v interface{}) bson.M {
	return bson.M{
		"$setOnInsert": v,
	}
}

var UpdatedAt = func() int64 {
	return time.Now().UnixNano()
}


func NewObjIDWithValue(id fmt.Stringer) {
	NewObjID = func() primitive.ObjectID {
		return objid.MustFromID(id)
	}
}
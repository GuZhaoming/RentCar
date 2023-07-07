package dao

import (
	"context"
	"fmt"
	rentalpb "rentCar/server/rental/api/gen/v1"
	"rentCar/server/shared/id"
	mgutil "rentCar/server/shared/mongo"
	"rentCar/server/shared/mongo/objid"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	tripField      = "trip"
	accountIDField = tripField + ".accountid"
	statusField    = tripField + ".status"
)

type Mongo struct {
	car *mongo.Collection
}

func NewMongo(db *mongo.Database) *Mongo {
	return &Mongo{
		car: db.Collection("trip"),
	}
}

type TripRecord struct {
	mgutil.IDField       `bson:"inline"`
	mgutil.UpdateAtField `bson:"inline"`
	Trip                 *rentalpb.Trip `bson:"trip"`
}

func (m *Mongo) CreateTrip(c context.Context, trip *rentalpb.Trip) (*TripRecord, error) {

	r := &TripRecord{
		Trip: trip,
	}
	r.ID = mgutil.NewObjID()
	r.UpdatedAt = mgutil.UpdatedAt()

	_, err := m.car.InsertOne(c, r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (m *Mongo) GetTrip(c context.Context, id id.TripID, accountID id.AccountID) (*TripRecord, error) {
	objID, err := objid.FromID(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id : %v", err)
	}
	res := m.car.FindOne(c, bson.M{
		mgutil.IDFieldName: objID,
		accountIDField:     accountID,
	})
	if err := res.Err(); err != nil {
		return nil, err
	}

	var tr TripRecord
	err = res.Decode(&tr)
	if err != nil {
		return nil, fmt.Errorf("cannot decode %v", err)
	}
	return &tr, nil
}

func (m *Mongo) GetTrips(c context.Context, accountID id.AccountID, status rentalpb.TripStatus) ([]*TripRecord, error) {
	filter := bson.M{
		accountIDField: accountID.String(),
	}
	if status != rentalpb.TripStatus_TS_NOT_SPECIFIED {
		filter[statusField] = status
	}

	res, err := m.car.Find(c, filter)
	if err != nil {
		return nil, err
	}

	var trips []*TripRecord

	for res.Next(c) {
		var trip TripRecord
		err := res.Decode(&trip)
		if err != nil {
			return nil, err
		}
		trips = append(trips, &trip)
	}
	fmt.Printf("----------------结果%v", trips)
	return trips, nil
}

func (m *Mongo) UpdateTrip(c context.Context, tid id.TripID, aid id.AccountID, updatedAt int64, trip *rentalpb.Trip) error {
	objID, err := objid.FromID(tid) // 将TripID转换为MongoDB的ObjectID类型
	if err != nil {
		return fmt.Errorf("错误invalid id : %v", err) // 如果转换失败，则返回错误信息
	}
	newUpdatedAt := mgutil.UpdatedAt() // 获取当前时间戳作为更新时间
	res, err := m.car.UpdateOne(c, bson.M{ // 在MongoDB的"car"集合中查找满足条件的文档
		mgutil.IDFieldName:       objID, // 根据ObjectID进行查找
		accountIDField:           aid.String(), // 根据AccountID进行查找
		mgutil.UpdateAtFieldName: updatedAt, // 根据更新时间进行查找
	}, mgutil.Set(bson.M{
		tripField:                trip, // 将传入的trip参数更新到对应的文档中
		mgutil.UpdateAtFieldName: newUpdatedAt, // 更新更新时间字段为当前时间戳
	}))
	if err != nil{
		return err // 如果更新过程中出现错误，则返回错误信息
	}
	if res.MatchedCount == 0{
		return mongo.ErrNoDocuments // 如果没有匹配到任何文档，则返回"无文档"错误
	}
	return nil // 更新成功，返回nil
}

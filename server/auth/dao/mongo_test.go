package dao

import (
	"context"
	"os"
	"rentCar/server/shared/id"
	mgutil "rentCar/server/shared/mongo"
	"rentCar/server/shared/mongo/objid"
	mongotesting "rentCar/server/shared/mongo/testing"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
)


// TestResolveAccountID 是一个测试函数，用于测试 ResolveAccountID 方法
func TestResolveAccountID(t *testing.T) {
	c := context.Background()
	mc, err := mongotesting.NewClient(c)
	if err != nil {
		t.Fatalf("cannot connect mongo %v", err)
	}

	m := NewMongo(mc.Database("rentCar"))

	// 在数据库中插入初始值
	_, err = m.car.InsertMany(c, []interface{}{
		bson.M{
			mgutil.IDFieldName:    objid.MustFromID(id.AccountID("648badb1b4600ccb093e5b29")) ,
			openIDField: "openid_1",
		},
		bson.M{
			mgutil.IDFieldName:     objid.MustFromID(id.AccountID("648badb1b4600ccb093e5b2a")),
			openIDField: "openid_2",
		},
	})

	if err != nil {
		t.Fatalf("cannot insert initial values : %v", err)
	}

	// 自定义生成新的 ObjectID 的函数
    mgutil.NewObjIDWithValue(id.AccountID("648badb1b4600ccb093e5b03"))


	cases := []struct {
		name   string
		openID string
		want   string
	}{
		{
			name:   "existing_user",
			openID: "openid_1",
			want:   "648badb1b4600ccb093e5b29",
		},
		{
			name:   "another_existing_user",
			openID: "openid_2",
			want:   "648badb1b4600ccb093e5b2a",
		},
		{
			name:   "new_user",
			openID: "openod_3",
			want:   "648badb1b4600ccb093e5b03",
		},
	}

	for _, cc := range cases {
		//运行一个子测试
		t.Run(cc.name, func(t *testing.T) {
			id, err := m.ResolveAccountID(context.Background(), cc.openID)
			if err != nil {
				t.Errorf("faild resolve account id for %q :%v", cc.openID, err)
			}
			if id.String() != cc.want {
				t.Errorf("resolve account id : want:%q , got:%q", cc.want, id)
			}
		})
	}

}



// TestMain 是测试的入口函数，负责启动 MongoDB 容器，并执行测试
func TestMain(m *testing.M) {
	os.Exit(mongotesting.RunWithMongoInDocker(m))
}

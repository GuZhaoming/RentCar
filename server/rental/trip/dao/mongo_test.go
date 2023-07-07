package dao

import (
	"context"
	"os"
	rentalpb "rentCar/server/rental/api/gen/v1"
	"rentCar/server/shared/id"
	mgutil "rentCar/server/shared/mongo"
	"rentCar/server/shared/mongo/objid"
	mongotesting "rentCar/server/shared/mongo/testing"
	"testing"

	"github.com/google/go-cmp/cmp"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestCreateTrip(t *testing.T) {
	c := context.Background()
	mc, err := mongotesting.NewClient(c)
	if err != nil {
		t.Fatalf("cannot connect mongodb %v", err)
	}

	db := mc.Database("rentCar")
	err = mongotesting.SetupIndexes(c, db)
	if err != nil {
		t.Fatalf("cannot setup indexes :%v", err)
	}

	m := NewMongo(db)

	cases := []struct {
		name       string
		tripID     id.TripID
		accountID  id.AccountID
		tripStatus rentalpb.TripStatus
		wantErr    bool
	}{
		{
			name:       "finished",
			tripID:     id.TripID("648badb1b4600ccb093e5b29"),
			accountID:  id.AccountID("account1"),
			tripStatus: rentalpb.TripStatus_FINISHED,
		},
		{
			name:       "another_finished",
			tripID:     id.TripID("648badb1b4600ccb093e5b21"),
			accountID:  id.AccountID("account1"),
			tripStatus: rentalpb.TripStatus_FINISHED,
		},
		{
			name:       "in_progress",
			tripID:     id.TripID("648badb1b4600ccb093e5b22"),
			accountID:  id.AccountID("account1"),
			tripStatus: rentalpb.TripStatus_IN_PROGRESS,
		},
		{
			name:       "anther_in_progress",
			tripID:     id.TripID("648badb1b4600ccb093e5b23"),
			accountID:  id.AccountID("account1"),
			tripStatus: rentalpb.TripStatus_IN_PROGRESS,
			wantErr:    true,
		},
		{
			name:       "in_progress_by_another_account",
			tripID:     id.TripID("648badb1b4600ccb093e5b24"),
			accountID:  id.AccountID("account2"),
			tripStatus: rentalpb.TripStatus_FINISHED,
		},
	}

	for _, cc := range cases {
		mgutil.NewObjIDWithValue(id.TripID(cc.tripID))
		tr, err := m.CreateTrip(c, &rentalpb.Trip{
			AccountId: string(cc.accountID),
			Status:    cc.tripStatus,
		})
		if cc.wantErr {
			if err == nil {
				t.Errorf("%s--error expected;got none", cc.name)
			}
			continue
		}
		if err != nil {
			t.Errorf("%s--error creating trip :%v", cc.name, err)
			continue
		}
		if tr.ID.Hex() != string(cc.tripID) {
			t.Errorf("%s--incorrect trip id; want-------: %q ; got------:%q", cc.name, cc.tripID, tr.ID.Hex())
		}
	}
}

func TestGetTrip(t *testing.T) {

	c := context.Background()
	mc, err := mongotesting.NewDefaultClient(c)
	if err != nil {
		t.Fatalf("cannot connect mongo %v", err)
	}

	m := NewMongo(mc.Database("rentCar"))
	acct := id.AccountID("account1")
	mgutil.NewObjID = primitive.NewObjectID
	tr, err := m.CreateTrip(c, &rentalpb.Trip{
		AccountId: acct.String(),
		CarId:     "car1",
		Start: &rentalpb.LocationStatus{
			PoiName: "startpoint",
			Location: &rentalpb.Location{
				Latitude:  30,
				Longitude: 120,
			},
		},
		End: &rentalpb.LocationStatus{
			PoiName:  "endpoint",
			FeeCent:  10000,
			KmDriven: 35,
			Location: &rentalpb.Location{
				Latitude:  35,
				Longitude: 121,
			},
		},
		Status: rentalpb.TripStatus_FINISHED,
	})
	if err != nil {
		t.Fatalf("cannot create trip err :%v", err)
	}

	got, err := m.GetTrip(c, objid.ToTripID(tr.ID), acct)
	if err != nil {
		t.Errorf("cannot get trip :%v", err)
	}

	if diff := cmp.Diff(tr, got, protocmp.Transform()); diff != "" {
		t.Errorf("result differs; -want +got:%s", diff)
	}
}

func TestGetTrips(t *testing.T) {
	rows := []struct {
		id        string
		accountID string
		status    rentalpb.TripStatus
	}{
		{
			id:        "64a27be7737b04d125af1a45",
			accountID: "account_id_for_get_trips",
			status:    rentalpb.TripStatus_FINISHED,
		},
		{
			id:        "648badb1b4600ccb093e5b24",
			accountID: "account_id_for_get_trips",
			status:    rentalpb.TripStatus_FINISHED,
		},
		{
			id:        "648badb1b4600ccb093e5b25",
			accountID: "account_id_for_get_trips",
			status:    rentalpb.TripStatus_FINISHED,
		},
		{
			id:        "648badb1b4600ccb093e5b26",
			accountID: "account_id_for_get_trips",
			status:    rentalpb.TripStatus_IN_PROGRESS,
		},
	}
	c := context.Background()
	mc, err := mongotesting.NewClient(c)
	if err != nil {
		t.Fatalf("cannot connect mongodb %v", err)
	}
	m := NewMongo(mc.Database("rentCar"))

	for _, r := range rows {
		mgutil.NewObjIDWithValue(id.TripID(r.id))
		_, err := m.CreateTrip(c, &rentalpb.Trip{
			AccountId: r.accountID,
			Status:    r.status,
		})
		if err != nil {
			t.Fatalf("cannot create rows:%v", err)
		}
	}

	cases := []struct {
		name       string
		accountID  string
		status     rentalpb.TripStatus
		wantCount  int
		wantOnlyID string
	}{
		{
			name:      "get_all",
			accountID: "account_id_for_get_trips",
			status:    rentalpb.TripStatus_TS_NOT_SPECIFIED,
			wantCount: 4,
		}, {
			name:       "get_one",
			accountID:  "account_id_for_get_trips",
			status:     rentalpb.TripStatus_IN_PROGRESS,
			wantCount:  1,
			wantOnlyID: "648badb1b4600ccb093e5b26",
		},
	}

	for _, cc := range cases {
		t.Run(cc.name, func(t *testing.T) {
			res, err := m.GetTrips(context.Background(), id.AccountID(cc.accountID), cc.status)
			if err != nil {
				t.Errorf("cannot get trips ,%v ", err)
			}
			if cc.wantCount != len(res) {
				t.Errorf("incorrect result count; want:%d,got:%d", cc.wantCount, len(res))
			}

			if cc.wantOnlyID != "" && len(res) > 0 {
				if cc.wantOnlyID != res[0].ID.Hex() {
					t.Errorf("only_id incorrect; want:%q, got:%q", cc.wantOnlyID, res[0].ID.Hex())
				}
			}
		})
	}
}

func TestUpdateTrip(t *testing.T) {
	c := context.Background()
	mc, err := mongotesting.NewClient(c)
	if err != nil {
		t.Fatalf("cannot connect mongodb %v", err)
	}

	m := NewMongo(mc.Database("rentCar"))
	tid := id.TripID("648badb1b4600ccb093e5b29")
	aid := id.AccountID("account_for_update")

	var now int64 = 10000
	mgutil.NewObjIDWithValue(tid)
	mgutil.UpdatedAt = func() int64 {
		return now
	}

	// 创建一个新的行程
	tr, err := m.CreateTrip(c, &rentalpb.Trip{
		AccountId: aid.String(),
		Status:    rentalpb.TripStatus_IN_PROGRESS,
		Start: &rentalpb.LocationStatus{
			PoiName: "start_poi",
		},
	})
	if err != nil {
		t.Fatalf("cannot create trip :%v", err)
	}

	if tr.UpdatedAt != 10000 {
		t.Fatalf("wrong updatedat want:10000, got:%v", tr.UpdatedAt)
	}

	// 更新行程的信息
	update := &rentalpb.Trip{
		AccountId: aid.String(),
		Status:    rentalpb.TripStatus_IN_PROGRESS,
		Start: &rentalpb.LocationStatus{
			PoiName: "start_poi_updated",
		},
	}
	cases := []struct {
		name          string
		now           int64
		withUpdatedAt int64
		wantErr       bool
	}{
		{
			name:          "normal_update",
			now:           20000,
			withUpdatedAt: 10000,
		},
		{
			name:          "normal_with_stale_timestamp",
			now:           30000,
			withUpdatedAt: 10000,
			wantErr:       true,
		},
		{
			name:          "update_with_refetch",
			now:           40000,
			withUpdatedAt: 20000,
		},
	}

	for _, cc := range cases {
		now = cc.now
		err := m.UpdateTrip(c, tid, aid, cc.withUpdatedAt, update)
		if cc.wantErr {
			if err == nil {
				t.Errorf("%s :want err ; got none", cc.name)
			} else {
				continue
			}
		} else {
			if err != nil {
				t.Errorf("%s :cannot update err", cc.name)
			}
		}
		// 获取更新后的行程
		updatedTrip, err := m.GetTrip(c, tid, aid)
		if err != nil {
			t.Errorf("%s cannot get trip after updated trip ,%v", cc.name, err)
		}
		if cc.now != updatedTrip.UpdatedAt {
			t.Errorf("%s: incorrect updatedat: want:%d ,got:%d", cc.name, cc.now, updatedTrip.UpdatedAt)
		}
	}
}

func TestMain(m *testing.M) {
	os.Exit(mongotesting.RunWithMongoInDocker(m))
}

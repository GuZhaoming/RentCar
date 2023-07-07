package trip

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	rentalpb "rentCar/server/rental/api/gen/v1"
	poi "rentCar/server/rental/trip/client/poi"
	"rentCar/server/rental/trip/dao"
	"rentCar/server/shared/auth"
	"rentCar/server/shared/id"
	mgutil "rentCar/server/shared/mongo"
	mongotesting "rentCar/server/shared/mongo/testing"
	"testing"

	"go.uber.org/zap"
)

func TestCreateTrip(t *testing.T) {
	c := auth.ContextWithAccountID(context.Background(), id.AccountID("account1"))
	mc, err := mongotesting.NewClient(c)
	if err != nil {
		t.Fatalf("cannot create mongo client: %v", err)
	}

	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatalf("cannot create mongo client : %v", err)
	}

	pm := &profileManager{}
	cm := &carManager{}
	s := &Service{
		ProfileManager: pm,
		CarManager:     cm,
		POIManager:     &poi.Manager{},
		Mongo:          dao.NewMongo(mc.Database("rentCar")),
		Logger:         logger,
	}

	req := &rentalpb.CreateTripRequest{
		CarId: "car1",
		Start: &rentalpb.Location{
			Latitude:  32.123,
			Longitude: 114.2525,
		},
	}
	pm.iID = "identity1"

	golden := `{"account_id":"account1","car_id":"car1","start":{"Location":{"latitude":32.123,"longitude":114.2525},"poi_name":"天安门"},"current":{"Location":{"latitude":32.123,"longitude":114.2525},"poi_name":"天安门"},"status":1,"identity_id":"identity1"}`

	cases := []struct {
		name         string
		tripID       string
		profileErr   error
		carVerifyErr error
		carUnlockErr error
		want         string
		wantErr      bool
	}{
		{
			name:   "normal_create",
			tripID: "648badb1b4600ccb093e5b29",
			want:   golden,
		},
		{
			name:       "profile_err",
			tripID:     "648badb1b4600ccb093e5b21",
			profileErr: fmt.Errorf("profile"),
			wantErr:    true,
		},
		{
			name:       "car_verify_err",
			tripID:     "648badb1b4600ccb093e5b22",
			profileErr: fmt.Errorf("verify"),
			wantErr:    true,
		},
		{
			name:       "car_unlock_err",
			tripID:     "648badb1b4600ccb093e5b28",
			want:       golden,
		},
	}

	for _, cc := range cases {
		t.Run(cc.name, func(t *testing.T) {
			mgutil.NewObjIDWithValue(id.TripID(cc.tripID))
			pm.err = cc.profileErr
			cm.unlockErr = cc.carUnlockErr
			cm.verifyErr = cc.carVerifyErr
			res, err := s.CreateTrip(c, req)
			if cc.wantErr {
				if err == nil {
					t.Errorf("want err ,got none")
				} else {
					return
				}
			}
			if err != nil {
				t.Errorf("error create trip: %v", err)
				return
			}

			if res.Id != cc.tripID {
				t.Errorf("incorrect id; want:%q,got:%q", cc.tripID, res.Id)
			}

			b, err := json.Marshal(res.Trip)
			if err != nil {
				t.Errorf("cannot marshal response: %v", err)
			}

			got := string(b)
			if cc.want != got {
				t.Errorf("incorrect response want:%s got:%s", cc.want, got)
			}
		})
	}
}

type profileManager struct {
	iID id.IdentityID
	err error
}

func (p *profileManager) Verify(context.Context, id.AccountID) (id.IdentityID, error) {
	return p.iID, p.err
}

type carManager struct {
	verifyErr error
	unlockErr error
}

func (c *carManager) Verify(context.Context, id.CarID, *rentalpb.Location) error {
	return c.verifyErr
}

func (c *carManager) Unlock(context.Context, id.CarID) error {
	return c.unlockErr
}

func TestMain(m *testing.M) {
	os.Exit(mongotesting.RunWithMongoInDocker(m))
}

package poi

import (
	"context"
	"hash/fnv"
	rentalpb "rentCar/server/rental/api/gen/v1"

	"google.golang.org/protobuf/proto"
)

var poi = []string{
	"中关村",
	"天安门",
	"陆家嘴",
	"迪士尼",
	"太南河体育中心",
	"广州塔",
}

type Manager struct {
}

func (*Manager) Resolve(c context.Context, loc *rentalpb.Location) (string, error) {
	b, err := proto.Marshal(loc)
	if err != nil {
		return "", err
	}

	h := fnv.New32()
	h.Write(b)
	return poi[int(h.Sum32())%len(poi)], nil
}

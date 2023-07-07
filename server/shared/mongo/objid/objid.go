package objid

import (
	"fmt"
	"rentCar/server/shared/id"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// convert id to objectid
func FromID(id fmt.Stringer) (primitive.ObjectID, error) {
	return primitive.ObjectIDFromHex(id.String())
}

func MustFromID(id fmt.Stringer) primitive.ObjectID {
	oid, err := FromID(id)
	if err != nil {
		panic(err)
	}
	return oid
}

// convert accountid
func ToAccountID(oid primitive.ObjectID) id.AccountID {
	return id.AccountID(oid.Hex())
}

func ToTripID(oid primitive.ObjectID) id.TripID {
	return id.TripID(oid.Hex())
}

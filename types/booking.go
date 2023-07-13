package types

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Booking struct {
	Id         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserId     primitive.ObjectID `bson:"userId,omitempty" json:"userId,omitempty"`
	RoomId     primitive.ObjectID `bson:"roomId,omitempty" json:"roomId,omitempty"`
	NumPersons int                `bson:"numPersons,omitempty" json:"numPersons,omitempty"`
	FromDate   time.Time          `bson:"fromDate,omitempty" json:"fromDate,omitempty"`
	TillDate   time.Time          `bson:"tillDate,omitempty" json:"tillDate,omitempty"`
	Canceled   bool               `bson:"canceled" json:"canceled"`
}

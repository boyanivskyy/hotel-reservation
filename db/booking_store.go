package db

import (
	"context"

	"github.com/boyanivskyy/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BookingStore interface {
	InsertBooking(ctx context.Context, booking *types.Booking) (*types.Booking, error)
	GetBookings(ctx context.Context, filter map[string]any) ([]*types.Booking, error)
	GetBooking(ctx context.Context, bookingId string) (*types.Booking, error)
	UpdateBooking(context.Context, string, map[string]any) error
}

type MongoBookingStore struct {
	client *mongo.Client
	coll   *mongo.Collection

	BookingStore
}

func NewMongoBookingStore(client *mongo.Client, dbName string) *MongoBookingStore {
	return &MongoBookingStore{
		client: client,
		coll:   client.Database(dbName).Collection("bookings"),
	}
}

func (s *MongoBookingStore) UpdateBooking(ctx context.Context, bookingId string, update map[string]any) error {
	bookingOid, err := primitive.ObjectIDFromHex(bookingId)
	if err != nil {
		return err
	}
	m := bson.M{
		"$set": update,
	}
	_, err = s.coll.UpdateByID(ctx, bookingOid, m)
	return err
}

func (s *MongoBookingStore) InsertBooking(ctx context.Context, booking *types.Booking) (*types.Booking, error) {
	resp, err := s.coll.InsertOne(ctx, booking)
	if err != nil {
		return nil, err
	}

	booking.Id = resp.InsertedID.(primitive.ObjectID)

	return booking, nil
}

func (s *MongoBookingStore) GetBookings(ctx context.Context, filter map[string]any) ([]*types.Booking, error) {
	curr, err := s.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	bookings := []*types.Booking{}
	if err := curr.All(ctx, &bookings); err != nil {
		return nil, err
	}
	return bookings, nil
}

func (s *MongoBookingStore) GetBooking(ctx context.Context, bookingId string) (*types.Booking, error) {
	bookingOid, err := primitive.ObjectIDFromHex(bookingId)
	if err != nil {
		return nil, err
	}

	booking := types.Booking{}
	if err := s.coll.FindOne(ctx, bson.M{"_id": bookingOid}).Decode(&booking); err != nil {
		return nil, err
	}

	return &booking, nil
}

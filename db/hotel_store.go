package db

import (
	"context"

	"github.com/boyanivskyy/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type HotelStore interface {
	Insert(ctx context.Context, hotel *types.Hotel) (*types.Hotel, error)
	Update(ctx context.Context, filter map[string]any, values map[string]any) error
	GetHotels(ctx context.Context, filter map[string]any) ([]*types.Hotel, error)
	GetHotelById(ctx context.Context, id string) (*types.Hotel, error)
}

type MongoHotelStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoHotelStore(client *mongo.Client, dbName string) *MongoHotelStore {
	return &MongoHotelStore{
		client: client,
		coll:   client.Database(dbName).Collection("hotels"),
	}
}

func (s *MongoHotelStore) Insert(ctx context.Context, hotel *types.Hotel) (*types.Hotel, error) {
	res, err := s.coll.InsertOne(ctx, hotel)
	if err != nil {
		return nil, err
	}
	hotel.Id = res.InsertedID.(primitive.ObjectID)
	return hotel, nil
}

func (s *MongoHotelStore) Update(ctx context.Context, filter map[string]any, values map[string]any) error {
	_, err := s.coll.UpdateOne(ctx, filter, values)
	if err != nil {
		return err
	}

	return nil
}

func (s *MongoHotelStore) GetHotels(ctx context.Context, filter map[string]any) ([]*types.Hotel, error) {
	resp, err := s.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	hotels := []*types.Hotel{}
	if err := resp.All(ctx, &hotels); err != nil {
		return nil, err
	}

	return hotels, err
}

func (s *MongoHotelStore) GetHotelById(ctx context.Context, id string) (*types.Hotel, error) {
	oId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	hotel := types.Hotel{}
	if err := s.coll.FindOne(ctx, bson.M{
		"_id": oId,
	}).Decode(&hotel); err != nil {
		return nil, err
	}

	return &hotel, nil
}

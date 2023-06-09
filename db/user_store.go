package db

import (
	"context"

	"github.com/boyanivskyy/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	userColl = "users"
)

type UserStore interface {
	GetUserById(ctx context.Context, id string) (*types.User, error)
	GetUsers(ctx context.Context) ([]*types.User, error)
	InsertUser(ctx context.Context, user *types.User) (*types.User, error)
	DeleteUser(ctx context.Context, id string) error
	UpdateUser(ctx context.Context, filter bson.M, values types.UpdateUserParams) error
}

type MongoUserStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func (store *MongoUserStore) UpdateUser(ctx context.Context, filter bson.M, params types.UpdateUserParams) error {
	values := params.ToBSON()
	update := bson.D{{"$set", values}}
	_, err := store.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (store *MongoUserStore) DeleteUser(ctx context.Context, id string) error {
	oId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = store.coll.DeleteOne(ctx, bson.M{
		"_id": oId,
	})
	if err != nil {
		return err
	}

	return nil
}

func (store *MongoUserStore) GetUserById(ctx context.Context, id string) (*types.User, error) {
	oId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	user := types.User{}
	if err := store.coll.FindOne(ctx, bson.M{"_id": oId}).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (store *MongoUserStore) GetUsers(ctx context.Context) ([]*types.User, error) {
	cur, err := store.coll.Find(ctx, bson.M{})
	if err != nil {
		return []*types.User{}, err
	}

	users := []*types.User{}
	if err := cur.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}

func (store *MongoUserStore) InsertUser(ctx context.Context, user *types.User) (*types.User, error) {
	res, err := store.coll.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	user.Id = res.InsertedID.(primitive.ObjectID)

	return user, nil
}

func NewMongoUserStore(client *mongo.Client) *MongoUserStore {
	return &MongoUserStore{
		client: client,
		coll:   client.Database(DBNAME).Collection(userColl),
	}
}

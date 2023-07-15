package db

import (
	"context"
	"fmt"

	"github.com/boyanivskyy/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	userColl = "users"
)

type Dropper interface {
	Drop(ctx context.Context) error
}

type UserStore interface {
	Dropper

	GetUserById(ctx context.Context, id string) (*types.User, error)
	GetUsers(ctx context.Context) ([]*types.User, error)
	InsertUser(ctx context.Context, user *types.User) (*types.User, error)
	DeleteUser(ctx context.Context, id string) error
	UpdateUser(ctx context.Context, filter map[string]any, values types.UpdateUserParams) error
	GetUserByEmail(ctx context.Context, email string) (*types.User, error)
}

type MongoUserStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func (store *MongoUserStore) Drop(ctx context.Context) error {
	fmt.Println("dropping user collection")
	return store.coll.Drop(ctx)
}

func (store *MongoUserStore) UpdateUser(ctx context.Context, filter map[string]any, params types.UpdateUserParams) error {
	oId, err := primitive.ObjectIDFromHex(filter["_id"].(string))
	if err != nil {
		return err
	}
	filter["_id"] = oId

	m := bson.M{
		"$set": params.ToBSON(),
	}
	_, err = store.coll.UpdateOne(ctx, filter, m)
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
	// TODO: check if exist by email???
	res, err := store.coll.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	user.Id = res.InsertedID.(primitive.ObjectID)

	return user, nil
}

func (store *MongoUserStore) GetUserByEmail(ctx context.Context, email string) (*types.User, error) {
	user := types.User{}
	if err := store.coll.FindOne(ctx, bson.M{"email": email}).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func NewMongoUserStore(client *mongo.Client, dbName string) *MongoUserStore {
	return &MongoUserStore{
		client: client,
		coll:   client.Database(dbName).Collection(userColl),
	}
}

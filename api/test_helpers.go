package api

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/boyanivskyy/hotel-reservation/db"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	testMongoUri = "mongodb://localhost:27017"
)

type testdb struct {
	client *mongo.Client
	*db.Store
}

func (tdb testdb) tearDown(t *testing.T) {
	dbname := os.Getenv("MONGO_TEST_DBNAME")
	if err := tdb.client.Database(dbname).Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}

func setup(t *testing.T) *testdb {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatal(err)
	}

	dbname := os.Getenv("MONGO_TEST_DBNAME")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(testMongoUri))
	if err != nil {
		return nil
	}

	hotelStore := db.NewMongoHotelStore(client, dbname)
	return &testdb{
		client: client,
		Store: &db.Store{
			Hotel:   hotelStore,
			User:    db.NewMongoUserStore(client, dbname),
			Room:    db.NewMongoRoomStore(client, hotelStore, dbname),
			Booking: db.NewMongoBookingStore(client, dbname),
		},
	}
}

package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/boyanivskyy/hotel-reservation/api"
	"github.com/boyanivskyy/hotel-reservation/db"
	"github.com/boyanivskyy/hotel-reservation/db/fixtures"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx := context.Background()
	mongoEndpoint := os.Getenv("MONGO_URI")
	mongoDBName := os.Getenv("MONGO_DBNAME")
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoEndpoint))
	if err != nil {
		log.Fatal(err)
	}

	if err := client.Database(mongoDBName).Drop(ctx); err != nil {
		log.Fatal(err)
	}

	hotelStore := db.NewMongoHotelStore(client, mongoDBName)
	store := &db.Store{
		User:    db.NewMongoUserStore(client, mongoDBName),
		Hotel:   hotelStore,
		Room:    db.NewMongoRoomStore(client, hotelStore, mongoDBName),
		Booking: db.NewMongoBookingStore(client, mongoDBName),
	}

	user := fixtures.AddUser(store, "test", "test", false)
	fmt.Println("user -> ", api.CreateTokenFromUser(user))
	admin := fixtures.AddUser(store, "admin", "admin", true)
	fmt.Println("admin -> ", api.CreateTokenFromUser(admin))
	hotel := fixtures.AddHotel(store, "Test Hotel", "anywhere", 5, []primitive.ObjectID{})
	room := fixtures.AddRoom(store, "small", true, 20.99, hotel.Id)
	booking := fixtures.AddBooking(store, user.Id, room.Id, 2, time.Now(), time.Now().AddDate(0, 0, 2))
	fmt.Println("booking ->", booking.Id)

	for i := 0; i < 100; i++ {
		fixtures.AddHotel(store, fmt.Sprintf("hotel %d", i), fmt.Sprintf("loc %d", i), rand.Intn(5)+1, nil)
	}
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
}

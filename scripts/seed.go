package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/boyanivskyy/hotel-reservation/api"
	"github.com/boyanivskyy/hotel-reservation/db"
	"github.com/boyanivskyy/hotel-reservation/db/fixtures"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	if err := client.Database(db.DBNAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}

	hotelStore := db.NewMongoHotelStore(client, db.DBNAME)
	store := &db.Store{
		User:    db.NewMongoUserStore(client, db.DBNAME),
		Hotel:   hotelStore,
		Room:    db.NewMongoRoomStore(client, hotelStore, db.DBNAME),
		Booking: db.NewMongoBookingStore(client, db.DBNAME),
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

package main

import (
	"context"
	"fmt"
	"log"

	"github.com/boyanivskyy/hotel-reservation/db"
	"github.com/boyanivskyy/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client     *mongo.Client
	roomStore  db.RoomStore
	hotelStore db.HotelStore
	userStore  db.UserStore
	ctx        = context.Background()
)

func seedUser(fname, lname, email string) {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		Email:     email,
		FirstName: fname,
		LastName:  lname,
		Password:  "qweasdzxc",
	})

	if err != nil {
		log.Fatal(err)
	}

	insertedUser, err := userStore.InsertUser(context.TODO(), user)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("insertedUser", insertedUser)
}

func seedHotel(name, location string, rating int) {
	hotel := types.Hotel{
		Name:     name,
		Location: location,
		Rooms:    []primitive.ObjectID{},
		Rating:   rating,
	}

	rooms := []types.Room{
		{
			Size:  "small",
			Price: 99.99,
		},
		{
			Size:  "normal",
			Price: 159.99,
		},
		{
			Size:  "kingsize",
			Price: 1099.99,
		},
	}

	insertedHotel, err := hotelStore.Insert(ctx, &hotel)
	if err != nil {
		log.Fatal(err)
	}

	for idx, room := range rooms {
		room.HotelId = insertedHotel.Id
		insertedRoom, err := roomStore.InsertRoom(ctx, &room)
		rooms[idx].Id = insertedRoom.Id
		rooms[idx].HotelId = insertedHotel.Id
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("insertedHotel", insertedHotel)
	fmt.Println("insertedRooms", rooms)
	fmt.Println()
}

func main() {
	seedHotel("Redisson Blue", "Bukovel", 3)
	seedHotel("HAY Hotel", "Bukovel", 5)
	seedHotel("Emem Resort", "Lviv", 4)
	seedUser("foo", "bar", "foo@bar.com")
}

func init() {
	var err error
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	if err := client.Database(db.DBNAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}

	hotelStore = db.NewMongoHotelStore(client)
	roomStore = db.NewMongoRoomStore(client, hotelStore)
	userStore = db.NewMongoUserStore(client)
}

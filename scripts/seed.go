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
	ctx        = context.Background()
)

func seedHotel(name, location string) {
	hotel := types.Hotel{
		Name:     name,
		Location: location,
		Rooms:    []primitive.ObjectID{},
	}

	rooms := []types.Room{
		{
			Type:      types.SingleRoomType,
			BasePrice: 99.99,
		},
		{
			Type:      types.DoubleRoomType,
			BasePrice: 159.99,
		},
		{
			Type:      types.DeluxeRoomType,
			BasePrice: 1099.99,
		},
	}

	insertedHotel, err := hotelStore.InsertHotel(ctx, &hotel)
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
	seedHotel("Redisson Blue", "Bukovel")
	seedHotel("HAY Hotel", "Bukovel")
	seedHotel("Emem Resort", "Lviv")
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
}

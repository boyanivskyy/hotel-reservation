package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/boyanivskyy/hotel-reservation/api"
	"github.com/boyanivskyy/hotel-reservation/db"
	"github.com/boyanivskyy/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client       *mongo.Client
	roomStore    db.RoomStore
	hotelStore   db.HotelStore
	userStore    db.UserStore
	bookingStore db.BookingStore
	ctx          = context.Background()
)

func seedUser(isAdmin bool, params types.CreateUserParams) *types.User {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		Email:     params.Email,
		FirstName: params.FirstName,
		LastName:  params.LastName,
		Password:  params.Password,
	})

	if err != nil {
		log.Fatal(err)
	}

	user.IsAdmin = isAdmin

	insertedUser, err := userStore.InsertUser(context.TODO(), user)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s - %s\n", insertedUser.Email, api.CreateTokenFromUser(user))

	return user
}

func seedHotel(name, location string, rating int) *types.Hotel {
	hotel := types.Hotel{
		Name:     name,
		Location: location,
		Rooms:    []primitive.ObjectID{},
		Rating:   rating,
	}

	insertedHotel, err := hotelStore.Insert(ctx, &hotel)
	if err != nil {
		log.Fatal(err)
	}
	return insertedHotel
}

func seedRoom(size string, seaside bool, price float64, hotelId primitive.ObjectID) *types.Room {
	room := &types.Room{
		Size:    size,
		Seaside: seaside,
		Price:   price,
		HotelId: hotelId,
	}

	insertedRoom, err := roomStore.InsertRoom(context.Background(), room)
	if err != nil {
		log.Fatal(err)
	}

	return insertedRoom
}

func seedBooking(userId, roomId primitive.ObjectID, from, till time.Time, numPersons int) {
	booking := types.Booking{
		UserId:     userId,
		RoomId:     roomId,
		FromDate:   from,
		TillDate:   till,
		NumPersons: 2,
	}

	insertedBooking, err := bookingStore.InsertBooking(context.Background(), &booking)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("booking id -> ", insertedBooking.Id)
}

func main() {
	user := seedUser(false, types.CreateUserParams{
		FirstName: "foo",
		LastName:  "bar",
		Email:     "foo@bar.com",
		Password:  "qweasdzxc",
	})
	seedUser(true, types.CreateUserParams{
		FirstName: "admin",
		LastName:  "admin",
		Email:     "admin@admin.com",
		Password:  "coolpass123word",
	})
	fmt.Println("_________________________________")

	seedHotel("Redisson Blue", "Bukovel", 3)
	seedHotel("HAY Hotel", "Bukovel", 5)
	hotel := seedHotel("Emem Resort", "Lviv", 4)

	seedRoom("small", true, 99.99, hotel.Id)
	seedRoom("medium", true, 199.99, hotel.Id)
	room := seedRoom("large", true, 399.99, hotel.Id)

	seedBooking(user.Id, room.Id, time.Now(), time.Now().Add(time.Hour*48), 2)
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
	bookingStore = db.NewMongoBookingStore(client)
}

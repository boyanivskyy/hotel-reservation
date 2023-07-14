package fixtures

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/boyanivskyy/hotel-reservation/db"
	"github.com/boyanivskyy/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddUser(store *db.Store, firstName, lastName string, isAdmin bool) *types.User {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		Email:     fmt.Sprintf("%s@%s.com", firstName, lastName),
		FirstName: firstName,
		LastName:  lastName,
		Password:  fmt.Sprintf("%s_%s", firstName, lastName),
	})
	if err != nil {
		log.Fatal(err)
	}

	user.IsAdmin = isAdmin
	insertedUser, err := store.User.InsertUser(context.Background(), user)
	if err != nil {
		log.Fatal(err)
	}
	return insertedUser
}

func AddBooking(store *db.Store, userId, roomId primitive.ObjectID, numbPersons int, from, till time.Time) *types.Booking {
	booking := types.Booking{
		UserId:     userId,
		RoomId:     roomId,
		FromDate:   from,
		TillDate:   till,
		NumPersons: 2,
	}

	insertedBooking, err := store.Booking.InsertBooking(context.Background(), &booking)
	if err != nil {
		log.Fatal(err)
	}

	return insertedBooking
}

func AddRoom(store *db.Store, size string, seaside bool, price float64, hotelId primitive.ObjectID) *types.Room {
	room := &types.Room{
		Size:    size,
		Seaside: seaside,
		Price:   price,
		HotelId: hotelId,
	}

	insertedRoom, err := store.Room.InsertRoom(context.Background(), room)
	if err != nil {
		log.Fatal(err)
	}

	return insertedRoom
}

func AddHotel(store *db.Store, name, location string, rating int, rooms []primitive.ObjectID) *types.Hotel {
	roomIds := rooms
	if rooms == nil {
		roomIds = []primitive.ObjectID{}
	}
	hotel := types.Hotel{
		Name:     name,
		Location: location,
		Rooms:    roomIds,
		Rating:   rating,
	}

	insertedHotel, err := store.Hotel.Insert(context.Background(), &hotel)
	if err != nil {
		log.Fatal(err)
	}
	return insertedHotel
}

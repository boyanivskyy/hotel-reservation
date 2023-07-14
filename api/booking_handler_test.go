package api

import (
	"fmt"
	"testing"
	"time"

	"github.com/boyanivskyy/hotel-reservation/db"
	"github.com/boyanivskyy/hotel-reservation/db/fixtures"
)

func TestGetBookings(t *testing.T) {
	db := setup(t, db.TestDBNAME)
	defer db.tearDown(t)

	user := fixtures.AddUser(db.Store, "test", "test", false)
	hotel := fixtures.AddHotel(db.Store, "bar hotel", "a", 4, nil)
	room := fixtures.AddRoom(db.Store, "small", true, 99.99, hotel.Id)
	booking := fixtures.AddBooking(db.Store, user.Id, room.Id, 2, time.Now(), time.Now().AddDate(0, 0, 2))
	fmt.Println(booking)
}

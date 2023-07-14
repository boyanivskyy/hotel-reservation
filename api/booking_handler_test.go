package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/boyanivskyy/hotel-reservation/db"
	"github.com/boyanivskyy/hotel-reservation/db/fixtures"
	"github.com/boyanivskyy/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
)

func TestUserGetBooking(t *testing.T) {
	db := setup(t, db.TestDBNAME)
	defer db.tearDown(t)

	var (
		nonAuthUser = fixtures.AddUser(db.Store, "nonauth", "nonauth", false)
		user        = fixtures.AddUser(db.Store, "test", "test", false)
		hotel       = fixtures.AddHotel(db.Store, "bar hotel", "a", 4, nil)
		room        = fixtures.AddRoom(db.Store, "small", true, 99.99, hotel.Id)
		booking     = fixtures.AddBooking(db.Store, user.Id, room.Id, 2, time.Now(), time.Now().AddDate(0, 0, 2))
		app         = fiber.New(fiber.Config{
			ErrorHandler: ErrorHandler,
		})
		api            = app.Group("/", JWTAuthentication(db.User))
		bookingHandler = NewBookingHandler(db.Store)
	)

	api.Get("/:id", bookingHandler.HandleGetBooking)
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/%s", booking.Id.Hex()), nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(user))
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("not 200 response, got %d", resp.StatusCode)
	}

	have := types.Booking{}
	if err := json.NewDecoder(resp.Body).Decode(&have); err != nil {
		t.Fatal(err)
	}
	if have.Id != booking.Id {
		t.Fatalf("expected bookingId %s, got %s", booking.Id, have.Id)
	}
	if have.UserId != booking.UserId {
		t.Fatalf("expected userId %s, got %s", booking.UserId, have.UserId)
	}

	// non auth user
	req.Header.Add("X-Api-Token", CreateTokenFromUser(nonAuthUser))

	resp, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode == http.StatusOK {
		t.Fatalf("expected a non 200 statusCode, got %d", resp.StatusCode)
	}
}

func TestAdminGetBookings(t *testing.T) {
	db := setup(t, db.TestDBNAME)
	defer db.tearDown(t)

	var (
		user      = fixtures.AddUser(db.Store, "test", "test", false)
		_         = user
		adminUser = fixtures.AddUser(db.Store, "admin", "admin", true)
		hotel     = fixtures.AddHotel(db.Store, "bar hotel", "a", 4, nil)
		room      = fixtures.AddRoom(db.Store, "small", true, 99.99, hotel.Id)
		booking   = fixtures.AddBooking(db.Store, adminUser.Id, room.Id, 2, time.Now(), time.Now().AddDate(0, 0, 2))
		app       = fiber.New(fiber.Config{
			ErrorHandler: ErrorHandler,
		})
		adminGroup     = app.Group("/", JWTAuthentication(db.User), AdminAuth)
		bookingHandler = NewBookingHandler(db.Store)
	)

	adminGroup.Get("/", bookingHandler.HandleGetBookings)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(adminUser))
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("not 200 response, got %d", resp.StatusCode)
	}

	bookings := []*types.Booking{}
	if err := json.NewDecoder(resp.Body).Decode(&bookings); err != nil {
		t.Fatal(err)
	}
	if len(bookings) < 1 {
		t.Fatalf("expected 1 booking, got %d", len(bookings))
	}
	have := bookings[0]
	if have.Id != booking.Id {
		t.Fatalf("expected bookingId %s, got %s", booking.Id, have.Id)
	}
	if have.UserId != booking.UserId {
		t.Fatalf("expected userId %s, got %s", booking.UserId, have.UserId)
	}

	// NOTE: test when user is not admin and can not access the endpoint
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(user))
	resp, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expeceted status 401 unathorized, got %d", resp.StatusCode)
	}
}

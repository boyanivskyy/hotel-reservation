package api

import (
	"fmt"
	"net/http"

	"github.com/boyanivskyy/hotel-reservation/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

type BookingHandler struct {
	store *db.Store
}

func NewBookingHandler(store *db.Store) *BookingHandler {
	return &BookingHandler{
		store: store,
	}
}

func (h *BookingHandler) HandleCancelBooking(c *fiber.Ctx) error {
	id := c.Params("id")
	fmt.Println("id", id)
	booking, err := h.store.Booking.GetBooking(c.Context(), id)
	if err != nil {
		return err
	}

	user, err := GetAuthUser(c)
	if err != nil {
		return err
	}

	if booking.UserId != user.Id {
		return c.Status(http.StatusUnauthorized).JSON(genericResp{
			Type: "error",
			Msg:  "not authorized",
		})
	}

	if booking.Canceled {
		return c.Status(http.StatusBadRequest).JSON(genericResp{
			Type: "error",
			Msg:  "booking already canceled",
		})
	}

	if err := h.store.Booking.UpdateBooking(c.Context(), id, bson.M{"canceled": true}); err != nil {
		return err
	}

	return c.JSON(genericResp{
		Type: "message",
		Msg:  "updated",
	})
}

func (h *BookingHandler) HandleGetBookings(c *fiber.Ctx) error {
	bookings, err := h.store.Booking.GetBookings(c.Context(), bson.M{})
	if err != nil {
		return err
	}
	return c.JSON(bookings)
}

func (h *BookingHandler) HandleGetBooking(c *fiber.Ctx) error {
	booking, err := h.store.Booking.GetBooking(c.Context(), c.Params("id"))
	if err != nil {
		return err
	}
	user, err := GetAuthUser(c)
	if err != nil {
		return err
	}

	if booking.UserId != user.Id {
		return c.Status(http.StatusUnauthorized).JSON(genericResp{
			Type: "error",
			Msg:  "not authorized",
		})
	}

	return c.JSON(booking)
}

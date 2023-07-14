package api

import (
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
	booking, err := h.store.Booking.GetBooking(c.Context(), id)
	if err != nil {
		return ErrorResourceNotFound()
	}

	user, err := GetAuthUser(c)
	if err != nil {
		return ErrorUnathorized()
	}

	if booking.UserId != user.Id {
		return ErrorUnathorized()
	}

	if booking.Canceled {
		return ErrorBadRequest()
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
		return ErrorResourceNotFound()
	}
	return c.JSON(bookings)
}

func (h *BookingHandler) HandleGetBooking(c *fiber.Ctx) error {
	booking, err := h.store.Booking.GetBooking(c.Context(), c.Params("id"))
	if err != nil {
		return ErrorResourceNotFound()
	}
	user, err := GetAuthUser(c)
	if err != nil {
		return ErrorUnathorized()
	}

	if booking.UserId != user.Id {
		return ErrorUnathorized()
	}

	return c.JSON(booking)
}

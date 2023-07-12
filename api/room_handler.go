package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/boyanivskyy/hotel-reservation/db"
	"github.com/boyanivskyy/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BookRoomParams struct {
	NumPersons int       `json:"numPersons"`
	FromDate   time.Time `json:"fromDate"`
	TillDate   time.Time `json:"tillDate"`
}

func (p BookRoomParams) validate() error {
	today := time.Now()
	yesterday := today.Add(time.Hour * 24 * -1)

	if yesterday.After(p.FromDate) || today.After(p.TillDate) {
		return fmt.Errorf("can not book a room in the past")
	}
	if p.FromDate.After(p.TillDate) {
		return fmt.Errorf("can not book a room, tillDate must be greater than fromDate")
	}
	return nil
}

type RoomHandler struct {
	store *db.Store
}

func NewRoomHandler(store *db.Store) *RoomHandler {
	return &RoomHandler{
		store: store,
	}
}

func (h *RoomHandler) HandleGetRooms(c *fiber.Ctx) error {
	rooms, err := h.store.Room.GetRooms(c.Context(), bson.M{})
	if err != nil {
		return err
	}
	return c.JSON(rooms)
}

func (h *RoomHandler) HandleBookRoom(c *fiber.Ctx) error {
	params := BookRoomParams{}
	if err := c.BodyParser(&params); err != nil {
		return err
	}
	if err := params.validate(); err != nil {
		return err
	}
	roomOId, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return err
	}
	user, ok := c.Context().Value("user").(*types.User)
	if !ok {
		return c.Status(http.StatusInternalServerError).JSON(genericResp{
			Type: "error",
			Msg:  "internal server error",
		})
	}
	ok, err = h.isRoomAvailableForBooking(c.Context(), roomOId, params)
	if err != nil {
		return err
	}
	if !ok {
		return c.Status(http.StatusBadRequest).JSON(genericResp{
			Type: "error",
			Msg:  fmt.Sprintf("room %s is already booked", c.Params("id")),
		})
	}
	booking := types.Booking{
		RoomId:     roomOId,
		UserId:     user.Id,
		NumPersons: params.NumPersons,
		FromDate:   params.FromDate,
		TillDate:   params.TillDate,
	}
	inserted, err := h.store.Booking.InsertBooking(c.Context(), &booking)
	if err != nil {
		return err
	}
	return c.JSON(inserted)
}

func (h *RoomHandler) isRoomAvailableForBooking(ctx context.Context, roomOId primitive.ObjectID, params BookRoomParams) (bool, error) {
	where := bson.M{
		"roomId": roomOId,
		"fromDate": bson.M{
			"$gte": params.FromDate,
		},
		"tillDate": bson.M{
			"$lte": params.TillDate,
		},
	}
	bookings, err := h.store.Booking.GetBookings(ctx, where)
	if err != nil {
		return false, err
	}

	ok := len(bookings) == 0

	return ok, nil
}

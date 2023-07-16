package api

import (
	"github.com/boyanivskyy/hotel-reservation/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HotelHandler struct {
	store *db.Store
}

func NewHotelHandler(store *db.Store) *HotelHandler {
	return &HotelHandler{
		store: store,
	}
}

func (h *HotelHandler) HandleGetRooms(c *fiber.Ctx) error {
	id := c.Params("id")
	oId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ErrorInvalidId()
	}

	filter := bson.M{
		"hotelId": oId,
	}

	rooms, err := h.store.Room.GetRooms(c.Context(), filter)
	if err != nil {
		return ErrorResourceNotFound()
	}

	return c.JSON(rooms)
}

type ApiSliceResponse struct {
	Data  any   `json:"data"`
	Total int   `json:"total,omitempty"`
	Page  int64 `json:"page,omitempty"`
}

type HotelQueryParams struct {
	db.PaginationFilter

	Rating int `json:"rating"`
}

func (h *HotelHandler) HandleGetHotels(c *fiber.Ctx) error {
	params := HotelQueryParams{}
	if err := c.QueryParser(&params); err != nil {
		return ErrorBadRequest()
	}

	filter := map[string]any{}

	if params.Rating != 0 {
		filter = map[string]any{
			"rating": params.Rating,
		}
	}

	hotels, err := h.store.Hotel.GetHotels(c.Context(), filter, params.PaginationFilter)
	if err != nil {
		return ErrorResourceNotFound()
	}

	resp := ApiSliceResponse{
		Data:  hotels,
		Total: len(hotels),
		Page:  params.Page,
	}
	return c.JSON(resp)
}

func (h *HotelHandler) HandleGetHotel(c *fiber.Ctx) error {
	id := c.Params("id")
	hotel, err := h.store.Hotel.GetHotelById(c.Context(), id)
	if err != nil {
		return ErrorResourceNotFound()
	}

	return c.JSON(hotel)
}

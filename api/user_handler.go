package api

import (
	"errors"
	"net/http"

	"github.com/boyanivskyy/hotel-reservation/db"
	"github.com/boyanivskyy/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserHandler struct {
	userStore db.UserStore
}

func NewUserHandler(userStore db.UserStore) *UserHandler {
	return &UserHandler{
		userStore: userStore,
	}
}

func (h *UserHandler) HandlePutUser(c *fiber.Ctx) error {
	userId := c.Params("id")

	params := types.UpdateUserParams{}
	if err := c.BodyParser(&params); err != nil {
		return ErrorBadRequest()
	}

	filter := bson.M{"_id": userId}
	if err := h.userStore.UpdateUser(c.Context(), filter, params); err != nil {
		return ErrorResourceNotFound()
	}

	return c.JSON(map[string]string{
		"id": userId,
	})
}

func (h *UserHandler) HandleDeleteUser(c *fiber.Ctx) error {
	userId := c.Params("id")

	err := h.userStore.DeleteUser(c.Context(), userId)
	if err != nil {
		return ErrorInvalidId()
	}

	return c.JSON(map[string]string{
		"id": userId,
	})
}

func (h *UserHandler) HandlePostUser(c *fiber.Ctx) error {
	params := types.CreateUserParams{}
	if err := c.BodyParser(&params); err != nil {
		return ErrorBadRequest()
	}

	if errors := params.Validate(); len(errors) > 0 {
		return c.Status(http.StatusBadRequest).JSON(errors)
	}

	user, err := types.NewUserFromParams(params)
	if err != nil {
		return ErrorBadRequest()
	}
	user, err = h.userStore.InsertUser(c.Context(), user)
	if err != nil {
		return ErrorBadRequest()
	}

	return c.JSON(user)
}

func (h *UserHandler) HandleGetUserById(c *fiber.Ctx) error {
	id := c.Params("id")

	user, err := h.userStore.GetUserById(c.Context(), id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrorResourceNotFound()
		}
		return ErrorResourceNotFound()
	}
	return c.JSON(user)
}

func (h *UserHandler) HandleGetUsers(c *fiber.Ctx) error {
	users, err := h.userStore.GetUsers(c.Context())
	if err != nil {
		return ErrorResourceNotFound()
	}

	return c.JSON(users)
}

package api

import (
	"github.com/boyanivskyy/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
)

func HandleGetUser(c *fiber.Ctx) error {
	user := types.User{
		FirstName: "Vitaliy",
		LastName:  "Boyanivskyy",
	}
	return c.JSON(user)
}

func HandleGetUserById(c *fiber.Ctx) error {
	userId := c.Params("id")
	return c.JSON(map[string]string{
		"id":   userId,
		"name": "Vitaliy",
	})
}

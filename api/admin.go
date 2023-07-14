package api

import (
	"github.com/boyanivskyy/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
)

func AdminAuth(c *fiber.Ctx) error {
	user, ok := c.Context().UserValue("user").(*types.User)
	if !ok {
		return ErrorUnathorized()
	}

	if !user.IsAdmin {
		return ErrorUnathorized()
	}

	return c.Next()
}

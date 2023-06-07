package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	app.Get("/foo", handleFoo)
	log.Fatal(app.Listen(":8080"))
}

func handleFoo(ctx *fiber.Ctx) error {
	return ctx.JSON(map[string]string{
		"msg": "working just fine",
	})
}

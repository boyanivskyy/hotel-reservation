package main

import (
	"context"
	"flag"
	"log"

	"github.com/boyanivskyy/hotel-reservation/api"
	"github.com/boyanivskyy/hotel-reservation/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const dbUri = "mongodb://localhost:27017"

var fiberConfig = fiber.Config{
	ErrorHandler: func(ctx *fiber.Ctx, err error) error {
		return ctx.JSON(map[string]string{
			"error": err.Error(),
		})
	},
}

func main() {
	// flag means when you run go run main.go you can add --listenAddr
	listenAddr := flag.String("listenAddr", ":8080", "The listen address of the API server")
	flag.Parse()

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dbUri))
	if err != nil {
		log.Fatal(err)
	}
	// handlers init
	userHandler := api.NewUserHandler(db.NewMongoUserStore(client))

	app := fiber.New(fiberConfig)

	apiv1 := app.Group("/api/v1")
	apiv1.Post("/user", userHandler.HandlePostUser)
	apiv1.Put("/user/:id", userHandler.HandlePutUser)
	apiv1.Get("/user/:id", userHandler.HandleGetUserById)
	apiv1.Get("/user", userHandler.HandleGetUsers)
	apiv1.Delete("/user/:id", userHandler.HandleDeleteUser)

	log.Fatal(app.Listen(*listenAddr))
}

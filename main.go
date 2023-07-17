package main

import (
	"context"
	"log"
	"os"

	"github.com/boyanivskyy/hotel-reservation/api"
	"github.com/boyanivskyy/hotel-reservation/db"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Configuration
// 1. MongoDB endpoint
// 2. ListenAddress of our http server
// 3. JWT secret
// 4. mongoDBname = hotel-reservation

var fiberConfig = fiber.Config{
	ErrorHandler: api.ErrorHandler,
}

func main() {
	mongoEndpoint := os.Getenv("MONGO_URI")
	mongoDBName := os.Getenv("MONGO_DBNAME")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoEndpoint))
	if err != nil {
		log.Fatal(err)
	}
	// stores list
	var (
		userStore    = db.NewMongoUserStore(client, mongoDBName)
		hotelStore   = db.NewMongoHotelStore(client, mongoDBName)
		roomStore    = db.NewMongoRoomStore(client, hotelStore, mongoDBName)
		bookingStore = db.NewMongoBookingStore(client, mongoDBName)
		store        = &db.Store{
			User:    userStore,
			Hotel:   hotelStore,
			Room:    roomStore,
			Booking: bookingStore,
		}
	)
	// handlers init
	var (
		authHandler    = api.NewAuthHandler(userStore)
		userHandler    = api.NewUserHandler(userStore)
		hotelHandler   = api.NewHotelHandler(store)
		roomHandler    = api.NewRoomHandler(store)
		bookingHandler = api.NewBookingHandler(store)
	)

	app := fiber.New(fiberConfig)
	auth := app.Group("/api")
	apiv1 := app.Group("/api/v1", api.JWTAuthentication(userStore))
	admin := apiv1.Group("/admin", api.AdminAuth)

	// auth handlers
	auth.Post("/auth", authHandler.HandleAuthenticate)

	// user handlers
	apiv1.Post("/user", userHandler.HandlePostUser)
	apiv1.Put("/user/:id", userHandler.HandlePutUser)
	apiv1.Get("/user/:id", userHandler.HandleGetUserById)
	apiv1.Get("/user", userHandler.HandleGetUsers)
	apiv1.Delete("/user/:id", userHandler.HandleDeleteUser)

	// hotel handlers
	apiv1.Get("/hotel", hotelHandler.HandleGetHotels)
	apiv1.Get("/hotel/:id", hotelHandler.HandleGetHotel)
	apiv1.Get("/hotel/:id/rooms", hotelHandler.HandleGetRooms)

	// room handler
	apiv1.Get("/room", roomHandler.HandleGetRooms)
	apiv1.Post("/room/:id/book", roomHandler.HandleBookRoom)
	// TODO: cancel a booking

	// booking handler
	apiv1.Get("/booking/:id", bookingHandler.HandleGetBooking)
	apiv1.Get("/booking/:id/cancel", bookingHandler.HandleCancelBooking)

	// admin handlers
	admin.Get("/booking", bookingHandler.HandleGetBookings)

	listenAddr := os.Getenv("HTTP_LISTEN_ADDRESS")
	log.Fatal(app.Listen(listenAddr))
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("unable to load .env: %s", err)
	}
}

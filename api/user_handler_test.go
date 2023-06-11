package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/boyanivskyy/hotel-reservation/db"
	"github.com/boyanivskyy/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	testMongoUri = "mongodb://localhost:27017"
	dbname       = "hotel-reservation-test"
)

type testdb struct {
	db.UserStore
}

func (tdb testdb) tearDown(t *testing.T) {
	if err := tdb.UserStore.Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}

func setup(t *testing.T) *testdb {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(testMongoUri))
	if err != nil {
		return nil
	}

	return &testdb{
		UserStore: db.NewMongoUserStore(client, dbname),
	}
}

func Test_HandlePostUser(t *testing.T) {
	tdb := setup(t)
	defer tdb.tearDown(t)

	app := fiber.New()
	userHandler := NewUserHandler(tdb.UserStore)
	app.Post("/", userHandler.HandlePostUser)

	params := types.CreateUserParams{
		Email:     "some@foo.com",
		FirstName: "test",
		LastName:  "test1",
		Password:  "random_password123",
	}

	b, _ := json.Marshal(params)
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}
	user := types.User{}
	json.NewDecoder(resp.Body).Decode(&user)
	if len(user.Id) == 0 {
		t.Error("expected some user id, but got nothing")
	}
	if len(user.EncryptedPassword) > 0 {
		t.Error("expected not to receive encrypted password in response")
	}
	if user.FirstName != params.FirstName {
		t.Errorf("expected firstname %s, but got %s", user.FirstName, params.FirstName)
	}
	if user.LastName != params.LastName {
		t.Errorf("expected lastname %s, but got %s", user.LastName, params.LastName)
	}
	if user.Email != params.Email {
		t.Errorf("expected email %s, but got %s", user.Email, params.Email)
	}
}

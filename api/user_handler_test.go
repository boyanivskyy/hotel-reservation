package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/boyanivskyy/hotel-reservation/db"
	"github.com/boyanivskyy/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
)

func Test_HandlePostUser(t *testing.T) {
	tdb := setup(t, db.TestDBNAME)
	defer tdb.tearDown(t)

	app := fiber.New(fiber.Config{
		ErrorHandler: ErrorHandler,
	})
	userHandler := NewUserHandler(tdb.User)
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
		t.Fatal(err)
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

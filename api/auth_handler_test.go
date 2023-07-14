package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/boyanivskyy/hotel-reservation/db"
	"github.com/boyanivskyy/hotel-reservation/db/fixtures"
	"github.com/gofiber/fiber/v2"
)

func TestAuthenticateSuccess(t *testing.T) {
	tdb := setup(t, db.TestDBNAME)
	defer tdb.tearDown(t)

	insertedUser := fixtures.AddUser(tdb.Store, "test", "test", false)

	app := fiber.New()

	authHandler := NewAuthHandler(tdb.User)
	app.Post("/authenticate", authHandler.HandleAuthenticate)

	params := AuthParams{
		Email:    "test@test.com",
		Password: "test_test",
	}
	b, _ := json.Marshal(params)

	req := httptest.NewRequest("POST", "/authenticate", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	users, _ := tdb.User.GetUsers(context.TODO())
	fmt.Println("users", users)

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected http status of 200 but got %d", resp.StatusCode)
	}

	var authResponse AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResponse); err != nil {
		t.Error(err)
	}

	if authResponse.Token == "" {
		t.Fatalf("expected the JWT token to be present in the authResponse")
	}

	// set EncryptedPassword to empty string as this was not
	// serialized to JSON format and does not remove EncryptedPassword from the object
	insertedUser.EncryptedPassword = ""
	if !reflect.DeepEqual(insertedUser, authResponse.User) {
		fmt.Println("insertedUser", insertedUser)
		fmt.Println("authResponse.User", authResponse.User)
		t.Fatal("expected the test user to same as authResponse.User")
	}
}

func TestAuthenticationWithWrongPasswordFailure(t *testing.T) {
	tdb := setup(t, db.TestDBNAME)
	defer tdb.tearDown(t)

	fixtures.AddUser(tdb.Store, "test", "test", false)

	app := fiber.New()

	authHandler := NewAuthHandler(tdb.User)
	app.Post("/authenticate", authHandler.HandleAuthenticate)

	params := AuthParams{
		Email:    "test@test.com",
		Password: "wrongpassword",
	}
	b, _ := json.Marshal(params)

	req := httptest.NewRequest("POST", "/authenticate", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected http status of 400 but got %d", resp.StatusCode)
	}

	var genResp genericResp
	if err := json.NewDecoder(resp.Body).Decode(&genResp); err != nil {
		t.Fatal(err)
	}

	if genResp.Type != "error" {
		t.Fatalf("expected genResp.Type to be error but got %s", genResp.Type)
	}

	if genResp.Msg != "invalid credentials" {
		t.Fatalf("expected genResp.Type to be <invalid credentials> but got %s", genResp.Msg)
	}
}

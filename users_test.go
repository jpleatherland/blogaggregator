package main

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/jpleatherland/blogaggregator/internal/database"
	_ "github.com/lib/pq"
)

var resources = Resources{}
var baseURL = "http://localhost:8080"

func TestMain(m *testing.M) {
	db, dbQueries, err := setupTestEnvironment()
	if err != nil {
		panic(err)
	}
	resources.DB = dbQueries
	code := m.Run()
	db.Close()
	teardownTestEnvironment()
	os.Exit(code)
}

func TestCreateUser(t *testing.T) {
	username := "Testing create user"
	user, err := createTestUser(&resources, username)
	if err != nil {
		t.Fatalf("test user creation failed: %v", err.Error())
	}
	if user.Name != username {
		t.Errorf("username does not match input, expected: %v, got: %v", username, user.Name)
	}
}

func TestGetUserByApiKey(t *testing.T) {
	username := "Testing get user"
	user, err := createTestUser(&resources, username)
	if err != nil {
		t.Fatalf("test user creation failed: %v", err.Error())
	}

	httpReq := httptest.NewRequest("GET", baseURL, nil)
	httpReq.Header.Add("Authorization", "ApiKey "+user.ApiKey)
	httpRecorder := httptest.NewRecorder()
	resources.getUser(httpRecorder, httpReq, user)
	responseBody, err := io.ReadAll(httpRecorder.Result().Body)
	if err != nil {
		t.Fatalf("unable to read response body: %v", err.Error())
	}

	dbUser := database.User{}
	err = json.Unmarshal(responseBody, &dbUser)
	if err != nil {
		t.Fatalf("unable to unmarshal response: %v", err.Error())
	}

	if dbUser != user {
		t.Errorf("users do not match expected %v, got %v", user, dbUser)
	}
}

package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/jpleatherland/blogaggregator/internal/database"
	_ "github.com/lib/pq"
)

var resources = Resources{}
var user = database.User{}
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
	httpRecorder, err := createTestUser(&resources)
	if err != nil {
		t.Fatalf("create test user failed: %v", err.Error())
	}
	if httpRecorder.Result().StatusCode != http.StatusCreated {
		t.Errorf("create user response status incorrect, expected: %v got: %v", http.StatusCreated, httpRecorder.Result().Status)
	}
	responseBody, err := io.ReadAll(httpRecorder.Result().Body)
	if err != nil {
		t.Fatalf("unable to read responseBody: %v", err.Error())
	}
	err = json.Unmarshal(responseBody, &user)
	if err != nil {
		t.Fatalf("unable to unmarshal response body: %v", err.Error())
	}
}

func TestGetUserByApiKey(t *testing.T) {
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

package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/joho/godotenv"
	"github.com/jpleatherland/blogaggregator/internal/database"
	"github.com/pressly/goose"
)

func setupTestEnvironment() (*sql.DB, *database.Queries, error) {
	godotenv.Load()
	DB_CONN_STRING := os.Getenv("TEST_DB_CONN_STRING")
	pgdb, err := sql.Open("postgres", DB_CONN_STRING+"?sslmode=disable")
	if err != nil {
		return nil, nil, err
	}
	// create database and required tables
	pgdb.Exec("CREATE DATABASE blogagg_test")
	db, err := sql.Open("postgres", DB_CONN_STRING+"/blogagg_test?sslmode=disable")
	if err != nil {
		return nil, nil, err
	}

	err = goose.Up(db, "./sql/schema")
	if err != nil {
		return nil, nil, err
	}

	dbQueries := database.New(db)

	log.Println("setup test database")
	return db, dbQueries, nil
}

func teardownTestEnvironment() {
	godotenv.Load()
	DB_CONN_STRING := os.Getenv("TEST_DB_CONN_STRING")
	pgdb, err := sql.Open("postgres", DB_CONN_STRING+"?sslmode=disable")
	if err != nil {
		log.Printf("failed to reopen database: %v", err.Error())
	}

	_, err = pgdb.Exec("DROP DATABASE blogagg_test")
	if err != nil {
		log.Printf("failed to teardown test environment: %v", err.Error())
	}
	log.Println("dropped test database")
}

func createTestUser(resources *Resources, username string) (database.User, error) {
	baseURL := "http://localhost:8080"
	user := database.User{}
	userBody := struct {
		Name string `json:"name"`
	}{
		Name: username,
	}

	payload, err := json.Marshal(userBody)
	if err != nil {
		return user, err
	}

	httpReq := httptest.NewRequest("POST", baseURL, bytes.NewBuffer(payload))
	httpRecorder := httptest.NewRecorder()
	resources.createUser(httpRecorder, httpReq)
	if httpRecorder.Result().StatusCode != http.StatusCreated {
		errorMsg := fmt.Sprintf("create user response status incorrect, expected: %v got: %v", http.StatusCreated, httpRecorder.Result().Status)
		return user, errors.New(errorMsg)
	}
	responseBody, err := io.ReadAll(httpRecorder.Result().Body)
	if err != nil {
		return user, errors.New("unable to read response body: " + err.Error())
	}
	err = json.Unmarshal(responseBody, &user)
	if err != nil {
		return user, errors.New("unable to unmarshal response body: " + err.Error())
	}
	return user, nil
}

func createTestFeed(resources *Resources, user database.User, feedName, feedUrl string) (feedResponse, error) {
	baseURL := "http://localhost:8080"
	feedResponse := feedResponse{
		Feed:        database.Feed{},
		Feed_Follow: database.FeedFollow{},
	}

	feedBody := struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}{
		Name: feedName,
		Url:  feedUrl,
	}

	payload, err := json.Marshal(feedBody)
	if err != nil {
		return feedResponse, errors.New("unable to marshal feed body: " + err.Error())
	}

	httpReq := httptest.NewRequest("GET", baseURL, bytes.NewBuffer(payload))
	httpReq.Header.Add("Authorization", "ApiKey "+user.ApiKey)
	httpRecorder := httptest.NewRecorder()
	resources.createFeed(httpRecorder, httpReq, user)

	expectedStatus := http.StatusCreated
	actualStatus := httpRecorder.Result().StatusCode

	if expectedStatus != actualStatus {
		errMsg := fmt.Sprintf("incorrect status received. expected %v got %v", expectedStatus, actualStatus)
		return feedResponse, errors.New(errMsg)
	}

	responseBody, err := io.ReadAll(httpRecorder.Result().Body)
	if err != nil {
		return feedResponse, errors.New("unable to read response body: " + err.Error())
	}

	err = json.Unmarshal(responseBody, &feedResponse)
	if err != nil {
		return feedResponse, errors.New("unable to unmarshal response body: " + err.Error())
	}

	return feedResponse, nil
}

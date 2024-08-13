package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"log"
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

func createTestUser(resources *Resources) (*httptest.ResponseRecorder, error) {
	var baseURL = "http://localhost:8080"
	userBody := struct {
		Name string `json:"name"`
	}{
		Name: "John",
	}

	payload, err := json.Marshal(userBody)
	if err != nil {
		return nil, err
	}

	httpReq := httptest.NewRequest("POST", baseURL, bytes.NewBuffer(payload))
	httpRecorder := httptest.NewRecorder()
	resources.createUser(httpRecorder, httpReq)
	return httpRecorder, nil
}

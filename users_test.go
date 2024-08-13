package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/jpleatherland/blogaggregator/internal/database"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

var resources = Resources{}

func TestMain(m *testing.M) {
	db, err := setupTestEnvironment()
	if err != nil {
		panic(err)
	}
	code := m.Run()
	db.Close()
	teardownTestEnvironment()
	os.Exit(code)
}

func TestCreateUser(t *testing.T) {
	userBody := struct {
		Name string `json:"name"`
	}{
		Name: "John",
	}

	payload, err := json.Marshal(userBody)
	if err != nil {
		t.Fatalf("unable to marshal payload: %v", err.Error())
	}

	httpReq := httptest.NewRequest("POST", "http://localhost:8080/", bytes.NewBuffer(payload))
	httpRecorder := httptest.NewRecorder()
	resources.createUser(httpRecorder, httpReq)
	if httpRecorder.Result().StatusCode != http.StatusCreated {
		t.Errorf("create user response status incorrect, expected: %v got: %v", http.StatusCreated, httpRecorder.Result().Status)
	}
}

func setupTestEnvironment() (*sql.DB, error) {
	godotenv.Load()
	DB_CONN_STRING := os.Getenv("TEST_DB_CONN_STRING")
	pgdb, err := sql.Open("postgres", DB_CONN_STRING+"?sslmode=disable")
	if err != nil {
		return nil, err
	}
	// create database and required tables
	pgdb.Exec("CREATE DATABASE blogagg_test")
	db, err := sql.Open("postgres", DB_CONN_STRING+"/blogagg_test?sslmode=disable")
	if err != nil {
		return nil, err
	}

	err = goose.Up(db, "./sql/schema")
	if err != nil {
		return nil, err
	}

	dbQueries := database.New(db)
	resources.DB = dbQueries

	log.Println("setup test database")
	return db, nil
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

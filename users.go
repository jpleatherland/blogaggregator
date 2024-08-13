package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jpleatherland/blogaggregator/internal/database"
	"golang.org/x/net/context"
)

func (resources *Resources) createUser(rw http.ResponseWriter, req *http.Request) {
	var dbUser = database.CreateUserParams{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&dbUser)
	if err != nil {
		respondWithError(rw, http.StatusBadRequest, err.Error())
		return
	}

	ctx := context.Background()

	currTime := time.Now()

	dbUser.ID = uuid.New()
	dbUser.CreatedAt = currTime
	dbUser.UpdatedAt = currTime

	writtenUser, err := resources.DB.CreateUser(ctx, dbUser)
	if err != nil {
		respondWithError(rw, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(rw, http.StatusCreated, writtenUser)
}

func (resources *Resources) getUser(rw http.ResponseWriter, req *http.Request, user database.User) {
	respondWithJSON(rw, http.StatusOK, user)
}

package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jpleatherland/blogaggregator/internal/database"
)

func (resources *Resources) createFeedFollow(rw http.ResponseWriter, req *http.Request, user database.User) {
	decoder := json.NewDecoder(req.Body)
	feedFollowBody := database.CreateFeedFollowParams{}
	err := decoder.Decode(&feedFollowBody)
	if err != nil {
		respondWithError(rw, http.StatusBadRequest, "unable to read body")
		return
	}

	feedFollowBody.UserID = user.ApiKey

	writtenFeed, err := resources.newFeedFollow(feedFollowBody)
	if err != nil {
		respondWithError(rw, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(rw, http.StatusCreated, writtenFeed)
}

func (resources *Resources) getFeedFollowsByUserId(rw http.ResponseWriter, _ *http.Request, user database.User) {
	ctx := context.Background()
	feeds, err := resources.DB.GetFeedFollowsByUserId(ctx, user.ApiKey)
	if err != nil {
		respondWithError(rw, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(rw, http.StatusOK, feeds)
}

func (resources *Resources) deleteFeedFollow(rw http.ResponseWriter, req *http.Request, user database.User) {
	feedFollowId := req.PathValue("feedFollowID")
	feedFollowUUID, err := uuid.Parse(feedFollowId)
	if err != nil {
		respondWithError(rw, http.StatusBadRequest, "bad feed follow id")
		return
	}
	ctx := context.Background()
	err = resources.DB.DeleteFeedFollow(ctx, feedFollowUUID)
	if err != nil {
		respondWithError(rw, http.StatusInternalServerError, err.Error())
		return
	}
	rw.WriteHeader(http.StatusNoContent)
}

func (resources *Resources) newFeedFollow(feedFollowPayload database.CreateFeedFollowParams) (database.FeedFollow, error) {
	currTime := time.Now()
	ctx := context.Background()

	feedFollowPayload.ID = uuid.New()
	feedFollowPayload.CreatedAt = currTime
	feedFollowPayload.UpdatedAt = currTime

	writtenFeed, err := resources.DB.CreateFeedFollow(ctx, feedFollowPayload)
	if err != nil {
		return database.FeedFollow{}, err
	}
	return writtenFeed, nil
}

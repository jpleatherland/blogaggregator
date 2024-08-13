package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jpleatherland/blogaggregator/internal/database"
)

func (resources *Resources) createFeed(rw http.ResponseWriter, req *http.Request, user database.User) {
	decoder := json.NewDecoder(req.Body)
	feedBody := database.CreateFeedParams{}
	err := decoder.Decode(&feedBody)
	if err != nil {
		respondWithError(rw, http.StatusBadRequest, "unable to read body")
		return
	}

	currTime := time.Now()
	ctx := context.Background()

	feedBody.ID = uuid.New()
	feedBody.CreatedAt = currTime
	feedBody.UpdatedAt = currTime
	feedBody.UserID = user.ApiKey

	writtenFeed, err := resources.DB.CreateFeed(ctx, feedBody)
	if err != nil {
		respondWithError(rw, http.StatusInternalServerError, err.Error())
		return
	}

	feedFollowBody := database.CreateFeedFollowParams{}
	feedFollowBody.FeedID = writtenFeed.ID
	feedFollowBody.UserID = user.ApiKey

	writtenFeedFollow, err := resources.newFeedFollow(feedFollowBody)
	if err != nil {
		log.Printf("feed created but follow creation failed with: %v", err.Error())
		respondWithJSON(rw, http.StatusCreated, map[string]interface{}{"feed": writtenFeed, "feed_follow": "feed created but failed to follow"})
		return
	}

	respondWithJSON(rw, http.StatusCreated, map[string]interface{}{"feed": writtenFeed, "feed_follow": writtenFeedFollow})
}

func (resources *Resources) getAllFeeds(rw http.ResponseWriter, _ *http.Request) {
	ctx := context.Background()
	feeds, err := resources.DB.GetFeeds(ctx)
	if err != nil {
		respondWithError(rw, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(rw, http.StatusOK, feeds)
}

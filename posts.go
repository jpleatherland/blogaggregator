package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/jpleatherland/blogaggregator/internal/database"
)

func (resources *Resources) createPost(post RSSItem, feedID uuid.UUID) error {
	currTime := time.Now()
	postParams := database.CreatePostParams{
		ID:        uuid.New(),
		CreatedAt: currTime,
		UpdatedAt: currTime,
		Title:     post.Title,
		Url:       post.Link,
		FeedID:    feedID,
	}
	if post.Description != "" {
		postParams.Description = sql.NullString{String: post.Description, Valid: true}
	}
	if post.PubDate != "" {
		t, err := identifyDateFormat(post.PubDate)
		if err != nil {
			return err
		} else {
			postParams.PublishedAt = sql.NullTime{Time: t, Valid: true}
		}
	}
	ctx := context.Background()
	err := resources.DB.CreatePost(ctx, postParams)
	if err != nil {
		errMsg := ("Failed to create post: " + err.Error())
		return errors.New(errMsg)
	}
	return nil
}

func (resources *Resources) getPostsByUser(rw http.ResponseWriter, req *http.Request, user database.User) {
	limit, err := strconv.Atoi(req.URL.Query().Get("limit"))
	if err != nil {
		limit = 10
	}
	ctx := context.Background()
	getPostParams := database.GetPostsByUserParams{UserID: user.ApiKey, Limit: int32(limit)}
	posts, err := resources.DB.GetPostsByUser(ctx, getPostParams)
	if err != nil {
		http.Error(rw, "unable to read posts from database: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if len(posts) == 0 {
		http.Error(rw, "no posts found for user", http.StatusNotFound)
		return
	}
	postResponse, err := json.Marshal(posts)
	if err != nil {
		http.Error(rw, "unable to parse response", http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusOK)
	rw.Write(postResponse)
}

func identifyDateFormat(dateString string) (time.Time, error) {
	layouts := []string{
		time.RFC3339,
		time.RFC1123Z,
		"2006-01-02 15:04:05",
		"02/01/2006 15:04:05",
		"02-Jan-2006",
		// Add more layouts as needed
	}

	var parsedTime time.Time
	var err error
	for _, layout := range layouts {
		parsedTime, err = time.Parse(layout, dateString)
		if err == nil {
			return parsedTime, nil
		}
	}
	return time.Time{}, fmt.Errorf("unknown date format")
}

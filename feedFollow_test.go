package main

import (
	"io"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jpleatherland/blogaggregator/internal/database"
)

func TestCreateFeedFollow(t *testing.T) {
	usernames := [2]string{"Test user create feed follow 1", "Test user create feed follow 2"}
	users := [2]database.User{}
	for i, user := range usernames {
		user, err := createTestUser(&resources, user)
		if err != nil {
			t.Fatal("unable to create user")
		}
		users[i] = user
	}
	feedName := "Test feed name create feed follow 1"
	feedUrl := "Test feed url create feed follow 1"

	dbFeed, err := createTestFeed(&resources, users[0], feedName, feedUrl)
	if err != nil {
		t.Fatalf("unable to create feed: %v", err.Error())
	}

	createFeedBody := struct {
		Feed_ID string `json:"feed_id"`
	}{
		Feed_ID: dbFeed.Feed.ID.String(),
	}

	payload, err := json.Marshal(createFeedBody)
	if err != nil {
		t.Fatalf("unable to marshall create feed payload: %v", err.Error())
	}

	httpReq := httptest.NewRequest("POST", baseURL, bytes.NewBuffer(payload))
	httpRecorder := httptest.NewRecorder()
	resources.createFeedFollow(httpRecorder, httpReq, users[1])
	if httpRecorder.Result().StatusCode != http.StatusCreated {
		t.Errorf("create feed follow status incorrect, expected: %v got: %v", http.StatusCreated, httpRecorder.Result().Status)
	}
	
	responseBody, err := io.ReadAll(httpRecorder.Result().Body)
	if err != nil {
		t.Errorf("unable to read response body: " + err.Error())
	}

	feedFollow := database.FeedFollow{}
	err = json.Unmarshal(responseBody, &feedFollow)
	if err != nil {
		t.Errorf("unable to unmarshal response body: " + err.Error())
	}
	if feedFollow.FeedID != dbFeed.Feed.ID {
		t.Errorf("feed follow ID does not match input. Expected %v, got %v", dbFeed.Feed.ID, feedFollow.FeedID)
	}

	if feedFollow.UserID != users[1].ApiKey {
		t.Errorf("feed follow user id does not match input. Expected %v, got %v", users[1].ApiKey, feedFollow.UserID)
	}

}
func TestGetFeedFollowsByUserId(t *testing.T) {
	username := "Test get feed follows by user id"
	feedNames := [2]string{"Test get feed names by user id name 1", "Test get feed names by user id name 2"}
	feedUrls := [2]string{"Test get feed names by user id url 1", "Test get feed names by user id url 2"}
	user, err := createTestUser(&resources, username)
	if err != nil {
		t.Fatalf("unable to create test user: %v", err.Error())
	}
	feeds := [2]feedResponse{}
	for i := range 2 {
		feed, err := createTestFeed(&resources, user, feedNames[i], feedUrls[i])
		if err != nil {
			t.Fatalf("unable to create test feed %d: %v", i, err.Error())
		}
		feeds[i] = feed
	}
	dbFeeds, err :=	getTestFeedFollowsByUserId(&resources, user)
	if err != nil {
		t.Fatalf("unable to get feeds by User ID: %v", err.Error())
	}

	for _, feed := range feeds {
		found := false
		for _, dbFeed := range dbFeeds {
			if dbFeed.FeedID == feed.Feed.ID {
				found = true
			}
			if dbFeed.UserID != user.ApiKey {
				t.Errorf("dbFeed user id does not match user id. Want %v got %v", user.ApiKey, dbFeed.UserID)
			}
		}
		if !found {
			t.Errorf("feed id not found in response: %v", feed.Feed.ID)
		}
	}
}

func TestDeleteFeedFollow(t *testing.T) {
	username := "Test delete follow feed user"
	feedNames := [2]string{"Test delete follow feed name 1", "Test delete follow feed name 2"}
	feedUrls := [2]string{"Test delete follow feed url 1", "Test delete follow feed url 2"}
	feeds := [2]feedResponse{}

	user, err := createTestUser(&resources, username)
	if err != nil {
		t.Fatalf("unable to create test user: %v", err.Error())
	}

	for i := range 2 {
		feed, err := createTestFeed(&resources, user, feedNames[i], feedUrls[i])
		if err != nil {
			t.Fatalf("unable to create test feed: %v", err.Error())
		}
		feeds[i] = feed
	}

	httpReq := httptest.NewRequest("DELETE", "http://localhost:8080/v1/feed_follows/"+feeds[1].Feed_Follow.ID.String(), nil)
	httpRecorder := httptest.NewRecorder()
	resources.deleteFeedFollow(httpRecorder, httpReq, user)

	expectedStatusCode := http.StatusNoContent
	actualStatusCode := httpRecorder.Result().StatusCode

	if actualStatusCode != expectedStatusCode {
		t.Errorf("incorrect status returned. want %v got %v", expectedStatusCode, actualStatusCode)
	}

	finalFeeds, err := getTestFeedFollowsByUserId(&resources, user)
	if err != nil {
		t.Fatalf("unable to get feed follows by user id: %v", err.Error())
	}

	for _, feed := range finalFeeds {
		if feed.FeedID == feeds[1].Feed_Follow.ID {
			t.Errorf("feed %v has not been deleted", feed.FeedID)
		}
	}
	
}

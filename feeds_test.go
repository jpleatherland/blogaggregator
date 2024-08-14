package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http/httptest"
	"testing"

	"github.com/jpleatherland/blogaggregator/internal/database"
)

func TestCreateFeed(t *testing.T) {
	username := "Testing create feed user"
	user, err := createTestUser(&resources, username)
	if err != nil {
		t.Fatalf("test user creation failed: %v", err.Error())
	}

	feedName := "test feed name 1"
	feedUrl := "test feed url 1"
	feed, err := createTestFeed(&resources, user, feedName, feedUrl)
	if err != nil {
		t.Fatalf("test feed creation failed: %v", err.Error())
	}

	if feed.Feed.Name != feedName {
		t.Errorf("test create feed incorrect name. expected %v got %v", feedName, feed.Feed.Name)
	}

	expectedUserId := user.ApiKey
	gotUserId := feed.Feed.UserID

	if expectedUserId != gotUserId {
		t.Errorf("test create feed incorrect user attribution. want %v got %v", expectedUserId, gotUserId)
	}
}

func TestGetAllFeeds(t *testing.T) {
	//create user & feed 1
	usernames := [2]string{"Testing get feed user", "Testing get feed user 2"}
	user, err := createTestUser(&resources, usernames[0])
	if err != nil {
		t.Fatalf("test user creation failed: %v", err.Error())
	}

	feedNames := [2]string{"test get feed name 1", "test get feed name 2"}
	feedUrls := [2]string{"test get feed url 1", "test get feed url 2"}
	_, err = createTestFeed(&resources, user, feedNames[0], feedUrls[0])
	if err != nil {
		t.Fatalf("test feed creation failed: %v", err.Error())
	}

	//create user & feed 2
	user2, err := createTestUser(&resources, usernames[1])
	if err != nil {
		t.Fatalf("test user creation failed: %v", err.Error())
	}

	_, err = createTestFeed(&resources, user2, feedNames[1], feedUrls[1])
	if err != nil {
		t.Fatalf("test feed creation failed: %v", err.Error())
	}

	httpReq := httptest.NewRequest("GET", baseURL, nil)
	httpRecorder := httptest.NewRecorder()
	resources.getAllFeeds(httpRecorder, httpReq)
	responseBody, err := io.ReadAll(httpRecorder.Result().Body)
	if err != nil {
		t.Fatalf("unable to read response body: %v", err.Error())
	}

	var dbFeeds []database.Feed
	err = json.Unmarshal(responseBody, &dbFeeds)
	if err != nil {
		t.Fatalf("unable to unmarshal response: %v", err.Error())
	}

	for _, feedName := range feedNames {
		found := false
		for _, dbFeed := range dbFeeds {
			if dbFeed.Name == feedName {
				found = true
			}
		}
		if !found {
			t.Errorf("test feed name not found in result: %v", feedName)
		}
	}
}

func TestFetchFeeds(t *testing.T) {
	username := "Test user fetch feeds"
	user, err := createTestUser(&resources, username)
	if err != nil {
		t.Fatalf("test user creation failed: %v", err.Error())
	}

	feedNames := [2]string{"boot.dev", "wagslane"}
	feedUrls := [2]string{"https://blog.boot.dev/index.xml", "https://wagslane.dev/index.xml"}

	for i := range 2 {
		_, err := createTestFeed(&resources, user, feedNames[i], feedUrls[i])
		if err != nil {
			t.Fatalf("failed to create test feed: %v with %v", feedNames[i], err.Error())
		}
	}

	log.Printf("set up ok")

	resources.fetchFeeds()
}

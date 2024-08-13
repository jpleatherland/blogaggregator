package main

import (
	"testing"
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

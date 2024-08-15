package main

import (
	"testing"

	"github.com/google/uuid"
)

func TestCreatePost(t *testing.T) {
	user, err := createTestUser(&resources, "Test post user name")
	if err != nil {
		t.Fatalf("failed to create user: %v", err.Error())
	}
	feedName := "Test post feed name"
	feedUrl := "Test post feed url"

	testFeed, err := createTestFeed(&resources, user, feedName, feedUrl)
	if err != nil {
		t.Fatalf("failed to create test feed: %v with %v", feedName, err.Error())
	}

	rss := RSSItem{
		Title:       "RSS Item Create Posts",
		Link:        "RSS Link Create Posts",
		Description: "Test RSS Desc Create Posts",
		Author:      "Test RSS Author Create Posts",
		Category:    "Test RSS Category Create Posts",
		PubDate:     "2023-05-02T09:34:01Z",
		Guid:        uuid.New().String(),
	}

	err = resources.createPost(rss, testFeed.Feed.ID)
	if err != nil {
		t.Errorf("failed to create post: %v", err.Error())
	}
}

func TestGetPostsByUser(t *testing.T) {
	user, err := createTestUser(&resources, "Test get post user name")
	if err != nil {
		t.Fatalf("failed to create user: %v", err.Error())
	}
	feedName := "Test get post feed name"
	feedUrl := "Test get post feed url"

	testFeed, err := createTestFeed(&resources, user, feedName, feedUrl)
	if err != nil {
		t.Fatalf("failed to create test feed: %v with %v", feedName, err.Error())
	}

	rss := RSSItem{
		Title:       "RSS Item Get Posts",
		Link:        "RSS Link Get Posts",
		Description: "Test RSS Desc Get Posts",
		Author:      "Test RSS Author Get Posts",
		Category:    "Test RSS Category Get Posts",
		PubDate:     "2023-05-02T09:34:01Z",
		Guid:        uuid.New().String(),
	}

	err = resources.createPost(rss, testFeed.Feed.ID)
	if err != nil {
		t.Fatalf("failed to create post: %v", err.Error())
	}

	posts, err := getTestPostsByUser(&resources, user)
	if err != nil {
		t.Fatalf("failed to get posts by user: %v", err.Error())
	}

	gotPostTitle := posts[0].Title
	wantPostTitle := rss.Title

	if gotPostTitle != wantPostTitle {
		t.Errorf("incorrect post title: got %v want %v", gotPostTitle, wantPostTitle)
	}

}

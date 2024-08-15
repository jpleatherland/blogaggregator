package main

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jpleatherland/blogaggregator/internal/database"
)

type feedResponse struct {
	Feed        database.Feed       `json:"feed"`
	Feed_Follow database.FeedFollow `json:"feed_follow"`
}

type RSS struct {
	XMLName xml.Name   `xml:"rss"`
	Version string     `xml:"version.attr"`
	Channel RSSChannel `xml:"channel"`
}

type RSSChannel struct {
	Title         string    `xml:"title"`
	Link          string    `xml:"link"`
	Description   string    `xml:"description"`
	Language      string    `xml:"language"`
	PubDate       string    `xml:"pubDate"`
	LastBuildDate string    `xml:"lastBuildDate"`
	Generator     string    `xml:"generator"`
	Items         []RSSItem `xml:"item"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description,omitempty"`
	Author      string `xml:"author,omitempty"`
	Category    string `xml:"category,omitempty"`
	PubDate     string `xml:"pubDate"`
	Guid        string `xml:"guid"`
}

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

	response := feedResponse{Feed: writtenFeed, Feed_Follow: writtenFeedFollow}

	respondWithJSON(rw, http.StatusCreated, response)
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

func getDataFromFeed(url string) (RSS, error) {
	rss := RSS{}
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error fetching RSS feed: ", err)
		return rss, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body: ", err)
		return rss, err
	}

	err = xml.Unmarshal(data, &rss)
	if err != nil {
		fmt.Println("Error unmarshalling RSS: ", err)
		return rss, err
	}

	return rss, nil
}

func (resources *Resources) fetchFeeds() {
	log.Print("starting feed fetch")
	ctx := context.Background()
	feeds, err := resources.DB.GetNextFeedsToFetch(ctx, 10)
	if err != nil {
		log.Print(err.Error())
		return
	}

	var waitGroup sync.WaitGroup
	results := make(chan RSS, len(feeds))
	waitGroup.Add(len(feeds))

	for i := range feeds {
		go func() {
			defer waitGroup.Done()
			log.Println("getting data for: " + feeds[i].Url)
			result, err := getDataFromFeed(feeds[i].Url)
			if err != nil {
				log.Printf("failed to get data from feed %v: %v ", feeds[i].Name, err.Error())
			}
			results <- result
		}()
	}

	waitGroup.Wait()
	close(results)
	log.Println("Wait group concluded, results channel closed")

	chanIdx := 0
	for result := range results {
		for _, item := range result.Channel.Items {
			err := resources.createPost(item, feeds[chanIdx].ID)
			if err != nil {
				log.Println(err.Error())
			}
		}
	}
	chanIdx++
}

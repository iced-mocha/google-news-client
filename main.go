package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/iced-mocha/shared/models"
)

var feedURL string

type Response struct {
	Status   string
	Source   string
	SortBy   string
	Articles []Article
}

type Article struct {
	Author      string
	Title       string
	Description string
	URL         string
	URLToImage  string
	PublishedAt time.Time
}

func GetPosts(w http.ResponseWriter, r *http.Request) {
	var err error

	resp, err := http.Get(feedURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	var respBody Response
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	postsToReturn := make([]models.Post, 0, 10)
	for _, a := range respBody.Articles {
		p := models.Post{
			Date: a.PublishedAt,
			Author: a.Author,
			Title: a.Title,
			Content: a.Description,
			HeroImg: a.URLToImage,
			PostLink: a.URL,
			Platform: models.PlatformGoogleNews,
		}
		postsToReturn = append(postsToReturn, p)
	}

	res, err := json.Marshal(postsToReturn)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

func main() {
	apiKey := os.Getenv("GOOGLE_NEWS_API_KEY")
	feedURL = "https://newsapi.org/v1/articles?source=google-news&sortBy=top&apiKey=" + apiKey
	r := mux.NewRouter()
	r.HandleFunc("/v1/posts", GetPosts).Methods(http.MethodGet)
	log.Fatal(http.ListenAndServe(":7000", r))
}

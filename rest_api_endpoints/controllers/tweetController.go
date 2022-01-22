package controllers

import (
	"encoding/json"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/amulcse/models"
	"github.com/amulcse/utils"
)

func GetTweets(w http.ResponseWriter, r *http.Request) {

	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))
	sortBy := r.FormValue("creationtime")
	username := r.FormValue("username")

	if count > 50 || count < 1 {
		count = 50
	}
	if start < 0 {
		start = 0
	}

	allTweets := models.FetchTweets()
	response := allTweets

	if len(allTweets) > 0 {

		// filter
		if username != "" {
			temp := []models.Tweet{}
			for _, tweet := range response {
				if tweet.Username == username {
					temp = append(temp, tweet)
				}
			}
			response = temp
		}

		// sorting
		if sortBy != "" {
			if sortBy == "asc" {
				sort.Slice(response, func(i, j int) bool {
					return response[i].Created.Before(response[j].Created)
				})
			} else {
				sort.Slice(response, func(i, j int) bool {
					return response[i].Created.After(response[j].Created)
				})

			}
		}

		// pagination
		end := start + count
		if end > len(response) {
			end = len(response)
		}
		response = response[start:end]

	}

	w.Header().Set("Content-Type", "applicaiton/json")
	json.NewEncoder(w).Encode((response))
}

func AddTweet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	tweet := models.Tweet{}
	if err := json.NewDecoder(r.Body).Decode(&tweet); err != nil {
		utils.InternalServerError(w, err)
		return
	}

	if !utils.ValidateRequest(tweet, w) {
		return
	}

	tweet.Created = time.Now()
	models.AddTweet(tweet)
	if err := json.NewEncoder(w).Encode(tweet); err != nil {
		utils.InternalServerError(w, err)
		return
	}
}

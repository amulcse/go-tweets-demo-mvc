package routes

import (
	"github.com/amulcse/controllers"

	"github.com/gorilla/mux"
)

var TweetRouter = func(router *mux.Router) {
	router.HandleFunc("/tweets", controllers.GetTweets).Methods("GET")
	router.HandleFunc("/tweet", controllers.AddTweet).Methods("POST")
}

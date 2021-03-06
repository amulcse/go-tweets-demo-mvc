package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

	textgenerator "github.com/liderman/text-generator"
)

var responses []string

var wg sync.WaitGroup
var mut sync.Mutex

const MAX = 20

func getTweets() {
	resp, err := http.Get("http://localhost:8000/tweets")
	if err != nil {
		log.Fatalln(err)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("All Tweets:", string(body))
}

func randomNumber(min int, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return min + rand.Intn(max-min)
}

func getAutoGeneratedTweet() string {
	tg := textgenerator.New()
	template := "{ Good {morning|evening|day}, {How are you | Whats new with you} I am learning { GO | PyThon | NodeJS | Rust }. It seems awesome language. { fuck | luck | New | Check | Work | Life | TEST | Mobile | Mouse | Pen | Lorem Ipsum is simply dummy text of the printing and typesetting indu Lorem Ipsum is simply dummy text of the printing and typesetting indu} }"
	message := tg.Generate(template)
	return message
	// fmt.Print(message)
}

func addTweet() {

	// defer wg.Done()

	usernames := [5]string{"amul", "ruslan", "jay", "sky", "alice"}

	tweet := map[string]string{"username": usernames[randomNumber(1, 5)], "message": getAutoGeneratedTweet()}
	tweetJSON, _ := json.Marshal(tweet)
	req, err := http.NewRequest("POST", "http://localhost:8000/tweet", bytes.NewBuffer(tweetJSON))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)

	// only bad request log
	if res.StatusCode != 200 {
		mut.Lock()
		fmt.Printf("[%d] - %s %s\n", res.StatusCode, time.Now(), body)
		defer mut.Unlock()
	}

	// keep record of all the responses
	responses = append(responses, string(body))
	// time.Sleep(time.Second)

}

func main() {
	queue := make(chan bool, MAX)
	for i := 1; i <= 10000; i++ {
		queue <- true
		go func() {
			defer func() {
				<-queue
			}()
			addTweet()
			// wg.Add(1)
		}()

	}

	for i := 1; i <= MAX; i++ {
		queue <- true
	}

	// wg.Wait()
	// getTweets()

}

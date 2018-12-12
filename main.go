package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"github.com/tidwall/gjson"
	"strings"
)

const SR_URL = "https://www.reddit.com/r/%s/new.json?sort=new"
const R_URL = "https://www.reddit.com/new.json?sort=new"

type RedditPost struct {
	Title		string
	Content		string
	Author		string
	URL			string
}

func main() {
	subreddit := flag.String("r", "reddit", "Subreddit to monitor for keyword")
	keyword := flag.String("k", "", "Words to monitor for Reddit posts")
	flag.Parse()

	if *keyword == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	// call API
	data := redditData(*subreddit)

	if data != nil {
		for _, val := range data {
			if strcheck(val.Title, *keyword) || strcheck(val.Content, *keyword) {
				fmt.Printf("Possible Match: %s by %s (%s)\n", val.Title, val.Author, val.URL)
			}
		}
	}
	fmt.Printf("Watching %s for %s", *subreddit, *keyword)
}

func redditData(subreddit string) []*RedditPost {
	rData := make([]*RedditPost, 0)
	client := &http.Client{}
	rURL := ""

	if subreddit == "reddit" {
		rURL = R_URL
	} else {
		rURL = fmt.Sprintf(SR_URL, subreddit)
	}

	req, err := http.NewRequest("GET", rURL, nil)
	req.Header.Add("User-Agent", "Reddit Tracker")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error reaching Reddit: %s", err.Error())
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed reading Reddit response: %s", err.Error())
	}

	var result map[string]interface{}
	json.Unmarshal([]byte(body), &result)

	if result["error"] != nil {
		log.Fatalf("Response error: %s", result["message"])
	}

	data := gjson.Get(string(body), "data.children")
	for _, value := range data.Array() {
		curPost := &RedditPost{
			Title: gjson.Get(value.String(), "data.title").String(),
			Content: gjson.Get(value.String(), "data.selftext").String(),
			Author: gjson.Get(value.String(), "data.author").String(),
			URL: "https://reddit.com" + gjson.Get(value.String(), "data.permalink").String(),
		}
		rData = append(rData, curPost)
	}

	return rData
}

func strcheck(str1, str2 string) bool {
	return strings.Contains(strings.ToLower(str1), strings.ToLower(str2))
}
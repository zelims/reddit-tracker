package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	subreddit := flag.String("r", "reddit", "Subreddit to monitor for keyword")
	keyword := flag.String("k", "", "Words to monitor for Reddit posts")
	flag.Parse()

	if *keyword == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	fmt.Printf("Watching %s for %s", *subreddit, *keyword)
}

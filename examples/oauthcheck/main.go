package main

import (
	"fmt"
	"os"

	"github.com/masci/flickr"
	"github.com/masci/flickr/auth/oauth"
)

func main() {
	// retrieve Flickr credentials from env vars
	apik := os.Getenv("FLICKRGO_API_KEY")
	apisec := os.Getenv("FLICKRGO_API_SECRET")
	token := os.Getenv("FLICKRGO_OAUTH_TOKEN")

	// do not proceed if credentials were not provided
	if apik == "" || apisec == "" || token == "" {
		fmt.Fprintln(os.Stderr, "Please set FLICKRGO_API_KEY, FLICKRGO_API_SECRET "+
			"and FLICKRGO_OAUTH_TOKEN env vars")
		os.Exit(1)
	}

	// create an API client with credentials
	client := flickr.NewFlickrClient(apik, apisec)

	response, _ := oauth.CheckToken(client, token)
	fmt.Println(fmt.Sprintf("%+v", *response))
}

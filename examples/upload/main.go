package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/masci/flickr.go/flickr"
)

func main() {
	// retrieve Flickr credentials from env vars
	apik := os.Getenv("FLICKRGO_API_KEY")
	apisec := os.Getenv("FLICKRGO_API_SECRET")
	token := os.Getenv("FLICKRGO_OAUTH_TOKEN")
	tokenSecret := os.Getenv("FLICKRGO_OAUTH_TOKEN_SECRET")

	// do not proceed if credentials were not provided
	if apik == "" || apisec == "" || token == "" || tokenSecret == "" {
		fmt.Fprintln(os.Stderr, "Please set FLICKRGO_API_KEY, FLICKRGO_API_SECRET "+
			"and FLICKRGO_OAUTH_TOKEN env vars")
		os.Exit(1)
	}

	// create an API client with credentials
	client := flickr.NewFlickrClient(apik, apisec)
	client.OAuthToken = token
	client.OAuthTokenSecret = tokenSecret

	path, _ := filepath.Abs("examples/upload/gopher.jpg")
	id, err := flickr.UploadPhoto(client, path, nil)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Photo uploaded, id:", id)
	}

	params := flickr.NewUploadParams()
	params.Title = "A Gopher"
	id, err = flickr.UploadPhoto(client, path, params)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Photo uploaded, id:", id)
	}
}

package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/masci/flickr.go/flickr"
	"github.com/masci/flickr.go/flickr/photos"
	"github.com/masci/flickr.go/flickr/photosets"
)

func main() {
	var pause = func() {
		var foo string
		fmt.Println("Press a key to continue")
		fmt.Scanln(&foo)
	}

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

	// upload a photo
	path, _ := filepath.Abs("examples/upload/gopher.jpg")
	params := flickr.NewUploadParams()
	params.Title = "A Gopher"
	resp, err := flickr.UploadPhoto(client, path, params)
	if err != nil {
		fmt.Println("Failed uploading:", err, resp.ErrorMsg())
		os.Exit(1)
	} else {
		fmt.Println("Photo uploaded, id:", resp.Id)
		pause()
	}

	// create a photoset
	respS, err := photosets.Create(client, "A Set", "", resp.Id)
	if err != nil {
		fmt.Println("Failed creating set:", respS.ErrorMsg())
		os.Exit(1)
	} else {
		fmt.Println("Set created, id:", respS.Set.Id, "url:", respS.Set.Url)
		pause()
	}

	// TODO assign photo to the photoset
	// TODO remove the photo from the photoset

	// delete the photo
	respD, err := photos.Delete(client, resp.Id)
	if err != nil {
		fmt.Println("Failed deleting photo:", err)
		fmt.Println(respD.ErrorMsg())
		os.Exit(1)
	} else {
		fmt.Println("Successfully removed photo")
		pause()
	}
}

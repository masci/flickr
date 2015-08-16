package main

import (
	"fmt"
	"os"

	"github.com/masci/flickr"
	"github.com/masci/flickr/photosets"
)

func main() {
	// retrieve Flickr credentials from env vars
	apik := os.Getenv("FLICKRGO_API_KEY")
	apisec := os.Getenv("FLICKRGO_API_SECRET")
	token := os.Getenv("FLICKRGO_OAUTH_TOKEN")
	tokenSecret := os.Getenv("FLICKRGO_OAUTH_TOKEN_SECRET")
	nsid := os.Getenv("FLICKRGO_USER_ID")

	// do not proceed if credentials were not provided
	if apik == "" || apisec == "" || token == "" || tokenSecret == "" {
		fmt.Fprintln(os.Stderr, "Please set FLICKRGO_API_KEY, FLICKRGO_API_SECRET "+
			"and FLICKRGO_OAUTH_TOKEN, FLICKRGO_OAUTH_TOKEN_SECRET env vars")
		os.Exit(1)
	}

	// create an API client with credentials
	client := flickr.NewFlickrClient(apik, apisec)
	client.OAuthToken = token
	client.OAuthTokenSecret = tokenSecret
	client.Id = nsid

	/*
		response, _ := photosets.GetList(client, false, "23148015@N00", 1)
		fmt.Println(fmt.Sprintf("%+v", *response))

		response, _ := photosets.GetPhotos(client, false, "72157632076344815", "23148015@N00", 1)
		fmt.Println(fmt.Sprintf("%+v", *response))

		response, _ := photosets.EditMeta(client, "72157654143356943", "bar", "Baz")
		fmt.Println(fmt.Sprintf("%+v", *response))

		response, _ := photosets.EditPhotos(client, "72157654143356943", "9518691684", []string{"9518691684", "19681581995"})
		fmt.Println(fmt.Sprintf("%+v", *response))

		response, _ := photosets.RemovePhotos(client, "72157654143356943", []string{"9518691684", "19681581995"})
		fmt.Println(fmt.Sprintf("%+v", *response))

		response, _ := photosets.SetPrimaryPhoto(client, "72157656097802609", "16438207896")
		fmt.Println(fmt.Sprintf("%+v", *response))

		response, _ := photosets.OrderSets(client, []string{"72157656097802609"})
		fmt.Println(fmt.Sprintf("%+v", *response))
	*/

	response, _ := photosets.GetInfo(client, true, "72157656097802609", "")
	fmt.Println(response.Set.Title)
}

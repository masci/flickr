package main

import (
	"fmt"
	"os"

	"github.com/masci/flickr"
	"github.com/masci/flickr/test"
)

func main() {
	// retrieve Flickr credentials from env vars
	apik := os.Getenv("FLICKRGO_API_KEY")
	apisec := os.Getenv("FLICKRGO_API_SECRET")
	// do not proceed if credentials were not provided
	if apik == "" || apisec == "" {
		fmt.Fprintln(os.Stderr, "Please set FLICKRGO_API_KEY and FLICKRGO_API_SECRET env vars")
		os.Exit(1)
	}

	// create an API client with credentials
	client := flickr.NewFlickrClient(apik, apisec)

	// ask user to authorize this application

	// first, get a request token
	tok, err := flickr.GetRequestToken(client)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	// build the authorizatin URL
	url, err := flickr.GetAuthorizeUrl(client, tok)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(3)
	}

	// ask user to hit the authorization url with
	// their browser, authorize this application and coming
	// back with the confirmation token
	var oauthVerifier string
	fmt.Println("Open your browser at this url:", url)
	fmt.Print("Then, insert the code:")
	fmt.Scanln(&oauthVerifier)

	// finally, get the access token
	accessTok, err := flickr.GetAccessToken(client, tok, oauthVerifier)
	fmt.Println("Successfully retrieved OAuth token", accessTok.OAuthToken)

	// check everything works
	resp, err := test.Login(client)
	if err != nil {
		fmt.Println(err)

	} else {
		fmt.Println(resp.Status, resp.User)
	}
}

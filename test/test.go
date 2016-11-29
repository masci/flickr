// Package implementing methods: flickr.test.*
package test

import (
	"gopkg.in/masci/flickr.v2"
)

// Response type used by Login function
type LoginResponse struct {
	flickr.BasicResponse
	// the user who provided authentication infos
	User struct {
		// Flickr ID
		ID string `xml:"id,attr"`
		// Flickr Username
		Username string `xml:"username"`
	} `xml:"user"`
}

// Response type used by Echo function
type EchoResponse struct {
	flickr.BasicResponse
	// API method name, dotted notation
	Method string `xml:"method"`
	// API Key provided
	ApiKey string `xml:"api_key"`
	// API data exchange format (ex. rest)
	Format string `xml:"format"`
}

// A testing method which checks if the caller is logged in then returns their username.
// This method requires authentication with 'read' permission.
func Login(client *flickr.FlickrClient) (*LoginResponse, error) {
	client.Init()
	client.SetOAuthDefaults()
	client.Args.Set("method", "flickr.test.login")
	client.OAuthSign()

	loginResponse := &LoginResponse{}
	err := flickr.DoGet(client, loginResponse)
	return loginResponse, err
}

// Noop method
// This method requires authentication with 'read' permission.
func Null(client *flickr.FlickrClient) (*flickr.BasicResponse, error) {
	client.Init()
	client.SetOAuthDefaults()
	client.Args.Set("method", "flickr.test.null")
	client.OAuthSign()

	response := &flickr.BasicResponse{}
	err := flickr.DoGet(client, response)
	return response, err
}

// A testing method which echo's all parameters back in the response.
// This method does not require authentication.
func Echo(client *flickr.FlickrClient) (*EchoResponse, error) {
	client.EndpointUrl = flickr.API_ENDPOINT
	client.Args.Set("method", "flickr.test.echo")
	client.Args.Set("oauth_consumer_key", client.ApiKey)

	response := &EchoResponse{}
	err := flickr.DoGet(client, response)
	return response, err
}

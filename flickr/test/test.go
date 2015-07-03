// flickr.test API methods
package test

import (
	"encoding/xml"

	"github.com/masci/flickr.go/flickr"
	flickErr "github.com/masci/flickr.go/flickr/error"
)

// Response type used by Login function
type LoginResponse struct {
	flickr.FlickrResponse
	// the user who provided authentication infos
	User struct {
		XMLName xml.Name `xml:"user"`
		// Flickr ID
		ID string `xml:"id,attr"`
		// Flickr Username
		Username string `xml:"username"`
	}
}

// Response type used by Echo function
type EchoResponse struct {
	flickr.FlickrResponse
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
	client.EndpointUrl = flickr.API_ENDPOINT // TODO move to SetDefaultArgs?
	client.SetDefaultArgs()
	client.Args.Set("method", "flickr.test.login")
	client.Args.Set("oauth_token", client.OAuthToken)
	client.Args.Set("oauth_consumer_key", client.ApiKey)
	client.Sign(client.OAuthTokenSecret)

	loginResponse := &LoginResponse{}
	err := flickr.GetResponse(client, loginResponse)

	if err != nil {
		return nil, err
	}

	if loginResponse.HasErrors() {
		return loginResponse, flickErr.NewError(10)
	}

	return loginResponse, nil
}

// Noop method
// This method requires authentication with 'read' permission.
func Null(client *flickr.FlickrClient) (*flickr.FlickrResponse, error) {
	client.EndpointUrl = flickr.API_ENDPOINT
	client.SetDefaultArgs()
	client.Args.Set("method", "flickr.test.null")
	client.Args.Set("oauth_token", client.OAuthToken)
	client.Args.Set("oauth_consumer_key", client.ApiKey)
	client.Sign(client.OAuthTokenSecret)

	response := &flickr.FlickrResponse{}
	err := flickr.GetResponse(client, response)

	if err != nil {
		return nil, err
	}

	if response.HasErrors() {
		return response, flickErr.NewError(10)
	}

	return response, nil
}

// A testing method which echo's all parameters back in the response.
// This method does not require authentication.
func Echo(client *flickr.FlickrClient) (*EchoResponse, error) {
	client.EndpointUrl = flickr.API_ENDPOINT
	client.Args.Set("method", "flickr.test.echo")
	client.Args.Set("oauth_consumer_key", client.ApiKey)

	response := &EchoResponse{}
	err := flickr.GetResponse(client, response)

	if err != nil {
		return nil, err
	}

	if response.HasErrors() {
		return response, flickErr.NewError(10)
	}

	return response, nil
}

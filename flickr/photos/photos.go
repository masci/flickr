package photos

import (
	"github.com/masci/flickr.go/flickr"
)

// TODO
// This method requires authentication with 'delete' permission.
func Delete(client *flickr.FlickrClient, id string) (*flickr.BasicResponse, error) {
	client.EndpointUrl = flickr.API_ENDPOINT
	client.HTTPVerb = "POST"
	client.SetDefaultArgs()
	client.Args.Set("oauth_token", client.OAuthToken)
	client.Args.Set("oauth_consumer_key", client.ApiKey)
	client.Args.Set("method", "flickr.photos.delete")
	client.Args.Set("photo_id", id)
	client.Args.Set("api_key", client.ApiKey)
	client.Sign(client.OAuthTokenSecret)

	response := &flickr.BasicResponse{}
	err := flickr.DoPost(client, response)
	return response, err
}

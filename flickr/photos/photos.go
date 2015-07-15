package photos

import (
	"github.com/masci/flickr.go/flickr"
)

// Delete a photo from Flickr
// This method requires authentication with 'delete' permission.
func Delete(client *flickr.FlickrClient, id string) (*flickr.BasicResponse, error) {
	client.EndpointUrl = flickr.API_ENDPOINT
	client.HTTPVerb = "POST"
	client.SetOAuthDefaults()
	client.Args.Set("method", "flickr.photos.delete")
	client.Args.Set("photo_id", id)
	client.OAuthSign()

	response := &flickr.BasicResponse{}
	err := flickr.DoPost(client, response)
	return response, err
}

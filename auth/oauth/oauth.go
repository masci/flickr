// Package implementing methods: flickr.auth.oauth.*
package oauth

import (
	"github.com/masci/flickr"
)

// Response type representing data returned by CheckToken
type CheckTokenResponse struct {
	flickr.BasicResponse
	OAuth struct {
		// OAuth token
		Token string `xml:"token"`
		// String containing permissions bonded to this token
		Perms string `xml:"perms"`
		// The owner of this token
		User struct {
			// Flickr ID
			ID string `xml:"nsid,attr"`
			// Flickr Username
			Username string `xml:"username,attr"`
			// Flickr full name
			Fullname string `xml:"fullname,attr"`
		} `xml:"user"`
	} `xml:"oauth"`
}

// Returns the credentials attached to an OAuth authentication token.
// This method does not require user authentication, but the request must be api-signed.
func CheckToken(client *flickr.FlickrClient, oauthToken string) (*CheckTokenResponse, error) {
	client.EndpointUrl = flickr.API_ENDPOINT
	client.ClearArgs()
	client.Args.Set("method", "flickr.auth.oauth.checkToken")
	client.Args.Set("oauth_token", oauthToken)
	client.ApiSign()

	response := &CheckTokenResponse{}
	err := flickr.DoGet(client, response)
	return response, err
}

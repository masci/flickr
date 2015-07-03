package oauth

import (
	"encoding/xml"

	"github.com/masci/flickr.go/flickr"
	flickErr "github.com/masci/flickr.go/flickr/error"
)

// Response type representing data returned by CheckToken
type CheckTokenResponse struct {
	flickr.FlickrResponse
	OAuth struct {
		XMLName xml.Name `xml:"oauth"`
		// OAuth token
		Token string `xml:"token"`
		// String containing permissions bonded to this token
		Perms string `xml:"perms"`
		// The owner of this token
		User struct {
			XMLName xml.Name `xml:"user"`
			// Flickr ID
			ID string `xml:"nsid,attr"`
			// Flickr Username
			Username string `xml:"username,attr"`
			// Flickr full name
			Fullname string `xml:"fullname,attr"`
		}
	}
}

// Returns the credentials attached to an OAuth authentication token.
// This method does not require user authentication, but the request must be api-signed.
func CheckToken(client *flickr.FlickrClient, oauthToken string) (*CheckTokenResponse, error) {
	client.EndpointUrl = flickr.API_ENDPOINT
	client.ClearArgs()
	client.Args.Set("method", "flickr.auth.oauth.checkToken")
	client.Args.Set("oauth_token", oauthToken)
	client.Args.Set("api_key", client.ApiKey)
	client.ApiSign(client.ApiSecret)

	response := &CheckTokenResponse{}
	err := flickr.GetResponse(client, response)

	if err != nil {
		return nil, err
	}

	if response.HasErrors() {
		return response, flickErr.NewError(10)
	}

	return response, nil
}

func getAccessToken() {

}

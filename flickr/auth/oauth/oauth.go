package oauth

import (
	"encoding/xml"
	"github.com/masci/flickr.go/flickr"
	flickErr "github.com/masci/flickr.go/flickr/error"
	"io/ioutil"
)

type CheckTokenResponse struct {
	flickr.FlickrResponse
	OAuth struct {
		XMLName xml.Name `xml:"oauth"`
		Token   string   `xml:"token"`
		Perms   string   `xml:"perms"`
		User    struct {
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

func CheckToken(client *flickr.FlickrClient, oauthToken string) (*CheckTokenResponse, error) {
	client.EndpointUrl = flickr.API_ENDPOINT
	client.ClearArgs()
	client.Args.Set("method", "flickr.auth.oauth.checkToken")
	client.Args.Set("oauth_token", oauthToken)
	client.Args.Set("api_key", client.ApiKey)
	client.ApiSign(client.ApiSecret)

	res, err := client.HTTPClient.Get(client.GetUrl())
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	response := CheckTokenResponse{}
	err = xml.Unmarshal([]byte(body), &response)
	if err != nil {
		return nil, err
	}

	if response.HasErrors() {
		return &response, flickErr.NewError(10)
	}

	return &response, nil
}

func getAccessToken() {

}

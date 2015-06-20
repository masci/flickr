package test

import (
	"encoding/xml"
	"github.com/masci/flick-rsync/flickr"
	"io/ioutil"
)

type LoginResponse struct {
	flickr.FlickrResponse
	User struct {
		XMLName  xml.Name `xml:"user"`
		ID       string   `xml:"id,attr"`
		Username string   `xml:"username"`
	}
}

func Login(client *flickr.FlickrClient) (*LoginResponse, error) {
	client.EndpointUrl = flickr.API_ENDPOINT

	client.SetDefaultArgs()
	client.Args.Set("method", "flickr.test.login")
	client.Args.Set("oauth_token", client.OAuthToken)
	client.Args.Set("oauth_consumer_key", client.ApiKey)

	client.Sign(client.OAuthTokenSecret)

	res, err := client.HTTPClient.Get(client.GetUrl())
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	loginResponse := LoginResponse{}
	err = xml.Unmarshal([]byte(body), &loginResponse)
	if err != nil {
		return nil, err
	}

	// TODO parse flickr errors

	return &loginResponse, nil
}

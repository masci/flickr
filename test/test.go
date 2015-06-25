package test

import (
	"encoding/xml"
	"github.com/masci/flick-rsync/flickr"
	flickErr "github.com/masci/flick-rsync/flickr/error"
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
	client.EndpointUrl = flickr.API_ENDPOINT // TODO move to SetDefaultArgs

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

	if loginResponse.HasErrors() {
		return &loginResponse, flickErr.NewError(10)
	}

	return &loginResponse, nil
}

func Null(client *flickr.FlickrClient) (*flickr.FlickrResponse, error) {
	client.EndpointUrl = flickr.API_ENDPOINT
	client.SetDefaultArgs()
	client.Args.Set("method", "flickr.test.null")
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

	response := flickr.FlickrResponse{}

	err = xml.Unmarshal([]byte(body), &response)
	if err != nil {
		return nil, err
	}

	if response.HasErrors() {
		return &response, flickErr.NewError(10)
	}

	return &response, nil
}

type EchoResponse struct {
	flickr.FlickrResponse
	Method string `xml:"method"`
	ApiKey string `xml:"api_key"`
	Format string `xml:"format"`
}

func Echo(client *flickr.FlickrClient) (*EchoResponse, error) {
	client.EndpointUrl = flickr.API_ENDPOINT
	client.SetDefaultArgs()
	client.Args.Set("method", "flickr.test.echo")
	client.Args.Set("oauth_consumer_key", client.ApiKey)

	res, err := client.HTTPClient.Get(client.GetUrl())
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	response := EchoResponse{}

	err = xml.Unmarshal([]byte(body), &response)
	if err != nil {
		return nil, err
	}

	if response.HasErrors() {
		return &response, flickErr.NewError(10)
	}

	return &response, nil
}

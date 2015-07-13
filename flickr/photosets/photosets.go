// Package implementing methods: flickr.photosets.*
package photosets

import (
	"fmt"

	"github.com/masci/flickr.go/flickr"
)

type Photoset struct {
	Id                string `xml:"id,attr"`
	Primary           string `xml:"primary,attr"`
	Secret            string `xml:"secret,attr"`
	Server            string `xml:"server,attr"`
	Farm              string `xml:"farm,attr"`
	Photos            int    `xml:"photos,attr"`
	Videos            int    `xml:"videos,attr"`
	NeedsInterstitial bool   `xml:"needs_interstitial,attr"`
	VisCanSeeSet      bool   `xml:"visibility_can_see_set,attr"`
	CountViews        int    `xml:"count_views,attr"`
	CountComments     int    `xml:"count_comments,attr"`
	CanComment        bool   `xml:"can_comment,attr"`
	DateCreate        int    `xml:"date_create,attr"`
	DateUpdate        int    `xml:"date_update,attr"`
	Title             string `xml:"title"`
	Description       string `xml:"description"`
	Url               string `xml:"url,attr"`
}

type PhotosetsListResponse struct {
	flickr.BasicResponse
	Photosets struct {
		Page    int        `xml:"page,attr"`
		Pages   int        `xml:"pages,attr"`
		Perpage int        `xml:"perpage,attr"`
		Total   int        `xml:"total,attr"`
		Items   []Photoset `xml:"photoset"`
	} `xml:"photosets"`
}

type PhotosetResponse struct {
	flickr.BasicResponse
	Set Photoset `xml:"photoset"`
}

// Return all the photosets belonging to the caller user
// This call must be signed to get both public and private sets
func GetOwnList(client *flickr.FlickrClient) (*PhotosetsListResponse, error) {
	client.EndpointUrl = flickr.API_ENDPOINT
	client.ClearArgs()
	client.Args.Set("method", "flickr.photosets.getList")
	client.Args.Set("api_key", client.ApiKey)

	client.ApiSign(client.ApiSecret)

	response := &PhotosetsListResponse{}
	err := flickr.DoGet(client, response)
	return response, err
}

// Return the public sets belonging to the user with userId
// This method does not require authentication.
func GetList(client *flickr.FlickrClient, userId string) (*PhotosetsListResponse, error) {
	client.EndpointUrl = flickr.API_ENDPOINT
	client.ClearArgs()
	client.Args.Set("method", "flickr.photosets.getList")
	client.Args.Set("api_key", client.ApiKey)
	client.Args.Set("user_id", userId)

	response := &PhotosetsListResponse{}
	err := flickr.DoGet(client, response)
	return response, err
}

// TODO docs
// This method requires authentication with 'write' permission.
func AddPhoto(client *flickr.FlickrClient, photosetId, photoId int) (*flickr.BasicResponse, error) {
	client.EndpointUrl = flickr.API_ENDPOINT
	client.ClearArgs()
	client.Args.Set("method", "flickr.photosets.getList")
	client.Args.Set("api_key", client.ApiKey)

	client.ApiSign(client.ApiSecret)

	response := &flickr.BasicResponse{}
	err := flickr.DoPost(client, response)
	return response, err
}

// TODO docs
// This method requires authentication with 'write' permission.
func Create(client *flickr.FlickrClient, title, description string, primaryPhotoId int) (*PhotosetResponse, error) {
	client.EndpointUrl = flickr.API_ENDPOINT
	client.HTTPVerb = "POST"
	client.SetDefaultArgs()
	client.Args.Set("method", "flickr.photosets.create")
	client.Args.Set("oauth_token", client.OAuthToken)
	client.Args.Set("oauth_consumer_key", client.ApiKey)
	client.Args.Set("api_key", client.ApiKey)
	client.Args.Set("title", title)
	client.Args.Set("description", description)
	client.Args.Set("primary_photo_id", fmt.Sprintf("%d", primaryPhotoId))

	client.Sign(client.OAuthTokenSecret)

	response := &PhotosetResponse{}
	err := flickr.DoPost(client, response)
	return response, err
}

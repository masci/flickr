// Package implementing methods: flickr.photosets.*
package photosets

import (
	"github.com/masci/flickr.go/flickr"
	flickErr "github.com/masci/flickr.go/flickr/error"
)

type PhotsetsListResponse struct {
	flickr.FlickrResponse
	Photosets struct {
		Page    int `xml:"page,attr"`
		Pages   int `xml:"pages,attr"`
		Perpage int `xml:"perpage,attr"`
		Total   int `xml:"total,attr"`
		Items   []struct {
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
		} `xml:"photoset"`
	} `xml:"photosets"`
}

// Return all the photosets belonging to the caller user
// This call must be signed with full permissions to get both public and private sets
func GetOwnList(client *flickr.FlickrClient) (*PhotsetsListResponse, error) {
	client.EndpointUrl = flickr.API_ENDPOINT
	client.ClearArgs()
	client.Args.Set("method", "flickr.photosets.getList")
	client.Args.Set("api_key", client.ApiKey)

	client.ApiSign(client.ApiSecret)

	response := &PhotsetsListResponse{}
	err := flickr.GetResponse(client, response)

	if err != nil {
		return nil, err
	}

	if response.HasErrors() {
		return response, flickErr.NewError(10)
	}

	return response, nil
}

// Return the public sets belonging to the user with userId
// This method does not require authentication.
func GetList(client *flickr.FlickrClient, userId string) (*PhotsetsListResponse, error) {
	client.EndpointUrl = flickr.API_ENDPOINT
	client.ClearArgs()
	client.Args.Set("method", "flickr.photosets.getList")
	client.Args.Set("api_key", client.ApiKey)
	client.Args.Set("user_id", userId)

	response := &PhotsetsListResponse{}
	err := flickr.GetResponse(client, response)

	if err != nil {
		return nil, err
	}

	if response.HasErrors() {
		return response, flickErr.NewError(10)
	}

	return response, nil
}

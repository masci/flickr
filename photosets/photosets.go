// Package implementing methods: flickr.photosets.*
package photosets

import (
	"strconv"
	"strings"

	"github.com/masci/flickr"
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
	Owner             string `xml:"owner,attr"`
}

type Photo struct {
	Id    string `xml:"id,attr"`
	Title string `xml:"title,attr"`
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

type PhotosListResponse struct {
	flickr.BasicResponse
	Photoset struct {
		Page    int     `xml:"page,attr"`
		Pages   int     `xml:"pages,attr"`
		Perpage int     `xml:"perpage,attr"`
		Total   int     `xml:"total,attr"`
		Photos  []Photo `xml:"photo"`
	} `xml:"photoset"`
}

// Return the public sets belonging to the user with userId.
// If userId is not provided it defaults to the caller user but call needs to be authenticated.
// This method requires authentication to retrieve private sets.
func GetList(client *flickr.FlickrClient, authenticate bool, userId string, page int) (*PhotosetsListResponse, error) {
	client.Init()
	client.Args.Set("method", "flickr.photosets.getList")
	if userId != "" {
		client.Args.Set("user_id", userId)
	}
	// if not provided, flickr defaults this argument to 1
	if page > 1 {
		client.Args.Set("page", strconv.Itoa(page))
	}
	// perform authentication if requested
	if authenticate {
		client.OAuthSign()
	} else {
		client.ApiSign()
	}

	response := &PhotosetsListResponse{}
	err := flickr.DoGet(client, response)
	return response, err
}

// Add a photo to a photoset
// This method requires authentication with 'write' permission.
func AddPhoto(client *flickr.FlickrClient, photosetId, photoId string) (*flickr.BasicResponse, error) {
	client.Init()
	client.HTTPVerb = "POST"
	client.Args.Set("method", "flickr.photosets.addPhoto")
	client.Args.Set("photoset_id", photosetId)
	client.Args.Set("photo_id", photoId)

	client.OAuthSign()

	response := &flickr.BasicResponse{}
	err := flickr.DoPost(client, response)
	return response, err
}

// Create a photoset specifying its primary photo
// This method requires authentication with 'write' permission.
func Create(client *flickr.FlickrClient, title, description, primaryPhotoId string) (*PhotosetResponse, error) {
	client.Init()
	client.HTTPVerb = "POST"
	client.Args.Set("method", "flickr.photosets.create")
	client.Args.Set("title", title)
	client.Args.Set("description", description)
	client.Args.Set("primary_photo_id", primaryPhotoId)

	client.OAuthSign()

	response := &PhotosetResponse{}
	err := flickr.DoPost(client, response)
	return response, err
}

// Delete a photoset
// This method requires authentication with 'write' permission.
func Delete(client *flickr.FlickrClient, photosetId string) (*flickr.BasicResponse, error) {
	client.Init()
	client.HTTPVerb = "POST"
	client.Args.Set("method", "flickr.photosets.delete")
	client.Args.Set("photoset_id", photosetId)

	client.OAuthSign()

	response := &flickr.BasicResponse{}
	err := flickr.DoPost(client, response)
	return response, err
}

// Remove a photo from a photoset
// This method requires authentication with 'write' permission.
func RemovePhoto(client *flickr.FlickrClient, photosetId, photoId string) (*flickr.BasicResponse, error) {
	client.Init()
	client.HTTPVerb = "POST"
	client.Args.Set("method", "flickr.photosets.removePhoto")
	client.Args.Set("photoset_id", photosetId)
	client.Args.Set("photo_id", photoId)

	client.OAuthSign()

	response := &flickr.BasicResponse{}
	err := flickr.DoPost(client, response)
	return response, err
}

// Get the photos in a set
// This method requires authentication to retrieve photos from private sets
func GetPhotos(client *flickr.FlickrClient, authenticate bool, photosetId, ownerID string, page int) (*PhotosListResponse, error) {
	client.Init()
	client.Args.Set("method", "flickr.photosets.getPhotos")
	client.Args.Set("photoset_id", photosetId)
	// this argument is optional but increases query performances
	if ownerID != "" {
		client.Args.Set("user_id", ownerID)
	}
	// if not provided, flickr defaults this argument to 1
	if page > 1 {
		client.Args.Set("page", strconv.Itoa(page))
	}
	// sign the client for authentication and authorization
	if authenticate {
		client.OAuthSign()
	} else {
		client.ApiSign()
	}

	response := &PhotosListResponse{}
	err := flickr.DoGet(client, response)
	return response, err
}

// Edit set name and description
// This method requires authentication with 'write' permission.
func EditMeta(client *flickr.FlickrClient, photosetId, title, description string) (*flickr.BasicResponse, error) {
	client.Init()
	client.HTTPVerb = "POST"
	client.Args.Set("method", "flickr.photosets.editMeta")
	client.Args.Set("photoset_id", photosetId)
	client.Args.Set("title", title)
	if description != "" {
		client.Args.Set("description", description)
	}

	client.OAuthSign()

	response := &flickr.BasicResponse{}
	err := flickr.DoPost(client, response)
	return response, err
}

// Modify the photos in a photoset. Use this method to add, remove and re-order photos.
// This method requires authentication with 'write' permission.
func EditPhotos(client *flickr.FlickrClient, photosetId, primaryId string, photoIds []string) (*flickr.BasicResponse, error) {
	client.Init()
	client.HTTPVerb = "POST"
	client.Args.Set("method", "flickr.photosets.editPhotos")
	client.Args.Set("photoset_id", photosetId)
	client.Args.Set("primary_photo_id", primaryId)
	photos := strings.Join(photoIds, ",")
	client.Args.Set("photo_ids", photos)

	client.OAuthSign()

	response := &flickr.BasicResponse{}
	err := flickr.DoPost(client, response)
	return response, err
}

// Gets information about a photoset.
// This method does not require authentication unless you want to access a private set
func GetInfo(client *flickr.FlickrClient, authenticate bool, photosetId, ownerID string) (*PhotosetResponse, error) {
	client.Init()
	client.Args.Set("method", "flickr.photosets.getInfo")
	client.Args.Set("photoset_id", photosetId)
	// this argument is optional but increases query performances
	if ownerID != "" {
		client.Args.Set("user_id", ownerID)
	}

	// sign the client for authentication and authorization
	if authenticate {
		client.OAuthSign()
	} else {
		client.ApiSign()
	}

	response := &PhotosetResponse{}
	err := flickr.DoGet(client, response)
	return response, err
}

// Set the order of photosets for the calling user.
// Any set IDs not given in the list will be set to appear at the end of the list, ordered by their IDs.
// This method requires authentication with 'write' permission.
func OrderSets(client *flickr.FlickrClient, photosetIds []string) (*flickr.BasicResponse, error) {
	client.Init()
	client.HTTPVerb = "POST"
	client.Args.Set("method", "flickr.photosets.orderSets")
	sets := strings.Join(photosetIds, ",")
	client.Args.Set("photoset_ids", sets)

	client.OAuthSign()

	response := &flickr.BasicResponse{}
	err := flickr.DoPost(client, response)
	return response, err
}

// Remove multiple photos from a photoset.
// This method requires authentication with 'write' permission.
func RemovePhotos(client *flickr.FlickrClient, photosetId string, photoIds []string) (*flickr.BasicResponse, error) {
	client.Init()
	client.HTTPVerb = "POST"
	client.Args.Set("method", "flickr.photosets.removePhotos")
	client.Args.Set("photoset_id", photosetId)
	photos := strings.Join(photoIds, ",")
	client.Args.Set("photo_ids", photos)

	client.OAuthSign()

	response := &flickr.BasicResponse{}
	err := flickr.DoPost(client, response)
	return response, err
}

// Alias for EditPhotos
func ReorderPhotos(client *flickr.FlickrClient, photosetId, primaryId string, photoIds []string) (*flickr.BasicResponse, error) {
	return EditPhotos(client, photosetId, primaryId, photoIds)
}

// Set photoset primary photo
// This method requires authentication with 'write' permission.
func SetPrimaryPhoto(client *flickr.FlickrClient, photosetId, primaryId string) (*flickr.BasicResponse, error) {
	client.Init()
	client.HTTPVerb = "POST"
	client.Args.Set("method", "flickr.photosets.setPrimaryPhoto")
	client.Args.Set("photoset_id", photosetId)
	client.Args.Set("photo_id", primaryId)

	client.OAuthSign()

	response := &flickr.BasicResponse{}
	err := flickr.DoPost(client, response)
	return response, err
}

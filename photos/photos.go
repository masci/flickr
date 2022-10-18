package photos

import (
	"strconv"
	"strings"

	"gopkg.in/masci/flickr.v2"
)

type PhotoInfo struct {
	Id           string `xml:"id,attr"`
	Secret       string `xml:"secret,attr"`
	Server       string `xml:"server,attr"`
	Farm         string `xml:"farm,attr"`
	DateUploaded string `xml:"dateuploaded,attr"`
	IsFavorite   bool   `xml:"isfavorite,attr"`
	License      string `xml:"license,attr"`
	// NOTE: one less than safety level set on upload (ie, here 0 = safe, 1 = moderate, 2 = restricted)
	//       while on upload, 1 = safe, 2 = moderate, 3 = restricted
	SafetyLevel    int    `xml:"safety_level,attr"`
	Rotation       int    `xml:"rotation,attr"`
	OriginalSecret string `xml:"originalsecret,attr"`
	OriginalFormat string `xml:"originalformat,attr"`
	Views          int    `xml:"views,attr"`
	Media          string `xml:"media,attr"`
	Title          string `xml:"title"`
	Description    string `xml:"description"`
	Visibility     struct {
		IsPublic bool `xml:"ispublic,attr"`
		IsFriend bool `xml:"isfriend,attr"`
		IsFamily bool `xml:"isfamily,attr"`
	} `xml:"visibility"`
	Dates struct {
		Posted           string `xml:"posted,attr"`
		Taken            string `xml:"taken,attr"`
		TakenGranularity string `xml:"takengranularity,attr"`
		TakenUnknown     string `xml:"takenunknown,attr"`
		LastUpdate       string `xml:"lastupdate,attr"`
	} `xml:"dates"`
	Permissions struct {
		PermComment string `xml:"permcomment,attr"`
		PermAdMeta  string `xml:"permadmeta,attr"`
	} `xml:"permissions"`
	Editability struct {
		CanComment string `xml:"cancomment,attr"`
		CanAddMeta string `xml:"canaddmeta,attr"`
	} `xml:"editability"`
	PublicEditability struct {
		CanComment string `xml:"cancomment,attr"`
		CanAddMeta string `xml:"canaddmeta,attr"`
	} `xml:"publiceditability"`
	Usage struct {
		CanDownload string `xml:"candownload,attr"`
		CanBlog     string `xml:"canblog,attr"`
		CanPrint    string `xml:"canprint,attr"`
		CanShare    string `xml:"canshare,attr"`
	} `xml:"usage"`
	Comments int   `xml:"comments"`
	Tags     []Tag `xml:"tags>tag"`
	// Notes XXX: not handled yet
	// People XXX: not handled yet
	// Urls XXX: not handled yet
}
type Tag struct {
	ID    string `xml:"id,attr"`
	Raw   string `xml:"raw,attr"`
	Value string `xml:",chardata"`
}

type PhotoInfoResponse struct {
	flickr.BasicResponse
	Photo PhotoInfo `xml:"photo"`
}
type PrivacyType int64

const (
	yes PrivacyType = 1
	no  PrivacyType = 0
)

type PhotoDownloadInfo struct {
	Label  string `xml:"label,attr"`
	Width  string `xml:"width,attr"`
	Height string `xml:"height,attr"`
	Source string `xml:"source,attr"`
	Url    string `xml:"url,attr"`
	Media  string `xml:"media,attr"`
}
type PhotoAccessInfo struct {
	flickr.BasicResponse
	Sizes []PhotoDownloadInfo `xml:"sizes>size"`
}

// GetSizes get all the downloadable link as
func GetSizes(client *flickr.FlickrClient, photoId string) (*PhotoAccessInfo, error) {

	client.Init()
	client.EndpointUrl = flickr.API_ENDPOINT
	client.HTTPVerb = "POST"

	client.Args.Set("method", "flickr.photos.getSizes")
	client.Args.Set("photo_id", photoId)
	client.OAuthSign()
	response := &PhotoAccessInfo{}
	err := flickr.DoPost(client, response)
	return response, err

}

// Set permission of a photo from flickr
// this method requires authentica with 'write' permission
func SetPerms(client *flickr.FlickrClient, id string, isPublic PrivacyType, IsFriend PrivacyType, isFamily PrivacyType) (*flickr.BasicResponse, error) {

	client.Init()
	client.EndpointUrl = flickr.API_ENDPOINT
	client.HTTPVerb = "POST"
	client.Args.Set("method", "flickr.photos.setPerms")
	client.Args.Set("photo_id", id)
	client.Args.Set("is_public", strconv.Itoa(int(isPublic)))
	client.Args.Set("is_friend", strconv.Itoa(int(IsFriend)))
	client.Args.Set("is_family", strconv.Itoa(int(isFamily)))
	client.OAuthSign()
	response := &flickr.BasicResponse{}
	err := flickr.DoPost(client, response)
	return response, err
}

// Delete a photo from Flickr
// This method requires authentication with 'delete' permission.
func Delete(client *flickr.FlickrClient, id string) (*flickr.BasicResponse, error) {
	client.Init()
	client.EndpointUrl = flickr.API_ENDPOINT
	client.HTTPVerb = "POST"
	client.Args.Set("method", "flickr.photos.delete")
	client.Args.Set("photo_id", id)
	client.OAuthSign()

	response := &flickr.BasicResponse{}
	err := flickr.DoPost(client, response)
	return response, err
}

// Get information about a Flickr photo
func GetInfo(client *flickr.FlickrClient, id string, secret string) (*PhotoInfoResponse, error) {
	client.Init()
	client.EndpointUrl = flickr.API_ENDPOINT
	client.HTTPVerb = "POST"
	client.Args.Set("method", "flickr.photos.getInfo")
	client.Args.Set("photo_id", id)
	if secret != "" {
		client.Args.Set("secret", secret)
	}
	client.OAuthSign()

	response := &PhotoInfoResponse{}
	err := flickr.DoPost(client, response)
	return response, err
}

// Set date posted and date taken on a Flickr photo
// datePosted and dateTaken are optional and may be set to ""
func SetDates(client *flickr.FlickrClient, id string, datePosted string, dateTaken string) (*flickr.BasicResponse, error) {
	client.Init()
	client.EndpointUrl = flickr.API_ENDPOINT
	client.HTTPVerb = "POST"
	client.Args.Set("method", "flickr.photos.setDates")
	client.Args.Set("photo_id", id)
	if datePosted != "" {
		client.Args.Set("date_posted", datePosted)
	}
	if dateTaken != "" {
		client.Args.Set("date_taken", dateTaken)
	}
	client.OAuthSign()

	response := &flickr.BasicResponse{}
	err := flickr.DoPost(client, response)
	return response, err
}

// AddTags add tags to an existing photo
func AddTags(client *flickr.FlickrClient, photoId string, tags []string) error {
	client.Init()
	client.EndpointUrl = flickr.API_ENDPOINT
	client.HTTPVerb = "POST"
	client.Args.Set("method", "flickr.photos.addTags")
	client.Args.Set("photo_id", photoId)
	client.Args.Set("tags", strings.Join(tags, ","))
	client.OAuthSign()
	response := &flickr.BasicResponse{}
	return flickr.DoPost(client, response)
}

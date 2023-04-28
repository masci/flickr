package people

import (
	"strconv"

	"gopkg.in/masci/flickr.v3"
)

type Photo struct {
	DateTaken      string `xml:"datetaken,attr"`
	DateUpload     string `xml:"dateupload,attr"`
	Description    string `xml:"description"`
	Geo            string `xml:"geo,attr"`
	IconServer     string `xml:"iconserver,attr"`
	Id             string `xml:"id,attr"`
	IsFamily       bool   `xml:"isfamily,attr"`
	IsFriend       bool   `xml:"isfriend,attr"`
	IsPublic       bool   `xml:"ispublic,attr"`
	LastUpdate     string `xml:"lastupdate,attr"`
	License        string `xml:"license,attr"`
	MachineTags    string `xml:"machine_tags,attr"`
	Media          string `xml:"media,attr"`
	OriginalFormat string `xml:"originalformat,attr"`
	Owner          string `xml:"owner,attr"`
	OwnerName      string `xml:"ownername,attr"`
	PathAlias      string `xml:"pathalias,attr"`
	Secret         string `xml:"secret,attr"`
	Server         string `xml:"server,attr"`
	Tags           string `xml:"tags,attr"`
	Title          string `xml:"title,attr"`
	URLC           string `xml:"url_c,attr"`  // medium 800
	URLL           string `xml:"url_l,attr"`  // large
	URLM           string `xml:"url_m,attr"`  // medium 500
	URLN           string `xml:"url_n,attr"`  // small
	URLO           string `xml:"url_o,attr"`  // URL of original size image
	URLQ           string `xml:"url_q,attr"`  // large square
	URLS           string `xml:"url_s,attr"`  // square
	URLSQ          string `xml:"url_sq,attr"` // square
	URLT           string `xml:"url_t,attr"`  // thumbnail
	URLZ           string `xml:"url_z,attr"`  // medium 640
	Views          string `xml:"views,attr"`
}

type Photos struct {
	Page    int     `xml:"page,attr"`
	Pages   int     `xml:"pages,attr"`
	PerPage int     `xml:"perpage,attr"`
	Total   int     `xml:"total,attr"`
	Photos  []Photo `xml:"photo"`
}

type GetPhotosResponse struct {
	flickr.BasicResponse
	Photos Photos `xml:"photos"`
}

type SafetyLevel int

const (
	NoSafetySpecified SafetyLevel = iota
	Safe
	Moderate
	Restricted
)

type ContentType int

const (
	NoContentTypeSpecified ContentType = iota
	PhotosOnly
	ScreenShotsOnly
	OtherOnly
	PhotosAndScreenshots
	ScreenShotsAndOther
	PhotosAndOther
	All
)

type PrivacyFilterType int

const (
	NoPrivacyFilterSpecified PrivacyFilterType = iota
	Public
	Friends
	Family
	FriendsAndFamily
	Private
)

type GetPhotosOptionalArgs struct {
	SafeSearch    SafetyLevel       // optional, set to NoneSpecified to ignore
	MinUploadDate string            // optional, set to "" to ignore. mysql datetime
	MaxUploadDate string            // optional, set to "" to ignore. mysql datetime
	MinTakenDate  string            // optional, set to "" to ignore. mysql datetime
	MaxTakenDate  string            // optional, set to "" to ignore. mysql datetime
	ContentType   ContentType       // optional, set to NoneSpecified to ignore
	PrivacyFilter PrivacyFilterType // optional, set to NoneSpecified to ignore
	Extras        string            // optional, set to "" to ignore. comma separated string.
	PerPage       int               // 0 to ignore
	Page          int               // 0 to ignore
}

func GetPhotos(client *flickr.FlickrClient,
	userId string, opts GetPhotosOptionalArgs) (*GetPhotosResponse, error) {
	client.Init()
	client.EndpointUrl = flickr.API_ENDPOINT
	client.Args.Set("method", "flickr.people.getPhotos")
	client.Args.Set("user_id", userId)
	if opts.SafeSearch != NoSafetySpecified {
		client.Args.Set("safe_search", strconv.Itoa(int(opts.SafeSearch)))
	}
	if opts.MinUploadDate != "" {
		client.Args.Set("min_upload_date", opts.MinUploadDate)
	}
	if opts.MaxUploadDate != "" {
		client.Args.Set("min_upload_date", opts.MaxUploadDate)
	}
	if opts.MinTakenDate != "" {
		client.Args.Set("min_taken_date", opts.MinTakenDate)
	}
	if opts.MaxTakenDate != "" {
		client.Args.Set("max_taken_date", opts.MaxTakenDate)
	}
	if opts.ContentType != NoContentTypeSpecified {
		client.Args.Set("content_type", strconv.Itoa(int(opts.ContentType)))
	}
	if opts.PrivacyFilter != NoPrivacyFilterSpecified {
		client.Args.Set("privacy_filter", strconv.Itoa(int(opts.PrivacyFilter)))
	}
	if opts.PerPage != 0 {
		client.Args.Set("per_page", strconv.Itoa(opts.PerPage))
	}
	if opts.Page != 0 {
		client.Args.Set("page", strconv.Itoa(opts.Page))
	}
	if opts.Extras != "" {
		client.Args.Set("extras", opts.Extras)
	}
	client.OAuthSign()

	response := &GetPhotosResponse{}
	err := flickr.DoGet(client, response)
	//	if err == nil {
	//		fmt.Println("API response:", response.Extra)
	//	} else {
	//		fmt.Println("API error:", err)
	//	}
	return response, err
}

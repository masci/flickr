package people

import (
	"fmt"
	"strconv"

	"gopkg.in/masci/flickr.v2"
)

type PhotoList struct {
	Page    int `xml:"page,attr"`
	Pages   int `xml:"pages,attr"`
	PerPage int `xml:"perpage,attr"`
	Total   int `xml:"total,attr"`
	Photo   struct {
		Id       string `xml:"id,attr"`
		Owner    string `xml:"owner,attr"`
		Secret   string `xml:"secret,attr"`
		Server   string `xml:"server,attr"`
		Farm     string `xml:"farm,attr"`
		Title    string `xml:"title,attr"`
		IsPublic bool   `xml:"ispublic,attr"`
		IsFriend bool   `xml:"isfriend,attr"`
		IsFamily bool   `xml:"isfamily,attr"`

		// if extras contains "url_o" these are populated
		UrlO    string `xml:"url_o,attr"`
		HeightO int    `xml:"height_o,attr"`
		WidthO  int    `xml:"width_o,attr"`

		Description    string `xml:"description,attr"`
		License        string `xml:"license,attr"`
		DateUpload     string `xml:"date_upload,attr"`
		DateTaken      string `xml:"date_taken,attr"`
		OwnerName      string `xml:"owner_name,attr"`
		IconServer     string `xml:"icon_server,attr"`
		OriginalFormat string `xml:"original_format",attr`
		LastUpdate     string `xml:"last_udpate",attr`

		// Geo - these attributes are provided when extras contains "geo"
		Latitude  string `xml:"latitude,attr"`
		Longitude string `xml:"longitude,attr"`
		Accuracy  string `xml:"accuracy,attr"`
		Context   string `xml:"context,attr"`

		// Tags - contains space-separated lists
		Tags        string `xml:"tags,attr"`
		MachineTags string `xml:"machine_tags,attr"`

		// Original Dimensions - these attributes are provided
		// when extras contains "o_dims"
		OWidth  int `xml:"o_width,attr"`
		OHeight int `xml:"o_height,attr"`

		Views     int    `xml:"views,attr"`
		Media     string `xml:"media,attr"`
		PathAlias string `xml:"path_alias,attr"`

		// Square Urls - these attributes are provided when
		// extras contains "url_sq"
		UrlSq    string `xml:"url_sq,attr"`
		HeightSq int    `xml:"height_sq,attr"`
		WidthSq  int    `xml:"width_sq,attr"`

		// Thumbnail Urls - these attributes are provided
		// when extras contains "url_t"
		UrlT    string `xml:"url_t,attr"`
		HeightT int    `xml:"height_t,attr"`
		WidthT  int    `xml:"width_t,attr"`

		// Q Urls - these attributes are provided when
		// extras contains "url_s"
		UrlS    string `xml:"url_s,attr"`
		HeightS int    `xml:"height_s,attr"`
		WidthS  int    `xml:"width_s,attr"`

		// M Urls - these attributes are provided when
		// extras contains "url_m"
		UrlM    string `xml:"url_m,attr"`
		HeightM int    `xml:"height_m,attr"`
		WidthM  int    `xml:"width_m,attr"`

		// N Urls - these attributes are provided when
		// extras contains "url_n"
		UrlN    string `xml:"url_n,attr"`
		HeightN int    `xml:"height_n,attr"`
		WidthN  int    `xml:"width_n,attr"`

		// Z Urls - these attributes are provided when
		// extras contains "url_z"
		UrlZ    string `xml:"url_z,attr"`
		HeightZ int    `xml:"height_z,attr"`
		WidthZ  int    `xml:"width_z,attr"`

		// C Urls - these attributes are provided when
		// extras contains "url_c"
		UrlC    string `xml:"url_c,attr"`
		HeightC int    `xml:"height_c,attr"`
		WidthC  int    `xml:"width_c,attr"`

		// L Urls - these attributes are provided when
		// extras contains "url_l"
		UrlL    string `xml:"url_l,attr"`
		HeightL int    `xml:"height_l,attr"`
		WidthL  int    `xml:"width_l,attr"`
	}
}

type PhotoListResponse struct {
	flickr.BasicResponse
	Photos PhotoList `xml:"photos"`
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
	userId string, opts GetPhotosOptionalArgs) (*PhotoListResponse, error) {
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
	fmt.Println("client", client)

	response := &PhotoListResponse{}
	err := flickr.DoGet(client, response)
	//	if err == nil {
	//		fmt.Println("API response:", response.Extra)
	//	} else {
	//		fmt.Println("API error:", err)
	//	}
	return response, err
}

package groups

import (
	"strconv"

	"gopkg.in/masci/flickr.v2"
)

type ThrottleInfo struct {
	Text      string `xml:",chardata"`
	Count     string `xml:"count,attr"`
	Mode      string `xml:"mode,attr"`
	Remaining string `xml:"remaining,attr"`
}
type RestrictionsInfo struct {
	Text         string `xml:",chardata"`
	PhotosOk     string `xml:"photos_ok,attr"`
	VideosOk     string `xml:"videos_ok,attr"`
	ImagesOk     string `xml:"images_ok,attr"`
	ScreensOk    string `xml:"screens_ok,attr"`
	ArtOk        string `xml:"art_ok,attr"`
	VirtualOk    string `xml:"virtual_ok,attr"`
	SafeOk       string `xml:"safe_ok,attr"`
	ModerateOk   string `xml:"moderate_ok,attr"`
	RestrictedOk string `xml:"restricted_ok,attr"`
	HasGeo       string `xml:"has_geo,attr"`
}

type Group struct {
	Text         string `xml:",chardata"`
	Nsid         string `xml:"nsid,attr"`
	ID           string `xml:"id,attr"`
	Name         string `xml:"name,attr"`
	Member       string `xml:"member,attr"`
	Moderator    string `xml:"moderator,attr"`
	Admin        string `xml:"admin,attr"`
	Privacy      string `xml:"privacy,attr"`
	Photos       string `xml:"photos,attr"`
	Iconserver   string `xml:"iconserver,attr"`
	Iconfarm     string `xml:"iconfarm,attr"`
	MemberCount  string `xml:"member_count,attr"`
	TopicCount   string `xml:"topic_count,attr"`
	PoolCount    string `xml:"pool_count,attr"`
	Restrictions RestrictionsInfo
	Throttle     ThrottleInfo
}

type GroupInfoResponse struct {
	flickr.BasicResponse
	Group struct {
		ID          string           `xml:"id,attr"`
		Throttle    ThrottleInfo     `xml:"throttle"`
		Restriction RestrictionsInfo `xml:"restrictions"`
	} `xml:"group"`
}
type GetGroupsResponse struct {
	flickr.BasicResponse
	Groups []Group `xml:"groups>group"`
}

func GetInfo(client *flickr.FlickrClient, groupId string) (*GroupInfoResponse, error) {
	client.Init()
	client.EndpointUrl = flickr.API_ENDPOINT
	client.HTTPVerb = "POST"
	client.Args.Set("method", "flickr.groups.getInfo")
	client.Args.Set("group_id", groupId)
	client.OAuthSign()
	response := &GroupInfoResponse{}
	err := flickr.DoPost(client, response)
	return response, err

}

// GetGroups Get all the groups for current user, ,currently it supports only fetching the first 400 groups
func GetGroups(client *flickr.FlickrClient, page int, perPage int) (*GetGroupsResponse, error) {
	// TODO impliment pagination
	client.Init()
	client.EndpointUrl = flickr.API_ENDPOINT
	client.HTTPVerb = "POST"

	client.Args.Set("method", "flickr.groups.pools.getGroups")
	if page > 0 {
		client.Args.Set("page", strconv.Itoa(page))
	}
	if page > 0 {
		client.Args.Set("per_page", strconv.Itoa(perPage))
	}
	client.OAuthSign()
	response := &GetGroupsResponse{}
	err := flickr.DoPost(client, response)
	return response, err
}

// AddPhoto  Add a photo to a particular group.
func AddPhoto(client *flickr.FlickrClient, groupId, photoId string) (*flickr.BasicResponse, error) {
	client.Init()
	client.EndpointUrl = flickr.API_ENDPOINT
	client.HTTPVerb = "POST"
	client.Args.Set("method", "flickr.groups.pools.add")
	client.Args.Set("photo_id", photoId)
	client.Args.Set("group_id", groupId)
	client.OAuthSign()
	response := &flickr.BasicResponse{}
	err := flickr.DoPost(client, response)
	return response, err
}

// CanAddPhotos verify if the use can add more photos to the group or not
func (group *GroupInfoResponse) CanAddPhotos() bool {
	// TODO also validate the Throttle capacity for the day/week/month
	val, err := strconv.Atoi(group.Group.Throttle.Remaining)
	if err == nil {
		return val > 0
	}
	return false
}

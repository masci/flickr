package people

import (
	"testing"

	"gopkg.in/masci/flickr.v3"
)

var (
	body = `<?xml version="1.0" encoding="utf-8" ?>
		<rsp stat="ok">
		<photos page="2" pages="89" perpage="10" total="881">
		<photo id="2636" owner="47058503995@N01"
			secret="a123456" server="2" title="test_04"
			ispublic="1" isfriend="0" isfamily="0" />
		<photo id="2635" owner="47058503995@N01"
			secret="b123456" server="2" title="test_03"
			ispublic="0" isfriend="1" isfamily="1" />
		<photo id="2633" owner="47058503995@N01"
			secret="c123456" server="2" title="test_01"
			ispublic="1" isfriend="0" isfamily="0" />
		<photo id="2610" owner="12037949754@N01"
			secret="d123456" server="2" title="00_tall"
			ispublic="1" isfriend="0" isfamily="0" />
	</photos>
		</rsp>`
)

func TestGetPhotos(t *testing.T) {
	fclient := flickr.GetTestClient()
	server, client := flickr.FlickrMock(200, body, "text/xml")
	defer server.Close()
	fclient.HTTPClient = client

	resp, err := GetPhotos(fclient, "123456@N00", GetPhotosOptionalArgs{})
	flickr.Expect(t, err, nil)
	flickr.Expect(t, resp.Photos.Page, 2)
	flickr.Expect(t, resp.Photos.Pages, 89)
	flickr.Expect(t, resp.Photos.PerPage, 10)
	flickr.Expect(t, resp.Photos.Total, 881)
	flickr.Expect(t, len(resp.Photos.Photos), 4)
	flickr.Expect(t, resp.Photos.Photos[0], Photo{
		Id:       "2636",
		Owner:    "47058503995@N01",
		Secret:   "a123456",
		Server:   "2",
		Title:    "test_04",
		IsPublic: true,
		IsFriend: false,
		IsFamily: false,
	})
	flickr.Expect(t, resp.Photos.Photos[3], Photo{
		Id:       "2610",
		Owner:    "12037949754@N01",
		Secret:   "d123456",
		Server:   "2",
		Title:    "00_tall",
		IsPublic: true,
		IsFriend: false,
		IsFamily: false,
	})
}

func TestGetPhotosExtras(t *testing.T) {
	body := `<?xml version="1.0" encoding="utf-8"?>
	<rsp stat="ok">
		<photos page="1" pages="5923" perpage="1" total="11845">
			<photo id="52840975319" owner="123456@N00" secret="redacted1" server="65535" farm="66"
				title="" ispublic="1" isfriend="0" isfamily="0" license="0" o_width="4032"
				o_height="3024" dateupload="1682290477" lastupdate="1682290500"
				datetaken="2023-04-23 18:54:13" datetakengranularity="0" datetakenunknown="0"
				ownername="Testy McTestFace" iconserver="65535" iconfarm="66" views="0" tags="unittesttags"
				machine_tags="unittestmachinetags" geo="unittestgeo" originalsecret="redactedoriginalsecret" originalformat="jpg"
				latitude="0" longitude="0" accuracy="0" context="0" media="photo" media_status="ready"
				url_sq="https://live.staticflickr.com/65535/52840975319_16de78c2d0_s.jpg" height_sq="75"
				width_sq="75" url_t="https://live.staticflickr.com/65535/52840975319_16de78c2d0_t.jpg"
				height_t="75" width_t="100"
				url_s="https://live.staticflickr.com/65535/52840975319_16de78c2d0_m.jpg" height_s="180"
				width_s="240" url_q="https://live.staticflickr.com/65535/52840975319_16de78c2d0_q.jpg"
				height_q="150" width_q="150"
				url_m="https://live.staticflickr.com/65535/52840975319_16de78c2d0.jpg" height_m="375"
				width_m="500" url_n="https://live.staticflickr.com/65535/52840975319_16de78c2d0_n.jpg"
				height_n="240" width_n="320"
				url_z="https://live.staticflickr.com/65535/52840975319_16de78c2d0_z.jpg" height_z="480"
				width_z="640" url_c="https://live.staticflickr.com/65535/52840975319_16de78c2d0_c.jpg"
				height_c="600" width_c="800"
				url_l="https://live.staticflickr.com/65535/52840975319_16de78c2d0_b.jpg" height_l="768"
				width_l="1024" url_o="https://live.staticflickr.com/65535/52840975319_02c20aba2d_o.jpg"
				height_o="3024" width_o="4032" pathalias="somealias">
				<description>New skeletyl keyboard</description>
			</photo>
		</photos>
	</rsp>
	`
	fclient := flickr.GetTestClient()
	server, client := flickr.FlickrMock(200, body, "text/xml")
	defer server.Close()
	fclient.HTTPClient = client

	resp, err := GetPhotos(fclient, "123456@N00", GetPhotosOptionalArgs{
		Extras: "description, license, date_upload, date_taken, owner_name, icon_server, original_format, last_update, geo, tags, machine_tags, views, media, path_alias, url_sq, url_t, url_s, url_q, url_m, url_n, url_z, url_c, url_l, url_o",
	})
	flickr.Expect(t, err, nil)
	flickr.Expect(t, resp.Photos.Page, 1)
	flickr.Expect(t, resp.Photos.Pages, 5923)
	flickr.Expect(t, resp.Photos.PerPage, 1)
	flickr.Expect(t, resp.Photos.Total, 11845)
	flickr.Expect(t, len(resp.Photos.Photos), 1)
	flickr.Expect(t, resp.Photos.Photos[0], Photo{
		DateTaken:      "2023-04-23 18:54:13",
		DateUpload:     "1682290477",
		Description:    "New skeletyl keyboard",
		Geo:            "unittestgeo",
		IconServer:     "65535",
		Id:             "52840975319",
		IsFamily:       false,
		IsFriend:       false,
		IsPublic:       true,
		LastUpdate:     "1682290500",
		License:        "0",
		MachineTags:    "unittestmachinetags",
		Media:          "photo",
		OriginalFormat: "jpg",
		Owner:          "123456@N00",
		OwnerName:      "Testy McTestFace",
		PathAlias:      "somealias",
		Secret:         "redacted1",
		Server:         "65535",
		Tags:           "unittesttags",
		Title:          "",
		URLC:           "https://live.staticflickr.com/65535/52840975319_16de78c2d0_c.jpg",
		URLL:           "https://live.staticflickr.com/65535/52840975319_16de78c2d0_b.jpg",
		URLM:           "https://live.staticflickr.com/65535/52840975319_16de78c2d0.jpg",
		URLN:           "https://live.staticflickr.com/65535/52840975319_16de78c2d0_n.jpg",
		URLO:           "https://live.staticflickr.com/65535/52840975319_02c20aba2d_o.jpg",
		URLQ:           "https://live.staticflickr.com/65535/52840975319_16de78c2d0_q.jpg",
		URLS:           "https://live.staticflickr.com/65535/52840975319_16de78c2d0_m.jpg",
		URLSQ:          "https://live.staticflickr.com/65535/52840975319_16de78c2d0_s.jpg",
		URLT:           "https://live.staticflickr.com/65535/52840975319_16de78c2d0_t.jpg",
		URLZ:           "https://live.staticflickr.com/65535/52840975319_16de78c2d0_z.jpg",
		Views:          "0",
	})
}

package people

import (
	"testing"

	"gopkg.in/masci/flickr.v2"
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

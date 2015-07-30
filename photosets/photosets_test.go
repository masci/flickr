package photosets

import (
	"testing"

	"github.com/masci/flickr"
	flickErr "github.com/masci/flickr/error"
)

var (
	body = `<?xml version="1.0" encoding="utf-8" ?>
		<rsp stat="ok">
			<photosets page="1" pages="1" perpage="10" total="2">
				<photoset id="1234567890" primary="123456" secret="abcdef" server="1234" farm="99" photos="1" videos="3" needs_interstitial="0" visibility_can_see_set="1" count_views="999" count_comments="777" can_comment="0" date_create="1361132046" date_update="1376079704">
					<title>A photoset</title>
					<description />
				</photoset>
				<photoset id="1234567890" primary="123456" secret="abcdef" server="1234" farm="1" photos="5" videos="0" needs_interstitial="0" visibility_can_see_set="1" count_views="17" count_comments="0" can_comment="0" date_create="1135438501" date_update="1375623695">
					<title>Portraits</title>
					<description>Another cool photosets with some pics inside</description>
				</photoset>
			</photosets>
		</rsp>`

	bodyKo = `<?xml version="1.0" encoding="utf-8" ?>
		<rsp stat="fail">
		  <err code="1" msg="User not found" />
		</rsp>`
)

func TestGetList(t *testing.T) {
	fclient := flickr.GetTestClient()
	server, client := flickr.FlickrMock(200, body, "text/xml")
	defer server.Close()
	fclient.HTTPClient = client

	resp, err := GetList(fclient, true, "123456@N00", 2)
	flickr.Expect(t, err, nil)
	flickr.Expect(t, resp.Photosets.Page, 1)
	flickr.Expect(t, resp.Photosets.Pages, 1)
	flickr.Expect(t, resp.Photosets.Perpage, 10)
	flickr.Expect(t, resp.Photosets.Total, 2)
	flickr.Expect(t, len(resp.Photosets.Items), 2)

	set1 := resp.Photosets.Items[0]
	flickr.Expect(t, set1.Id, "1234567890")
	flickr.Expect(t, set1.Primary, "123456")
	flickr.Expect(t, set1.Secret, "abcdef")
	flickr.Expect(t, set1.Server, "1234")
	flickr.Expect(t, set1.Farm, "99")
	flickr.Expect(t, set1.Photos, 1)
	flickr.Expect(t, set1.Videos, 3)
	flickr.Expect(t, set1.NeedsInterstitial, false)
	flickr.Expect(t, set1.VisCanSeeSet, true)
	flickr.Expect(t, set1.CountViews, 999)
	flickr.Expect(t, set1.CountComments, 777)
	flickr.Expect(t, set1.CanComment, false)
	flickr.Expect(t, set1.DateCreate, 1361132046)
	flickr.Expect(t, set1.DateUpdate, 1376079704)
	flickr.Expect(t, set1.Title, "A photoset")
	flickr.Expect(t, set1.Description, "")

	set2 := resp.Photosets.Items[1]
	flickr.Expect(t, set2.Description, "Another cool photosets with some pics inside")

	params := []string{"user_id", "page"}
	flickr.AssertParamsInBody(t, fclient, params)

	server, client = flickr.FlickrMock(200, bodyKo, "text/xml")
	defer server.Close()
	fclient.HTTPClient = client
	resp, err = GetList(fclient, false, "", 1)

	if err == nil {
		t.Error("Unexpected nil error")
		t.FailNow()
	}

	ee, ok := err.(*flickErr.Error)
	if !ok {
		t.Error("err is not a flickErr.Error!")
	}

	flickr.Expect(t, ee.ErrorCode, 10)
	flickr.Expect(t, resp.HasErrors(), true)
	flickr.Expect(t, resp.ErrorCode(), 1)
}

func TestAddPhoto(t *testing.T) {
	fclient := flickr.GetTestClient()
	server, client := flickr.FlickrMock(200, `<rsp stat="ok"></rsp>`, "text/xml")
	defer server.Close()
	fclient.HTTPClient = client

	_, err := AddPhoto(fclient, "123456", "123")
	flickr.Expect(t, err, nil)

	server, client = flickr.FlickrMock(200, `<rsp stat="fail"></rsp>`, "text/xml")
	defer server.Close()
	fclient.HTTPClient = client

	resp, err := AddPhoto(fclient, "123456", "123")
	_, ok := err.(*flickErr.Error)
	flickr.Expect(t, ok, true)
	flickr.Expect(t, resp.HasErrors(), true)

	// check params, reset Flickr client to dismiss mocked responses
	fclient = flickr.GetTestClient()
	AddPhoto(fclient, "123456", "123")
	params := []string{"photoset_id", "photo_id"}
	flickr.AssertParamsInBody(t, fclient, params)
}

func TestCreate(t *testing.T) {
	fclient := flickr.GetTestClient()
	server, client := flickr.FlickrMock(200, `<rsp stat="ok"></rsp>`, "text/xml")
	defer server.Close()
	fclient.HTTPClient = client

	_, err := Create(fclient, "title", "desc", "123456")
	flickr.Expect(t, err, nil)

	server, client = flickr.FlickrMock(200, `<rsp stat="fail"></rsp>`, "text/xml")
	defer server.Close()
	fclient.HTTPClient = client

	resp, err := Create(fclient, "title", "desc", "123456")
	_, ok := err.(*flickErr.Error)
	flickr.Expect(t, ok, true)
	flickr.Expect(t, resp.HasErrors(), true)

	// check params, reset Flickr client to dismiss mocked responses
	fclient = flickr.GetTestClient()
	Create(fclient, "title", "desc", "123456")
	params := []string{"title", "description", "primary_photo_id"}
	flickr.AssertParamsInBody(t, fclient, params)
}

func TestDelete(t *testing.T) {
	fclient := flickr.GetTestClient()
	server, client := flickr.FlickrMock(200, `<rsp stat="ok"></rsp>`, "text/xml")
	defer server.Close()
	fclient.HTTPClient = client

	_, err := Delete(fclient, "123456")
	flickr.Expect(t, err, nil)

	server, client = flickr.FlickrMock(200, `<rsp stat="fail"></rsp>`, "text/xml")
	defer server.Close()
	fclient.HTTPClient = client

	resp, err := Delete(fclient, "123456")
	_, ok := err.(*flickErr.Error)
	flickr.Expect(t, ok, true)
	flickr.Expect(t, resp.HasErrors(), true)

	// check params, reset Flickr client to dismiss mocked responses
	fclient = flickr.GetTestClient()
	Delete(fclient, "123456")
	params := []string{"photoset_id"}
	flickr.AssertParamsInBody(t, fclient, params)
}

func TestRemovePhoto(t *testing.T) {
	fclient := flickr.GetTestClient()
	server, client := flickr.FlickrMock(200, `<rsp stat="ok"></rsp>`, "text/xml")
	defer server.Close()
	fclient.HTTPClient = client

	_, err := RemovePhoto(fclient, "123456", "123456")
	flickr.Expect(t, err, nil)

	server, client = flickr.FlickrMock(200, `<rsp stat="fail"></rsp>`, "text/xml")
	defer server.Close()
	fclient.HTTPClient = client

	resp, err := RemovePhoto(fclient, "123456", "123456")
	_, ok := err.(*flickErr.Error)
	flickr.Expect(t, ok, true)
	flickr.Expect(t, resp.HasErrors(), true)

	// check params, reset Flickr client to dismiss mocked responses
	fclient = flickr.GetTestClient()
	RemovePhoto(fclient, "123456", "123456")
	params := []string{"photoset_id", "photo_id"}
	flickr.AssertParamsInBody(t, fclient, params)
}

func TestGetPhotos(t *testing.T) {
	rspOk := `<?xml version="1.0" encoding="utf-8" ?>
	<rsp stat="ok">
	  <photoset id="72157654991267328" primary="18497456039" owner="126545133@N08" ownername="Caleb4ever" page="1" per_page="500" perpage="500" pages="1" total="20" title="Landscape">
		<photo id="18497456039" secret="e590ac1028" server="410" farm="1" title="Heaven sent" isprimary="1" ispublic="1" isfriend="0" isfamily="0" />
		<photo id="17217350039" secret="4fbc01db5b" server="8751" farm="9" title="" isprimary="0" ispublic="1" isfriend="0" isfamily="0" />
		<photo id="16492421763" secret="5a08237214" server="8794" farm="9" title="The Green Mile with deep roots" isprimary="0" ispublic="1" isfriend="0" isfamily="0" />
	  </photoset>
	</rsp>`

	fclient := flickr.GetTestClient()
	server, client := flickr.FlickrMock(200, rspOk, "text/xml")
	defer server.Close()
	fclient.HTTPClient = client

	resp, err := GetPhotos(fclient, false, "72157654991267328", "126545133@N08", 1)
	flickr.Expect(t, err, nil)
	flickr.Expect(t, len(resp.Photoset.Photos), 3)

	server, client = flickr.FlickrMock(200, `<rsp stat="fail"><err code="1" msg="Photoset not found" /></rsp>`, "text/xml")
	defer server.Close()
	fclient.HTTPClient = client

	resp, err = GetPhotos(fclient, false, "72157654991267328", "126545133@N08", 3)
	_, ok := err.(*flickErr.Error)
	flickr.Expect(t, ok, true)
	flickr.Expect(t, resp.HasErrors(), true)

	// check params, reset Flickr client to dismiss mocked responses
	fclient = flickr.GetTestClient()
	GetPhotos(fclient, false, "72157654991267328", "126545133@N08", 3)
	params := []string{"photoset_id", "user_id", "page"}
	flickr.AssertParamsInBody(t, fclient, params)
}

func TestEditMeta(t *testing.T) {
	fclient := flickr.GetTestClient()
	server, client := flickr.FlickrMock(200, `<rsp stat="ok"></rsp>`, "text/xml")
	defer server.Close()
	fclient.HTTPClient = client

	_, err := EditMeta(fclient, "72157654991267328", "name", "long description")
	flickr.Expect(t, err, nil)

	server, client = flickr.FlickrMock(200, `<rsp stat="fail"><err code="99" msg="Insufficient permissions."/></rsp>`, "text/xml")
	defer server.Close()
	fclient.HTTPClient = client

	resp, err := EditMeta(fclient, "72157654991267328", "name", "long description")
	_, ok := err.(*flickErr.Error)
	flickr.Expect(t, ok, true)
	flickr.Expect(t, resp.HasErrors(), true)

	// check params, reset Flickr client to dismiss mocked responses
	fclient = flickr.GetTestClient()
	EditMeta(fclient, "72157654991267328", "name", "long description")
	params := []string{"photoset_id", "title", "description"}
	flickr.AssertParamsInBody(t, fclient, params)
}

func TestEditPhotos(t *testing.T) {
	fclient := flickr.GetTestClient()
	server, client := flickr.FlickrMock(200, `<rsp stat="ok"></rsp>`, "text/xml")
	defer server.Close()
	fclient.HTTPClient = client

	_, err := EditPhotos(fclient, "72157654991267328", "123456", []string{"123456", "23456"})
	flickr.Expect(t, err, nil)

	server, client = flickr.FlickrMock(200, `<rsp stat="fail"><err code="99" msg="Insufficient permissions."/></rsp>`, "text/xml")
	defer server.Close()
	fclient.HTTPClient = client

	resp, err := EditPhotos(fclient, "72157654991267328", "123456", []string{"123456", "23456"})
	_, ok := err.(*flickErr.Error)
	flickr.Expect(t, ok, true)
	flickr.Expect(t, resp.HasErrors(), true)

	// check params, reset Flickr client to dismiss mocked responses
	fclient = flickr.GetTestClient()
	EditPhotos(fclient, "72157654991267328", "123456", []string{"123456", "23456"})
	params := []string{"photoset_id", "primary_photo_id", "photo_ids"}
	flickr.AssertParamsInBody(t, fclient, params)
}

func TestRemovePhotos(t *testing.T) {
	fclient := flickr.GetTestClient()
	server, client := flickr.FlickrMock(200, `<rsp stat="ok"></rsp>`, "text/xml")
	defer server.Close()
	fclient.HTTPClient = client

	_, err := RemovePhotos(fclient, "72157654991267328", []string{"123456", "23456"})
	flickr.Expect(t, err, nil)

	server, client = flickr.FlickrMock(200, `<rsp stat="fail"><err code="99" msg="Insufficient permissions."/></rsp>`, "text/xml")
	defer server.Close()
	fclient.HTTPClient = client

	resp, err := RemovePhotos(fclient, "72157654991267328", []string{"123456", "23456"})
	_, ok := err.(*flickErr.Error)
	flickr.Expect(t, ok, true)
	flickr.Expect(t, resp.HasErrors(), true)

	// check params, reset Flickr client to dismiss mocked responses
	fclient = flickr.GetTestClient()
	RemovePhotos(fclient, "72157654991267328", []string{"123456", "23456"})
	params := []string{"photoset_id", "photo_ids"}
	flickr.AssertParamsInBody(t, fclient, params)
}

func TestSetPrimaryPhoto(t *testing.T) {
	fclient := flickr.GetTestClient()
	server, client := flickr.FlickrMock(200, `<rsp stat="ok"></rsp>`, "text/xml")
	defer server.Close()
	fclient.HTTPClient = client

	_, err := SetPrimaryPhoto(fclient, "72157654991267328", "123456")
	flickr.Expect(t, err, nil)

	server, client = flickr.FlickrMock(200, `<rsp stat="fail"><err code="99" msg="Insufficient permissions."/></rsp>`, "text/xml")
	defer server.Close()
	fclient.HTTPClient = client

	resp, err := SetPrimaryPhoto(fclient, "72157654991267328", "123456")
	_, ok := err.(*flickErr.Error)
	flickr.Expect(t, ok, true)
	flickr.Expect(t, resp.HasErrors(), true)

	// check params, reset Flickr client to dismiss mocked responses
	fclient = flickr.GetTestClient()
	SetPrimaryPhoto(fclient, "72157654991267328", "123456")
	params := []string{"photoset_id", "photo_id"}
	flickr.AssertParamsInBody(t, fclient, params)
}

func TestGetInfo(t *testing.T) {
	respBody := `<?xml version="1.0" encoding="utf-8" ?>
		<rsp stat="ok">
		<photoset id="72157656097802609" owner="23148015@N00" username="Massimiliano Pippi" primary="16438207896" secret="abababa" server="0000" farm="0" photos="2" count_views="0" count_comments="0" count_photos="2" count_videos="0" can_comment="1" date_create="1438183533" date_update="1438183843" coverphoto_server="0" coverphoto_farm="0">
		<title>FooBarBaz</title>
		<description />
		</photoset>
		</rsp>`

	fclient := flickr.GetTestClient()
	server, client := flickr.FlickrMock(200, respBody, "text/xml")
	defer server.Close()
	fclient.HTTPClient = client

	resp, err := GetInfo(fclient, true, "72157654991267328", "")
	flickr.Expect(t, err, nil)
	flickr.Expect(t, resp.Set.Id, "72157656097802609")

	server, client = flickr.FlickrMock(200, `<rsp stat="fail"><err code="1" msg="Photoset not found" /></rsp>`, "text/xml")
	defer server.Close()
	fclient.HTTPClient = client

	resp, err = GetInfo(fclient, true, "72157654991267328", "")
	_, ok := err.(*flickErr.Error)
	flickr.Expect(t, ok, true)
	flickr.Expect(t, resp.HasErrors(), true)

	// check params, reset Flickr client to dismiss mocked responses
	fclient = flickr.GetTestClient()
	GetInfo(fclient, true, "72157654991267328", "")
	params := []string{"photoset_id"}
	flickr.AssertParamsInBody(t, fclient, params)
	GetInfo(fclient, true, "72157654991267328", "uuid")
	params = append(params, "user_id")
	flickr.AssertParamsInBody(t, fclient, params)
}

func TestOrderSet(t *testing.T) {
	fclient := flickr.GetTestClient()
	server, client := flickr.FlickrMock(200, `<rsp stat="ok"></rsp>`, "text/xml")
	defer server.Close()
	fclient.HTTPClient = client

	_, err := OrderSets(fclient, []string{"72157654991267328", "123456"})
	flickr.Expect(t, err, nil)

	server, client = flickr.FlickrMock(200, `<rsp stat="fail"><err code="99" msg="Insufficient permissions."/></rsp>`, "text/xml")
	defer server.Close()
	fclient.HTTPClient = client

	resp, err := OrderSets(fclient, []string{"72157654991267328", "123456"})
	_, ok := err.(*flickErr.Error)
	flickr.Expect(t, ok, true)
	flickr.Expect(t, resp.HasErrors(), true)

	// check params, reset Flickr client to dismiss mocked responses
	fclient = flickr.GetTestClient()
	OrderSets(fclient, []string{"72157654991267328", "123456"})
	params := []string{"photoset_ids"}
	flickr.AssertParamsInBody(t, fclient, params)

}

package photosets

import (
	"testing"

	"github.com/masci/flickr.go/flickr"
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

func TestGetOwnList(t *testing.T) {
	fclient := flickr.GetTestClient()
	server, client := flickr.FlickrMock(200, body, "text/xml")
	defer server.Close()
	fclient.HTTPClient = client

	resp, err := GetOwnList(fclient)
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
}

func TestGetList(t *testing.T) {

}

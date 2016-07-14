package oauth

import (
	"testing"

	"gopkg.in/masci/flickr.v2"
	flickErr "gopkg.in/masci/flickr.v2/error"
)

func TestCheckToken(t *testing.T) {
	body := `<?xml version="1.0" encoding="utf-8" ?>
	<rsp stat="ok">
	<oauth>
		<token>12345678901234567-12abc345def67890</token>
		<perms>delete</perms>
		<user nsid="12345678@N00" username="Massimiliano Pippi" fullname="Masci" />
	</oauth>
	</rsp>`

	fclient := flickr.GetTestClient()
	server, client := flickr.FlickrMock(200, body, "text/xml")
	defer server.Close()
	fclient.HTTPClient = client

	resp, err := CheckToken(fclient, "12345678901234567-12abc345def67890")

	if err != nil {
		t.Error("Unexpected error", err)
	}

	flickr.Expect(t, resp.OAuth.Token, "12345678901234567-12abc345def67890")
	flickr.Expect(t, resp.OAuth.Perms, "delete")
	flickr.Expect(t, resp.OAuth.User.ID, "12345678@N00")
	flickr.Expect(t, resp.OAuth.User.Username, "Massimiliano Pippi")
	flickr.Expect(t, resp.OAuth.User.Fullname, "Masci")
}

func TestCheckTokenKo(t *testing.T) {
	body := `<?xml version="1.0" encoding="utf-8" ?>
	<rsp stat="fail">
	  <err code="98" msg="Invalid token" />
	</rsp>`

	fclient := flickr.GetTestClient()
	server, client := flickr.FlickrMock(200, body, "text/xml")
	defer server.Close()
	fclient.HTTPClient = client

	resp, err := CheckToken(fclient, "12345678901234567-12abc345def67890")

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
	flickr.Expect(t, resp.ErrorCode(), 98)
}

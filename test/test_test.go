package test

import (
	"github.com/masci/flick-rsync/flickr"
	flickErr "github.com/masci/flick-rsync/flickr/error"
	"testing"
)

func TestLoginKo(t *testing.T) {
	body := `<?xml version="1.0" encoding="utf-8" ?>
	<rsp stat="fail">
		<err code="98" msg="Invalid auth token" />
	</rsp>`

	fclient := flickr.GetTestClient()
	server, client := flickr.FlickrMock(200, body, "text/xml")
	defer server.Close()
	fclient.HTTPClient = client

	resp, err := Login(fclient)

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

func TestLogin(t *testing.T) {
	body := `<?xml version="1.0" encoding="utf-8" ?>
	<rsp stat="ok">
	  <user id="21156022@N00">
		<username>John Doe</username>
	  </user>
	</rsp>`

	fclient := flickr.GetTestClient()
	server, client := flickr.FlickrMock(200, body, "text/xml")
	defer server.Close()
	fclient.HTTPClient = client

	resp, err := Login(fclient)

	if err != nil {
		t.Error("Unexpected error", err)
	}

	flickr.Expect(t, resp.HasErrors(), false)
	flickr.Expect(t, resp.User.ID, "21156022@N00")
	flickr.Expect(t, resp.User.Username, "John Doe")
}

func TestNullKo(t *testing.T) {
	body := `<?xml version="1.0" encoding="utf-8" ?>
	<rsp stat="fail">
	  <err code="99" msg="Insufficient permissions. Method requires read privileges; none granted." />
	</rsp>`

	fclient := flickr.GetTestClient()
	server, client := flickr.FlickrMock(200, body, "text/xml")
	defer server.Close()
	fclient.HTTPClient = client

	resp, err := Null(fclient)

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
	flickr.Expect(t, resp.ErrorCode(), 99)
}

func TestNull(t *testing.T) {
	body := `<?xml version="1.0" encoding="utf-8" ?>
	<rsp stat="ok"></rsp>`

	fclient := flickr.GetTestClient()
	server, client := flickr.FlickrMock(200, body, "text/xml")
	defer server.Close()
	fclient.HTTPClient = client

	resp, err := Login(fclient)

	if err != nil {
		t.Error("Unexpected error", err)
	}

	flickr.Expect(t, resp.HasErrors(), false)
}

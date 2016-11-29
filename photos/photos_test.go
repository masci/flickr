package photos

import (
	"testing"

	"gopkg.in/masci/flickr.v2"
	flickErr "gopkg.in/masci/flickr.v2/error"
)

func TestDelete(t *testing.T) {
	fclient := flickr.GetTestClient()
	server, client := flickr.FlickrMock(200, `<?xml version="1.0" encoding="utf-8" ?><rsp stat="ok"></rsp>`, "")
	defer server.Close()
	fclient.HTTPClient = client
	resp, err := Delete(fclient, "123456")
	flickr.Expect(t, err, nil)
	flickr.Expect(t, resp.HasErrors(), false)
}

func TestDeleteKo(t *testing.T) {
	fclient := flickr.GetTestClient()
	server, client := flickr.FlickrMock(200, `<?xml version="1.0" encoding="utf-8" ?><rsp stat="fail"></rsp>`, "")
	defer server.Close()
	fclient.HTTPClient = client
	resp, err := Delete(fclient, "123456")
	_, ok := err.(*flickErr.Error)
	flickr.Expect(t, ok, true)
	flickr.Expect(t, resp.HasErrors(), true)
}

package flickr

import (
	"bytes"
	"testing"
)

func TestDoGet(t *testing.T) {
	bodyStr := `<?xml version="1.0" encoding="utf-8" ?><rsp stat="ok"></rsp>`

	fclient := GetTestClient()
	server, client := FlickrMock(200, bodyStr, "")
	defer server.Close()
	fclient.HTTPClient = client

	err := DoGet(fclient, &FooResponse{})
	Expect(t, err, nil)
}

func TestDoPostBody(t *testing.T) {
	bodyStr := `<?xml version="1.0" encoding="utf-8" ?><rsp stat="ok"></rsp>`

	fclient := GetTestClient()
	server, client := FlickrMock(200, bodyStr, "")
	defer server.Close()
	fclient.HTTPClient = client

	err := DoPostBody(fclient, bytes.NewBufferString("foo"), "", &FooResponse{})
	Expect(t, err, nil)
}

func TestDoPost(t *testing.T) {
	fclient := GetTestClient()
	fclient.Args.Set("fooArg", "foo way")

	params := []string{"fooArg"}
	AssertParamsInBody(t, fclient, params)
}

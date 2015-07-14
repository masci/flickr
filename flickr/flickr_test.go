package flickr

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
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
	var handler = func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		fmt.Fprintln(w, "Hello, client")
		Expect(t, strings.Contains(string(body), `Content-Disposition: form-data; name="fooArg"`), true)
		Expect(t, strings.Contains(string(body), "foo way"), true)
	}

	ts := httptest.NewServer(http.HandlerFunc(handler))
	defer ts.Close()

	fclient := GetTestClient()
	fclient.EndpointUrl = ts.URL
	fclient.Args.Set("fooArg", "foo way")

	DoPost(fclient, &FooResponse{})
}

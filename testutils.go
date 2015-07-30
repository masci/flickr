package flickr

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path"
	"reflect"
	"strings"
	"testing"
)

func Expect(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("Expected %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

// testing keys were published at http://www.wackylabs.net/2011/12/oauth-and-flickr-part-2/
func GetTestClient() *FlickrClient {
	args := url.Values{}
	args.Set("oauth_nonce", "C2F26CD5C075BA9050AD8EE90644CF29")
	args.Set("oauth_timestamp", "1316657628")
	args.Set("oauth_consumer_key", "768fe946d252b119746fda82e1599980")
	args.Set("oauth_signature_method", "HMAC-SHA1")
	args.Set("oauth_version", "1.0")
	args.Set("oauth_callback", "http://www.wackylabs.net/oauth/test")

	return &FlickrClient{
		EndpointUrl: "http://www.flickr.com/services/oauth/request_token",
		HTTPClient:  &http.Client{},
		HTTPVerb:    "GET",
		Args:        args,
		ApiSecret:   "1a3c208e172d3edc",
	}
}

// mock the Flickr API
type RewriteTransport struct {
	Transport http.RoundTripper
	URL       *url.URL
}

func (t RewriteTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.URL.Scheme = t.URL.Scheme
	req.URL.Host = t.URL.Host
	req.URL.Path = path.Join(t.URL.Path, req.URL.Path)
	rt := t.Transport
	if rt == nil {
		rt = http.DefaultTransport
	}
	return rt.RoundTrip(req)
}

func FlickrMock(code int, body string, contentType string) (*httptest.Server, *http.Client) {
	if contentType == "" {
		contentType = "text/plain;charset=UTF-8"
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(code)
		w.Header().Set("Content-Type", contentType)
		fmt.Fprintln(w, body)
	}))

	u, _ := url.Parse(server.URL)

	return server, &http.Client{Transport: RewriteTransport{URL: u}}
}

// A ReaderCloser to fake http.Response Body field
type FakeBody struct {
	content *bytes.Buffer
}

func (f *FakeBody) Read(p []byte) (n int, err error) {
	return f.content.Read(p)
}

func (f FakeBody) Close() error {
	// noop
	return nil
}

func NewFakeBody(s string) *FakeBody {
	return &FakeBody{content: bytes.NewBufferString(s)}
}

// TODO
func AssertParamsInBody(t *testing.T, client *FlickrClient, params []string) {
	var handler = func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		bodystr := string(body)
		fmt.Fprintln(w, "Hello, client")
		for _, p := range params {
			needle := fmt.Sprintf(`Content-Disposition: form-data; name="%s"`, p)
			Expect(t, strings.Contains(bodystr, needle), true)
		}
	}

	ts := httptest.NewServer(http.HandlerFunc(handler))
	defer ts.Close()

	client.EndpointUrl = ts.URL
	DoPost(client, &BasicResponse{})
}

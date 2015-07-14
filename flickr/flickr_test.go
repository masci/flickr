package flickr

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	flickErr "github.com/masci/flickr.go/flickr/error"
)

type FooResponse struct {
	BasicResponse
	Foo string `xml:"foo"`
}

func TestParseRequestToken(t *testing.T) {
	in := "oauth_callback_confirmed=true&oauth_token=72157654304937659-8eedcda57d9d57e3&oauth_token_secret=8700d234e3fc00c6"
	expected := RequestToken{true, "72157654304937659-8eedcda57d9d57e3", "8700d234e3fc00c6", ""}

	tok, err := ParseRequestToken(in)
	Expect(t, nil, err)
	Expect(t, *tok, expected)
}

func TestParseRequestTokenKo(t *testing.T) {
	in := "oauth_problem=foo"
	tok, err := ParseRequestToken(in)

	ee, ok := err.(*flickErr.Error)
	if !ok {
		t.Error("err is not a flickErr.Error!")
	}

	Expect(t, ee.ErrorCode, 20)
	Expect(t, tok.OAuthProblem, "foo")

	tok, err = ParseRequestToken("notA%%%ValidUrl")
	if err == nil {
		t.Error("Parsing an invalid URL string should rise an error")
	}
}

func TestGetRequestToken(t *testing.T) {
	fclient := GetTestClient()
	mocked_body := "oauth_callback_confirmed=true&oauth_token=72157654304937659-8eedcda57d9d57e3&oauth_token_secret=8700d234e3fc00c6"
	server, client := FlickrMock(200, mocked_body, "")
	defer server.Close()
	// use the mocked client
	fclient.HTTPClient = client

	tok, err := GetRequestToken(fclient)
	if err != nil {
		t.Error("Unexpected error:", err)
	}

	Expect(t, tok.OauthCallbackConfirmed, true)
	Expect(t, tok.OauthToken, "72157654304937659-8eedcda57d9d57e3")
	Expect(t, tok.OauthTokenSecret, "8700d234e3fc00c6")
}

func TestGetAuthorizeUrl(t *testing.T) {
	client := GetTestClient()
	tok := &RequestToken{true, "token", "token_secret", ""}
	url, err := GetAuthorizeUrl(client, tok)
	Expect(t, err, nil)
	Expect(t, url, "https://www.flickr.com/services/oauth/authorize?oauth_token=token&perms=delete")
}

func TestParseOAuthToken(t *testing.T) {
	response := "fullname=Jamal%20Fanaian" +
		"&oauth_token=72157626318069415-087bfc7b5816092c" +
		"&oauth_token_secret=a202d1f853ec69de" +
		"&user_nsid=21207597%40N07" +
		"&username=jamalfanaian"

	tok, _ := ParseOAuthToken(response)

	Expect(t, tok.OAuthToken, "72157626318069415-087bfc7b5816092c")
	Expect(t, tok.OAuthTokenSecret, "a202d1f853ec69de")
	Expect(t, tok.UserNsid, "21207597@N07")
	Expect(t, tok.Username, "jamalfanaian")
	Expect(t, tok.Fullname, "Jamal Fanaian")
}

func TestParseOAuthTokenKo(t *testing.T) {
	response := "oauth_problem=foo"
	tok, err := ParseOAuthToken(response)

	ee, ok := err.(*flickErr.Error)
	if !ok {
		t.Error("err is not a flickErr.Error!")
	}

	Expect(t, ee.ErrorCode, 30)
	Expect(t, tok.OAuthProblem, "foo")

	tok, err = ParseOAuthToken("notA%%%ValidUrl")
	if err == nil {
		t.Error("Parsing an invalid URL string should rise an error")
	}

}

func TestGetAccessToken(t *testing.T) {
	body := "fullname=Jamal%20Fanaian" +
		"&oauth_token=72157626318069415-087bfc7b5816092c" +
		"&oauth_token_secret=a202d1f853ec69de" +
		"&user_nsid=21207597%40N07" +
		"&username=jamalfanaian"
	fclient := GetTestClient()

	server, client := FlickrMock(200, body, "")
	defer server.Close()
	// use the mocked client
	fclient.HTTPClient = client

	rt := &RequestToken{true, "token", "token_secret", ""}

	_, err := GetAccessToken(fclient, rt, "fooVerifier")
	if err != nil {
		t.Error("Unexpected error:", err)
	}
}

func TestFlickrResponse(t *testing.T) {
	failure := `<?xml version="1.0" encoding="utf-8" ?>
<rsp stat="fail">
  <err code="99" msg="Insufficient permissions. Method requires read privileges; none granted." />
</rsp>
`
	resp := FooResponse{}
	err := xml.Unmarshal([]byte(failure), &resp)
	if err != nil {
		t.Error("Error unmarsshalling", failure)
	}

	Expect(t, resp.HasErrors(), true)
	Expect(t, resp.ErrorCode(), 99)
	Expect(t, resp.ErrorMsg(), "Insufficient permissions. Method requires read privileges; none granted.")

	ok := `<?xml version="1.0" encoding="utf-8" ?>
<rsp stat="ok">
  <user id="23148015@N00">
    <username>Massimiliano Pippi</username>
  </user>
  <foo>Foo!</foo>
</rsp>`

	resp = FooResponse{}
	err = xml.Unmarshal([]byte(ok), &resp)
	if err != nil {
		t.Error("Error unmarsshalling", ok)
	}

	Expect(t, resp.HasErrors(), false)
	Expect(t, resp.Foo, "Foo!")
	Expect(t, resp.ErrorCode(), 0)
	Expect(t, resp.ErrorMsg(), "")

	resp = FooResponse{}
	resp.SetErrorStatus(true)
	resp.SetErrorMsg("a message")
	resp.SetErrorCode(999)
	Expect(t, resp.HasErrors(), true)
	Expect(t, resp.ErrorMsg(), "a message")
	Expect(t, resp.ErrorCode(), 999)
	resp.SetErrorStatus(false)
	Expect(t, resp.HasErrors(), false)
}

func TestParseResponse(t *testing.T) {
	bodyStr := `<?xml version="1.0" encoding="utf-8" ?>
<rsp stat="ok">
  <user id="23148015@N00">
    <username>Massimiliano Pippi</username>
  </user>
  <foo>Foo!</foo>
</rsp>`

	flickrResp := &FooResponse{}
	response := &http.Response{}
	response.Body = NewFakeBody(bodyStr)

	err := parseApiResponse(response, flickrResp)

	Expect(t, err, nil)
	Expect(t, flickrResp.Foo, "Foo!")

	response = &http.Response{}
	response.Body = NewFakeBody("a_non_rest_format_error")

	err = parseApiResponse(response, flickrResp)
	ferr, ok := err.(*flickErr.Error)
	Expect(t, ok, true)
	Expect(t, ferr.ErrorCode, 10)

	response = &http.Response{}
	response.Body = NewFakeBody(`<?xml version="1.0" encoding="utf-8" ?><rsp stat="fail"></rsp>`)
	err = parseApiResponse(response, flickrResp)
	//ferr, ok := err.(*flickErr.Error)
	//Expect(t, ok, true)
	//Expect(t, ferr.ErrorCode, 10)
}

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

package flickr

import (
	"encoding/xml"
	flickErr "github.com/masci/flick-rsync/flickr/error"
	"testing"
)

func TestGetSigningBaseString(t *testing.T) {
	c := GetTestClient()

	ret := c.getSigningBaseString()
	expected := "GET&http%3A%2F%2Fwww.flickr.com%2Fservices%2Foauth%2Frequest_token&" +
		"oauth_callback%3Dhttp%253A%252F%252Fwww.wackylabs.net%252F" +
		"oauth%252Ftest%26oauth_consumer_key%3D768fe946d252b119746fda82e1599980%26" +
		"oauth_nonce%3DC2F26CD5C075BA9050AD8EE90644CF29%26" +
		"oauth_signature_method%3DHMAC-SHA1%26oauth_timestamp%3D1316657628%26" +
		"oauth_version%3D1.0"

	Expect(t, ret, expected)
}

func TestSign(t *testing.T) {
	c := GetTestClient()

	c.Sign("token12345secret")
	expected := "dXyfrCetFSTpzD3djSrkFhj0MIQ="
	signed := c.Args.Get("oauth_signature")
	Expect(t, signed, expected)

	// test empty token_secret
	c.Sign("")
	expected = "0fhNGlzpFNAsTme/hDfUb5HPB5U="
	signed = c.Args.Get("oauth_signature")
	Expect(t, signed, expected)
}

func TestClearArgs(t *testing.T) {
	c := GetTestClient()
	c.SetDefaultArgs()
	c.ClearArgs()
	Expect(t, len(c.Args), 0)
}

func TestGenerateNonce(t *testing.T) {
	var nonce string
	nonce = generateNonce()
	Expect(t, 8, len(nonce))
}

func TestSetDefaultArgs(t *testing.T) {
	c := GetTestClient()
	c.SetDefaultArgs()
	check := func(key string) {
		val := c.Args.Get(key)
		if val == "" {
			t.Error("Found empty string for", key)
		}
	}

	check("oauth_version")
	check("oauth_signature_method")
	check("oauth_nonce")
	check("oauth_timestamp")
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

func TestNewFlickrClient(t *testing.T) {
	tok := NewFlickrClient("apikey", "apisecret")
	Expect(t, tok.ApiKey, "apikey")
	Expect(t, tok.ApiSecret, "apisecret")
	Expect(t, tok.HTTPVerb, "GET")
	Expect(t, len(tok.Args), 0)
}

func TestNewOAuthToken(t *testing.T) {
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
	type FooResponse struct {
		FlickrResponse
		Foo string `xml:"foo"`
	}

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
}

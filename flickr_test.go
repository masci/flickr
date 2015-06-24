package flickr

import (
	"encoding/xml"
	"testing"
)

func TestGetSigningBaseString(t *testing.T) {
	c := getTestClient()

	ret := c.getSigningBaseString()
	expected := "GET&http%3A%2F%2Fwww.flickr.com%2Fservices%2Foauth%2Frequest_token&" +
		"oauth_callback%3Dhttp%253A%252F%252Fwww.wackylabs.net%252F" +
		"oauth%252Ftest%26oauth_consumer_key%3D768fe946d252b119746fda82e1599980%26" +
		"oauth_nonce%3DC2F26CD5C075BA9050AD8EE90644CF29%26" +
		"oauth_signature_method%3DHMAC-SHA1%26oauth_timestamp%3D1316657628%26" +
		"oauth_version%3D1.0"

	expect(t, ret, expected)
}

func TestSign(t *testing.T) {
	c := getTestClient()

	c.Sign("token12345secret")
	expected := "dXyfrCetFSTpzD3djSrkFhj0MIQ="
	signed := c.Args.Get("oauth_signature")
	expect(t, signed, expected)

	// test empty token_secret
	c.Sign("")
	expected = "0fhNGlzpFNAsTme/hDfUb5HPB5U="
	signed = c.Args.Get("oauth_signature")
	expect(t, signed, expected)
}

func TestGenerateNonce(t *testing.T) {
	var nonce string
	nonce = generateNonce()
	expect(t, 8, len(nonce))
}

func TestGetDefaultArgs(t *testing.T) {
	args := getDefaultArgs()
	check := func(key string) {
		val := args.Get(key)
		if val == "" {
			t.Error("Found empty string for", key)
		}
	}

	check("oauth_version")
	check("oauth_signature_method")
	check("oauth_nonce")
	check("oauth_timestamp")
}

func TestNewRequestToken(t *testing.T) {
	in := "oauth_callback_confirmed=true&oauth_token=72157654304937659-8eedcda57d9d57e3&oauth_token_secret=8700d234e3fc00c6"
	expected := RequestToken{true, "72157654304937659-8eedcda57d9d57e3", "8700d234e3fc00c6"}

	tok, err := NewRequestToken(in)
	expect(t, nil, err)
	expect(t, *tok, expected)

	tok, err = NewRequestToken("notA%%%ValidUrl")
	if err == nil {
		t.Error("Parsing an invalid URL string should rise an error")
	}
}

func TestGetRequestToken(t *testing.T) {
	fclient := getTestClient()
	mocked_body := "oauth_callback_confirmed=true&oauth_token=72157654304937659-8eedcda57d9d57e3&oauth_token_secret=8700d234e3fc00c6"
	server, client := flickrMock(200, mocked_body, "")
	defer server.Close()
	// use the mocked client
	fclient.HTTPClient = client

	tok, err := GetRequestToken(fclient)
	if err != nil {
		t.Error("Unexpected error:", err)
	}

	expect(t, tok.OauthCallbackConfirmed, true)
	expect(t, tok.OauthToken, "72157654304937659-8eedcda57d9d57e3")
	expect(t, tok.OauthTokenSecret, "8700d234e3fc00c6")
}

func TestGetAuthorizeUrl(t *testing.T) {
	client := getTestClient()
	tok := &RequestToken{true, "token", "token_secret"}
	url, err := GetAuthorizeUrl(client, tok)
	expect(t, err, nil)
	expect(t, url, "https://www.flickr.com/services/oauth/authorize?oauth_token=token&perms=delete")
}

func TestNewFlickrClient(t *testing.T) {
	tok := NewFlickrClient("apikey", "apisecret")
	expect(t, tok.ApiKey, "apikey")
	expect(t, tok.ApiSecret, "apisecret")
	expect(t, tok.HTTPVerb, "GET")
	expect(t, len(tok.Args), 0)
}

func TestNewOAuthToken(t *testing.T) {
	response := "fullname=Jamal%20Fanaian" +
		"&oauth_token=72157626318069415-087bfc7b5816092c" +
		"&oauth_token_secret=a202d1f853ec69de" +
		"&user_nsid=21207597%40N07" +
		"&username=jamalfanaian"

	tok, _ := NewOAuthToken(response)

	expect(t, tok.OAuthToken, "72157626318069415-087bfc7b5816092c")
	expect(t, tok.OAuthTokenSecret, "a202d1f853ec69de")
	expect(t, tok.UserNsid, "21207597@N07")
	expect(t, tok.Username, "jamalfanaian")
	expect(t, tok.Fullname, "Jamal Fanaian")
}

func TestGetAccessToken(t *testing.T) {
	body := "fullname=Jamal%20Fanaian" +
		"&oauth_token=72157626318069415-087bfc7b5816092c" +
		"&oauth_token_secret=a202d1f853ec69de" +
		"&user_nsid=21207597%40N07" +
		"&username=jamalfanaian"
	fclient := getTestClient()

	server, client := flickrMock(200, body, "")
	defer server.Close()
	// use the mocked client
	fclient.HTTPClient = client

	rt := &RequestToken{true, "token", "token_secret"}

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

	expect(t, resp.HasErrors(), true)
	expect(t, resp.ErrorCode(), 99)
	expect(t, resp.ErrorMsg(), "Insufficient permissions. Method requires read privileges; none granted.")

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

	expect(t, resp.HasErrors(), false)
	expect(t, resp.Foo, "Foo!")
	expect(t, resp.ErrorCode(), 0)
	expect(t, resp.ErrorMsg(), "")
}

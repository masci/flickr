package flickr

import (
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
	c.SetOAuthDefaults()
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
	c.SetOAuthDefaults()
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

func TestNewFlickrClient(t *testing.T) {
	tok := NewFlickrClient("apikey", "apisecret")
	Expect(t, tok.ApiKey, "apikey")
	Expect(t, tok.ApiSecret, "apisecret")
	Expect(t, tok.HTTPVerb, "GET")
	Expect(t, len(tok.Args), 0)
	Expect(t, tok.Id, "")
}

func TestApiSign(t *testing.T) {
	client := NewFlickrClient("1234567890", "SECRET")
	client.Args.Set("foo", "1")
	client.Args.Set("bar", "2")
	client.Args.Set("baz", "3")

	client.ApiSign()

	Expect(t, client.Args.Get("api_sig"), "0a55ae496d1db08f39deb5d894ae3849")
}

func TestInit(t *testing.T) {
	client := GetTestClient()
	client.Args.Set("foo", "bar")
	client.EndpointUrl = ""
	client.Init()
	Expect(t, len(client.Args), 0)
	Expect(t, client.EndpointUrl != "", true)
}

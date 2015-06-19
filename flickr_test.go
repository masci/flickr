package flickr

import (
	"fmt"
	"net/url"
	"testing"
)

// testing keys were published at http://www.wackylabs.net/2011/12/oauth-and-flickr-part-2/
func getTestClient() *FlickrClient {
	args := url.Values{}
	args.Set("oauth_nonce", "C2F26CD5C075BA9050AD8EE90644CF29")
	args.Set("oauth_timestamp", "1316657628")
	args.Set("oauth_consumer_key", "768fe946d252b119746fda82e1599980")
	args.Set("oauth_signature_method", "HMAC-SHA1")
	args.Set("oauth_version", "1.0")
	args.Set("oauth_callback", "http://www.wackylabs.net/oauth/test")

	return &FlickrClient{
		EndpointUrl: "http://www.flickr.com/services/oauth/request_token",
		Method:      "GET",
		Args:        args,
		ApiSecret:   "1a3c208e172d3edc",
	}
}

func TestGetSigningBaseString(t *testing.T) {
	c := getTestClient()

	ret := getSigningBaseString(c)
	expected := "GET&http%3A%2F%2Fwww.flickr.com%2Fservices%2Foauth%2Frequest_token&" +
		"oauth_callback%3Dhttp%253A%252F%252Fwww.wackylabs.net%252F" +
		"oauth%252Ftest%26oauth_consumer_key%3D768fe946d252b119746fda82e1599980%26" +
		"oauth_nonce%3DC2F26CD5C075BA9050AD8EE90644CF29%26" +
		"oauth_signature_method%3DHMAC-SHA1%26oauth_timestamp%3D1316657628%26" +
		"oauth_version%3D1.0"

	if ret != expected {
		t.Error("Expected", expected, "found", ret)
	}
}

func TestSign(t *testing.T) {
	c := getTestClient()

	c.Sign("token12345secret")
	expected := "dXyfrCetFSTpzD3djSrkFhj0MIQ="
	signed := c.Args.Get("oauth_signature")

	if signed != expected {
		t.Error("Expected", expected, "found", signed)
	}
	fmt.Println(signed)

	// test empty token_secret
	c.Sign("")
	expected = "0fhNGlzpFNAsTme/hDfUb5HPB5U="
	signed = c.Args.Get("oauth_signature")

	if signed != expected {
		t.Error("Expected", expected, "found", signed)
	}
}

func TestGenerateNonce(t *testing.T) {
	var nonce string
	nonce = generateNonce()
	if len(nonce) != 8 {
		t.Error("Expected string with length of 8, found", nonce)
	}
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

func TestParseRequestToken(t *testing.T) {
	tok := RequestToken{}
	in := "oauth_callback_confirmed=true&oauth_token=72157654304937659-8eedcda57d9d57e3&oauth_token_secret=8700d234e3fc00c6"
	expected := RequestToken{true, "72157654304937659-8eedcda57d9d57e3", "8700d234e3fc00c6"}

	err := tok.Parse(in)
	if err != nil {
		t.Error("Error:", err)
	}

	if tok != expected {
		t.Error("Expected", expected, "found", tok)
	}

	err = tok.Parse("notA%%%ValidUrl")
	if err == nil {
		t.Error("Parsing an invalid URL string should rise an error")
	}
}

func TestGetRequestToken(t *testing.T) {
	//GetRequestToken("a70c23170443bac2c189b92fc6439ef0", "82c542eaba4f56c9")
}

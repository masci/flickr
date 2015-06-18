package flickr

import (
	"net/url"
	"testing"
)

// testing keys were published at http://www.wackylabs.net/2011/12/oauth-and-flickr-part-2/
func getTestRequest() *Request {
	args := url.Values{}
	args.Add("oauth_nonce", "C2F26CD5C075BA9050AD8EE90644CF29")
	args.Add("oauth_timestamp", "1316657628")
	args.Add("oauth_consumer_key", "768fe946d252b119746fda82e1599980")
	args.Add("oauth_signature_method", "HMAC-SHA1")
	args.Add("oauth_version", "1.0")
	args.Add("oauth_callback", "http://www.wackylabs.net/oauth/test")

	return NewRequest("http://www.flickr.com/services/oauth/request_token", "GET", args)
}

func TestGetSigningBaseString(t *testing.T) {
	r := getTestRequest()

	ret := getSigningBaseString(r)
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
	r := getTestRequest()

	signed := Sign(r, "1a3c208e172d3edc", "token12345secret")
	expected := "dXyfrCetFSTpzD3djSrkFhj0MIQ="
	if signed != expected {
		t.Error("Expected", expected, "found", signed)
	}

	// test empty token_secret
	signed = Sign(r, "1a3c208e172d3edc", "")
	expected = "0fhNGlzpFNAsTme/hDfUb5HPB5U="
	if signed != expected {
		t.Error("Expected", expected, "found", signed)
	}
}

func TestParseRequestToken(t *testing.T) {
	in := "oauth_callback_confirmed=true&oauth_token=72157654304937659-8eedcda57d9d57e3&oauth_token_secret=8700d234e3fc00c6"
	expected := RequestToken{true, "72157654304937659-8eedcda57d9d57e3", "8700d234e3fc00c6"}

	rt, err := parseRequestToken(in)
	if err != nil {
		t.Error("Error:", err)
	}

	if *rt != expected {
		t.Error("Expected", expected, "found", rt)
	}
}

func TestGetRequestToken(t *testing.T) {
	//GetRequestToken("a70c23170443bac2c189b92fc6439ef0", "82c542eaba4f56c9")
}

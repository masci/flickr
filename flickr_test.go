package flickr

import (
	"net/url"
	"testing"
)

func TestSign(t *testing.T) {
	args := url.Values{}
	args.Add("oauth_nonce", "89601180")
	args.Add("oauth_timestamp", "1305583298")
	args.Add("oauth_consumer_key", "653e7a6ecc1d528c516cc8f92cf98611")
	args.Add("oauth_signature_method", "HMAC-SHA1")
	args.Add("oauth_version", "1.0")
	args.Add("oauth_callback", "http%3A%2F%2Fwww.example.com")

	r := NewRequest("https://www.flickr.com/services/oauth/request_token", "GET", args)

	signed := Sign(r, "api1234567secret", "token12345secret")
	expected := "+G489dWwcLcsN7yn3BrKG9DlWTk="
	if signed != expected {
		t.Error("Expected", expected, "found", signed)
	}

	// test empty token_secret
	signed = Sign(r, "api1234567secret", "")
	expected = "//1z5OyAwlLKdQZL5vG2wBIx5CM="
	if signed != expected {
		t.Error("Expected", expected, "found", signed)
	}
}

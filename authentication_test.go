package flickr

import (
	"testing"

	flickErr "gopkg.in/masci/flickr.v2/error"
)

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

	Expect(t, fclient.Id, "21207597@N07")
	Expect(t, fclient.OAuthToken, "72157626318069415-087bfc7b5816092c")
	Expect(t, fclient.OAuthTokenSecret, "a202d1f853ec69de")
}

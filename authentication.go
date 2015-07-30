package flickr

import (
	"io/ioutil"
	"net/url"
	"strconv"
	"strings"

	flickErr "github.com/masci/flickr/error"
)

// Type representing a request token during the exchange process
type RequestToken struct {
	// Whether the callback url matches the one provided in Flickr dashboard
	OauthCallbackConfirmed bool
	// Request token
	OauthToken string
	// Request token secret
	OauthTokenSecret string
	// OAuth failing reason in case of errors
	OAuthProblem string
}

// Extract a RequestToken from the response body
func ParseRequestToken(response string) (*RequestToken, error) {
	val, err := url.ParseQuery(strings.TrimSpace(response))
	if err != nil {
		return nil, err
	}

	ret := &RequestToken{}

	oauth_problem := val.Get("oauth_problem")
	if oauth_problem != "" {
		ret.OAuthProblem = oauth_problem
		return ret, flickErr.NewError(20)
	}

	confirmed, _ := strconv.ParseBool(val.Get("oauth_callback_confirmed"))
	ret.OauthCallbackConfirmed = confirmed
	ret.OauthToken = val.Get("oauth_token")
	ret.OauthTokenSecret = val.Get("oauth_token_secret")

	return ret, nil
}

// Type representing a OAuth access token along with its owner's data
type OAuthToken struct {
	// OAuth access token
	OAuthToken string
	// OAuth access token secret
	OAuthTokenSecret string
	// Flickr ID of token's owner
	UserNsid string
	// Flickr Username of token's owner
	Username string
	// Flickr full name of token's owner
	Fullname string
	// OAuth failing reason in case of errors
	OAuthProblem string
}

// Extract a OAuthToken from the response body
func ParseOAuthToken(response string) (*OAuthToken, error) {
	val, err := url.ParseQuery(strings.TrimSpace(response))
	if err != nil {
		return nil, err
	}

	ret := &OAuthToken{}

	oauth_problem := val.Get("oauth_problem")
	if oauth_problem != "" {
		ret.OAuthProblem = oauth_problem
		return ret, flickErr.NewError(30)
	}

	ret.OAuthToken = val.Get("oauth_token")
	ret.OAuthTokenSecret = val.Get("oauth_token_secret")
	ret.Fullname = val.Get("fullname")
	ret.UserNsid = val.Get("user_nsid")
	ret.Username = val.Get("username")

	return ret, nil
}

// Retrieve a request token: this is the first step to get a fully functional
// access token from Flickr
func GetRequestToken(client *FlickrClient) (*RequestToken, error) {
	client.EndpointUrl = REQUEST_TOKEN_URL
	client.SetOAuthDefaults()
	client.Args.Set("oauth_consumer_key", client.ApiKey)
	client.Args.Set("oauth_callback", "oob")

	// we don't have token secret at this stage, pass an empty string
	client.Sign("")

	res, err := client.HTTPClient.Get(client.GetUrl())
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return ParseRequestToken(string(body))
}

// Returns the URL users need to reach to grant permission to our application
func GetAuthorizeUrl(client *FlickrClient, reqToken *RequestToken) (string, error) {
	client.EndpointUrl = AUTHORIZE_URL
	client.Args = url.Values{}
	client.Args.Set("oauth_token", reqToken.OauthToken)
	// TODO make permission value parametric
	client.Args.Set("perms", "delete")

	return client.GetUrl(), nil
}

// Get an access token providing an OAuth verifier provided by Flickr once the user
// authorizes your application
func GetAccessToken(client *FlickrClient, reqToken *RequestToken, oauthVerifier string) (*OAuthToken, error) {
	client.EndpointUrl = ACCESS_TOKEN_URL
	client.SetOAuthDefaults()
	client.Args.Set("oauth_verifier", oauthVerifier)
	client.Args.Set("oauth_consumer_key", client.ApiKey)
	client.Args.Set("oauth_token", reqToken.OauthToken)
	// use the request token for signing
	client.Sign(reqToken.OauthTokenSecret)

	res, err := client.HTTPClient.Get(client.GetUrl())
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	accessTok, err := ParseOAuthToken(string(body))

	// set client params for convenience
	client.OAuthToken = accessTok.OAuthToken
	client.OAuthTokenSecret = accessTok.OAuthTokenSecret
	client.Id = accessTok.UserNsid

	return accessTok, err
}

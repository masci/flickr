package flickr

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	API_ENDPOINT = "https://api.flickr.com/services/rest"
)

// Generate a random string of 8 chars, needed for OAuth signature
func generateNonce() string {
	rand.Seed(time.Now().UTC().UnixNano())
	// For convenience, use a set of chars we don't need to url-escape
	var letters = []rune("123456789abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ")
	b := make([]rune, 8)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

type FlickrClient struct {
	ApiKey           string
	ApiSecret        string
	HTTPClient       *http.Client
	EndpointUrl      string
	HTTPVerb         string
	Args             url.Values
	OAuthToken       string
	OAuthTokenSecret string
}

func NewFlickrClient(apiKey string, apiSecret string) *FlickrClient {
	return &FlickrClient{
		ApiKey:     apiKey,
		ApiSecret:  apiSecret,
		HTTPClient: &http.Client{},
		HTTPVerb:   "GET",
		Args:       url.Values{},
	}
}

// Sign the next request performed by the FlickrClient
func (c *FlickrClient) Sign(tokenSecret string) {
	// the "oauth_signature" param must not be included in the signing process
	c.Args.Del("oauth_signature")
	c.Args.Set("oauth_signature", c.getSignature(tokenSecret))
}

func (c *FlickrClient) GetUrl() string {
	return fmt.Sprintf("%s?%s", c.EndpointUrl, c.Args.Encode())
}

func (c *FlickrClient) ClearArgs() {
	c.Args = url.Values{}
}

// Set a default set of args needed for signing a request
func (c *FlickrClient) SetDefaultArgs() {
	c.Args = url.Values{}
	c.Args.Add("oauth_version", "1.0")
	c.Args.Add("oauth_signature_method", "HMAC-SHA1")
	c.Args.Add("oauth_nonce", generateNonce())
	c.Args.Add("oauth_timestamp", fmt.Sprintf("%d", time.Now().Unix()))
}

func (c *FlickrClient) getSigningBaseString() string {
	request_url := url.QueryEscape(c.EndpointUrl)
	query := url.QueryEscape(c.Args.Encode())

	return fmt.Sprintf("%s&%s&%s", c.HTTPVerb, request_url, query)
}

func (c *FlickrClient) getSignature(token_secret string) string {
	key := fmt.Sprintf("%s&%s", url.QueryEscape(c.ApiSecret), url.QueryEscape(token_secret))
	base_string := c.getSigningBaseString()

	mac := hmac.New(sha1.New, []byte(key))
	mac.Write([]byte(base_string))

	ret := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	return ret
}

type FlickrResponse struct {
	XMLName xml.Name `xml:"rsp"`
	Status  string   `xml:"stat,attr"`
	Error   struct {
		XMLName xml.Name `xml:"err"`
		Code    int      `xml:"code,attr"`
		Message string   `xml:"msg,attr"`
	}
}

func (r *FlickrResponse) HasErrors() bool {
	return r.Status == "fail"
}

func (r *FlickrResponse) ErrorCode() int {
	return r.Error.Code
}

func (r *FlickrResponse) ErrorMsg() string {
	return r.Error.Message
}

type RequestToken struct {
	OauthCallbackConfirmed bool
	OauthToken             string
	OauthTokenSecret       string
}

func NewRequestToken(response string) (*RequestToken, error) {
	// TODO parse flickr errors inside the body
	val, err := url.ParseQuery(strings.TrimSpace(response))
	if err != nil {
		return nil, err
	}

	confirmed, _ := strconv.ParseBool(val.Get("oauth_callback_confirmed"))

	return &RequestToken{
		confirmed,
		val.Get("oauth_token"),
		val.Get("oauth_token_secret"),
	}, nil
}

type OAuthToken struct {
	OAuthToken       string
	OAuthTokenSecret string
	UserNsid         string
	Username         string
	Fullname         string
}

func NewOAuthToken(response string) (*OAuthToken, error) {
	// TODO parse flickr errors inside the body
	val, err := url.ParseQuery(strings.TrimSpace(response))
	if err != nil {
		return nil, err
	}

	return &OAuthToken{
		OAuthToken:       val.Get("oauth_token"),
		OAuthTokenSecret: val.Get("oauth_token_secret"),
		Fullname:         val.Get("fullname"),
		UserNsid:         val.Get("user_nsid"),
		Username:         val.Get("username"),
	}, nil
}

func GetRequestToken(client *FlickrClient) (*RequestToken, error) {
	client.EndpointUrl = "https://www.flickr.com/services/oauth/request_token"
	client.SetDefaultArgs()
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

	return NewRequestToken(string(body))
}

func GetAuthorizeUrl(client *FlickrClient, reqToken *RequestToken) (string, error) {
	client.EndpointUrl = "https://www.flickr.com/services/oauth/authorize"
	client.Args = url.Values{}
	client.Args.Set("oauth_token", reqToken.OauthToken)
	client.Args.Set("perms", "delete")

	return client.GetUrl(), nil
}

func GetAccessToken(client *FlickrClient, reqToken *RequestToken, oauthVerifier string) (*OAuthToken, error) {
	client.EndpointUrl = "https://www.flickr.com/services/oauth/access_token"
	client.SetDefaultArgs()
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

	return NewOAuthToken(string(body))
}

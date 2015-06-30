package flickr

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"sort"
	"time"
)

const (
	API_ENDPOINT      = "https://api.flickr.com/services/rest"
	AUTHORIZE_URL     = "https://www.flickr.com/services/oauth/authorize"
	REQUEST_TOKEN_URL = "https://www.flickr.com/services/oauth/request_token"
	ACCESS_TOKEN_URL  = "https://www.flickr.com/services/oauth/access_token"
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

// An utility type to wrap all resources and data needed to complete requests
// to the Flickr API
type FlickrClient struct {
	// Flickr application api key
	ApiKey string
	// Flickr application api secret
	ApiSecret string
	// A generic HTTP client to perform GET and POST requests
	HTTPClient *http.Client
	// The base url for API endpoints
	EndpointUrl string
	// A string containing POST or GET, needed for OAuth signing
	HTTPVerb string
	// A set of url params to query the API
	Args url.Values
	// User access token
	OAuthToken string
	// User secret token
	OAuthTokenSecret string
}

// Create a Flickr client, apiKey and apiSecret are mandatory
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

// Specific signing process for API calls, it's not the same as OAuth sign
func (c *FlickrClient) ApiSign(tokenSecret string) {
	// the "api_sig" param must not be included in the signing process
	c.Args.Del("api_sig")
	c.Args.Set("api_sig", c.getApiSignature(tokenSecret))
}

// Evaluate the complete URL to make requests (base url + params)
func (c *FlickrClient) GetUrl() string {
	return fmt.Sprintf("%s?%s", c.EndpointUrl, c.Args.Encode())
}

// Remove all query params
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

// Get the base string to compose the signature
func (c *FlickrClient) getSigningBaseString() string {
	request_url := url.QueryEscape(c.EndpointUrl)
	query := url.QueryEscape(c.Args.Encode())

	return fmt.Sprintf("%s&%s&%s", c.HTTPVerb, request_url, query)
}

// Compute the signature of a signed request
func (c *FlickrClient) getSignature(token_secret string) string {
	key := fmt.Sprintf("%s&%s", url.QueryEscape(c.ApiSecret), url.QueryEscape(token_secret))
	base_string := c.getSigningBaseString()

	mac := hmac.New(sha1.New, []byte(key))
	mac.Write([]byte(base_string))

	ret := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	return ret
}

func (c *FlickrClient) getApiSignature(token_secret string) string {
	var buf bytes.Buffer
	buf.WriteString(token_secret)

	keys := make([]string, 0, len(c.Args))
	for k := range c.Args {
		keys = append(keys, k)
	}
	// args needs to be in alphabetical order
	sort.Strings(keys)

	for _, k := range keys {
		arg := c.Args[k][0]
		buf.WriteString(k)
		buf.WriteString(arg)
	}

	base := buf.String()

	data := []byte(base)
	return fmt.Sprintf("%x", md5.Sum(data))
}

// Base type representing responses from Flickr API
type FlickrResponse struct {
	XMLName xml.Name `xml:"rsp"`
	// Status might contain "err" or "ok" strings
	Status string `xml:"stat,attr"`
	// Flickr API error detail
	Error struct {
		XMLName xml.Name `xml:"err"`
		Code    int      `xml:"code,attr"`
		Message string   `xml:"msg,attr"`
	}
}

// Return whether a response contains errors
func (r *FlickrResponse) HasErrors() bool {
	return r.Status == "fail"
}

// Return the error code (0 if no errors)
func (r *FlickrResponse) ErrorCode() int {
	return r.Error.Code
}

// Return error message string (empty string if no errors)
func (r *FlickrResponse) ErrorMsg() string {
	return r.Error.Message
}

// TODO docs
func GetResponse(client *FlickrClient, r interface{}) error {
	res, err := client.HTTPClient.Get(client.GetUrl())
	if err != nil {
		return err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	err = xml.Unmarshal([]byte(body), r)
	if err != nil {
		return err
	}

	return nil
}

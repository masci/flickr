package flickr

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
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
	// User flickr ID
	Id string
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

// Set the mandatory params for an OAuth request
func (c *FlickrClient) SetOAuthDefaults() {
	c.Args.Add("oauth_version", "1.0")
	c.Args.Add("oauth_signature_method", "HMAC-SHA1")
	c.Args.Add("oauth_nonce", generateNonce())
	c.Args.Add("oauth_timestamp", fmt.Sprintf("%d", time.Now().Unix()))
}

// Sign the request with a default set of OAuth parameters, needed to authorize
// users for certain writing/destructive operations.
func (c *FlickrClient) OAuthSign() {
	c.SetOAuthDefaults()
	c.Args.Set("oauth_token", c.OAuthToken)
	c.Args.Set("oauth_consumer_key", c.ApiKey)
	c.Args.Set("api_key", c.ApiKey)

	c.Sign(c.OAuthTokenSecret)
}

// Specific signing process for API calls: not the same as OAuth sign, used
// for requests that don't need user authorizations.
func (c *FlickrClient) ApiSign() {
	c.Args.Set("api_key", c.ApiKey)
	// the "api_sig" param must not be included in the signing process
	c.Args.Del("api_sig")
	c.Args.Set("api_sig", c.getApiSignature(c.ApiSecret))
}

// Evaluate the complete URL to make requests (base url + params)
func (c *FlickrClient) GetUrl() string {
	return fmt.Sprintf("%s?%s", c.EndpointUrl, c.Args.Encode())
}

// Remove all query params
func (c *FlickrClient) ClearArgs() {
	c.Args = url.Values{}
}

// Reset Args and set the default endpoint
func (c *FlickrClient) Init() {
	c.ClearArgs()
	c.EndpointUrl = API_ENDPOINT
}

// Get the base string to compose the signature
func (c *FlickrClient) getSigningBaseString() string {
	request_url := url.QueryEscape(c.EndpointUrl)
	flickr_encoded := strings.Replace(c.Args.Encode(), "+", "%20", -1)
	query := url.QueryEscape(flickr_encoded)

	ret := fmt.Sprintf("%s&%s&%s", c.HTTPVerb, request_url, query)
	return ret
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

// Sign API requests. This method differs from the signing process needed for
// OAuth authenticated requests.
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

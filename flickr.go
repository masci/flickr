package flickr

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type FlickrClient struct {
	ApiKey      string
	ApiSecret   string
	HTTPClient  *http.Client
	EndpointUrl string
	Method      string
	Args        url.Values
}

func NewFlickrClient(apiKey string, apiSecret string) *FlickrClient {
	return &FlickrClient{
		ApiKey:     apiKey,
		ApiSecret:  apiSecret,
		HTTPClient: &http.Client{},
		Method:     "GET",
		Args:       url.Values{},
	}
}

func (c *FlickrClient) Sign(tokenSecret string) {
	// the "oauth_signature" param should not be included in the signing process
	c.Args.Del("oauth_signature")
	c.Args.Set("oauth_signature", getSignature(c, tokenSecret))
}

func (c *FlickrClient) GetUrl() string {
	return fmt.Sprintf("%s?%s", c.EndpointUrl, c.Args.Encode())
}

type RequestToken struct {
	OauthCallbackConfirmed bool
	OauthToken             string
	OauthTokenSecret       string
}

func (rt *RequestToken) Parse(response string) error {
	val, err := url.ParseQuery(strings.TrimSpace(response))
	if err != nil {
		return err
	}

	confirmed, _ := strconv.ParseBool(val.Get("oauth_callback_confirmed"))
	rt.OauthCallbackConfirmed = confirmed
	rt.OauthToken = val.Get("oauth_token")
	rt.OauthTokenSecret = val.Get("oauth_token_secret")

	return nil
}

func getSigningBaseString(client *FlickrClient) string {
	request_url := url.QueryEscape(client.EndpointUrl)
	query := url.QueryEscape(client.Args.Encode())

	return fmt.Sprintf("%s&%s&%s", client.Method, request_url, query)
}

func getSignature(client *FlickrClient, token_secret string) string {
	key := fmt.Sprintf("%s&%s", url.QueryEscape(client.ApiSecret), url.QueryEscape(token_secret))
	base_string := getSigningBaseString(client)

	mac := hmac.New(sha1.New, []byte(key))
	mac.Write([]byte(base_string))

	ret := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	return ret
}

func generateNonce() string {
	rand.Seed(time.Now().UTC().UnixNano())
	var letters = []rune("123456789abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ")
	b := make([]rune, 8)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func getDefaultArgs() url.Values {
	args := url.Values{}
	args.Add("oauth_version", "1.0")
	args.Add("oauth_signature_method", "HMAC-SHA1")
	args.Add("oauth_nonce", generateNonce())
	args.Add("oauth_timestamp", fmt.Sprintf("%d", time.Now().Unix()))

	return args
}

func GetRequestToken(client *FlickrClient) (*RequestToken, error) {
	client.EndpointUrl = "https://www.flickr.com/services/oauth/request_token"
	client.Args = getDefaultArgs()
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

	token := RequestToken{}
	token.Parse(string(body))

	return &token, nil
}

func GetAuthorizeUrl(client *FlickrClient, reqToken *RequestToken) (string, error) {
	client.EndpointUrl = "https://www.flickr.com/services/oauth/authorize"
	client.Args = url.Values{}
	client.Args.Set("oauth_token", reqToken.OauthToken)
	client.Args.Set("perms", "delete")

	return client.GetUrl(), nil
}

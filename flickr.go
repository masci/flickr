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
	"time"
)

type Request struct {
	Url    string
	Method string
	Args   url.Values
}

func NewRequest(url string, method string, args url.Values) *Request {
	r := Request{url, method, args}

	return &r
}

type RequestToken struct {
	OauthCallbackConfirmed bool
	OauthToken             string
	OauthTokenSecret       string
}

func (rt *RequestToken) Parse(response string) error {
	val, err := url.ParseQuery(response)
	if err != nil {
		return err
	}

	confirmed, _ := strconv.ParseBool(val.Get("oauth_callback_confirmed"))
	rt.OauthCallbackConfirmed = confirmed
	rt.OauthToken = val.Get("oauth_token")
	rt.OauthTokenSecret = val.Get("oauth_token_secret")

	return nil
}

func getSigningBaseString(request *Request) string {
	request_url := url.QueryEscape(request.Url)
	query := url.QueryEscape(request.Args.Encode())

	return fmt.Sprintf("%s&%s&%s", request.Method, request_url, query)
}

func Sign(request *Request, consumer_secret string, token_secret string) string {
	key := fmt.Sprintf("%s&%s", url.QueryEscape(consumer_secret), url.QueryEscape(token_secret))
	base_string := getSigningBaseString(request)

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

func GetRequestToken(api_key string, api_secret string) (*RequestToken, error) {
	base_url := "https://www.flickr.com/services/oauth/request_token"

	args := getDefaultArgs()
	args.Add("oauth_consumer_key", api_key)
	args.Add("oauth_callback", "oob")

	request := NewRequest(base_url, "GET", args)
	// we don't have token secret at this stage, pass an empty string
	request.Args.Add("oauth_signature", Sign(request, api_secret, ""))

	api_url := fmt.Sprintf("%s?%s", base_url, request.Args.Encode())
	res, err := http.Get(api_url)
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

	fmt.Println(api_url)
	fmt.Println(res)

	return &token, nil
}

func GetAuthorizeUrl() {
}

package flickr

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"math/rand"
	"net/url"
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

func GetRequestToken(api_key string, api_secret string) {
	base_url := "https://www.flickr.com/services/oauth/request_token"

	args := getDefaultArgs()
	args.Add("oauth_consumer_key", api_key)
	args.Add("oauth_callback", "oob")

	request := NewRequest(base_url, "GET", args)
	// we don't have token secret at this stage, pass an empty string
	request.Args.Add("oauth_signature", Sign(request, api_secret, ""))

	api_url := fmt.Sprintf("%s?%s", base_url, request.Args.Encode())
	fmt.Println(api_url)
}

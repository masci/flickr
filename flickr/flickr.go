// Flickr.go is a Go library for accessing Flickr API https://www.flickr.com/services/api
package flickr

import (
	"bytes"
	"mime/multipart"
)

const (
	API_ENDPOINT      = "https://api.flickr.com/services/rest"
	UPLOAD_ENDPOINT   = "https://up.flickr.com/services/upload/"
	AUTHORIZE_URL     = "https://www.flickr.com/services/oauth/authorize"
	REQUEST_TOKEN_URL = "https://www.flickr.com/services/oauth/request_token"
	ACCESS_TOKEN_URL  = "https://www.flickr.com/services/oauth/access_token"
)

// Perform a GET request to the Flickr API with the configured FlickrClient passed as first
// parameter. Results will be unmarshalled to fill in a FlickrResponse struct passed as
// second parameter.
func DoGet(client *FlickrClient, r FlickrResponse) error {
	res, err := client.HTTPClient.Get(client.GetUrl())
	if err != nil {
		return err
	}

	return parseApiResponse(res, r)
}

// Perform a POST request to the Flickr API with the configured FlickrClient, the
// request body and the body content type. Results will be unmarshalled in a FlickrResponse
// struct.
func DoPostBody(client *FlickrClient, body *bytes.Buffer, bodyType string, r FlickrResponse) error {
	res, err := client.HTTPClient.Post(client.EndpointUrl, bodyType, body)
	if err != nil {
		return err
	}

	return parseApiResponse(res, r)
}

// Perform a POST request to the Flickr API with the configured FlickrClient,
// dumping client Args into the request Body.
func DoPost(client *FlickrClient, r FlickrResponse) error {
	// instance an empty request body
	body := &bytes.Buffer{}
	// multipart writer to fill the body
	writer := multipart.NewWriter(body)
	// dump params
	for key, val := range client.Args {
		_ = writer.WriteField(key, val[0])
	}
	err := writer.Close()
	if err != nil {
		return err
	}
	// evaluate the content type and the boundary
	contentType := writer.FormDataContentType()

	return DoPostBody(client, body, contentType, r)
}

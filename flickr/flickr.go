// Flickr.go is a Go library for accessing Flickr API https://www.flickr.com/services/api
package flickr

import (
	"bytes"
	"encoding/xml"
	"io/ioutil"
	"mime/multipart"
	"net/http"

	flickErr "github.com/masci/flickr.go/flickr/error"
)

const (
	API_ENDPOINT      = "https://api.flickr.com/services/rest"
	UPLOAD_ENDPOINT   = "https://up.flickr.com/services/upload/"
	AUTHORIZE_URL     = "https://www.flickr.com/services/oauth/authorize"
	REQUEST_TOKEN_URL = "https://www.flickr.com/services/oauth/request_token"
	ACCESS_TOKEN_URL  = "https://www.flickr.com/services/oauth/access_token"
)

// TODO docs
type FlickrResponse interface {
	HasErrors() bool
	ErrorCode() int
	ErrorMsg() string
	SetErrorStatus(bool)
	SetErrorCode(int)
	SetErrorMsg(string)
}

// Base type representing responses from Flickr API
type BasicResponse struct {
	XMLName xml.Name `xml:"rsp"`
	// Status might contain "fail" or "ok" strings
	Status string `xml:"stat,attr"`
	// Flickr API error detail
	Error struct {
		Code    int    `xml:"code,attr"`
		Message string `xml:"msg,attr"`
	} `xml:"err"`
}

// Return whether a response contains errors
func (r *BasicResponse) HasErrors() bool {
	return r.Status != "ok"
}

// Return the error code (0 if no errors)
func (r *BasicResponse) ErrorCode() int {
	return r.Error.Code
}

// Return error message string (empty string if no errors)
func (r *BasicResponse) ErrorMsg() string {
	return r.Error.Message
}

// TODO
func (r *BasicResponse) SetErrorStatus(hasErrors bool) {
	if hasErrors {
		r.Status = "fail"
	} else {
		r.Status = "ok"
	}
}

// TODO
func (r *BasicResponse) SetErrorCode(code int) {
	r.Error.Code = code
}

// TODO
func (r *BasicResponse) SetErrorMsg(msg string) {
	r.Error.Message = msg
}

// Given an http.Response retrieved from Flickr, unmarshal results
// into a FlickrResponse struct.
func parseApiResponse(res *http.Response, r FlickrResponse) error {
	defer res.Body.Close()
	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	err = xml.Unmarshal(responseBody, r)
	if err != nil {
		// In case of OAuth errors (signature, parameters, etc) Flicker does not
		// return a REST response but raw text (!).
		// We need to artificially build a FlickrResponse and manually fill in
		// the error string
		r.SetErrorStatus(true)
		r.SetErrorCode(-1)
		r.SetErrorMsg(string(responseBody))
	}

	if r.HasErrors() {
		return flickErr.NewError(10)
	}

	return nil
}

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

// TODO docs
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

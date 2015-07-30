package flickr

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"

	flickErr "github.com/masci/flickr/error"
)

// Interface for Flickr request objects
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
	Extra string `xml:",innerxml"`
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

// Set error status explicitly
func (r *BasicResponse) SetErrorStatus(hasErrors bool) {
	if hasErrors {
		r.Status = "fail"
	} else {
		r.Status = "ok"
	}
}

// Set error code explicitly
func (r *BasicResponse) SetErrorCode(code int) {
	r.Error.Code = code
}

// Set error message explicitly
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

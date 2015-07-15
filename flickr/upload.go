package flickr

import (
	"bytes"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Encode the file and request parameters in a multipart body
func getUploadBody(client *FlickrClient, photo io.Reader, fileName string) (*bytes.Buffer, string, error) {
	// instance an empty request body
	body := &bytes.Buffer{}
	// multipart writer to fill the body
	writer := multipart.NewWriter(body)
	// dump the file in the "photo" field
	part, err := writer.CreateFormFile("photo", filepath.Base(fileName))
	if err != nil {
		return nil, "", err
	}
	_, err = io.Copy(part, photo)
	// dump other params
	for key, val := range client.Args {
		_ = writer.WriteField(key, val[0])
	}
	err = writer.Close()
	if err != nil {
		return nil, "", err
	}
	// evaluate the content type and the boundary
	contentType := writer.FormDataContentType()

	return body, contentType, nil
}

// A convenience struct wrapping all optional upload parameters
type UploadParams struct {
	Title, Description           string
	Tags                         []string
	IsPublic, IsFamily, IsFriend bool
	ContentType                  int
	Hidden                       int
	SafetyLevel                  int
}

// Provide meaningful default values
func NewUploadParams() *UploadParams {
	ret := &UploadParams{}
	ret.ContentType = 1 // photo
	ret.Hidden = 2      // hidden from public searchesi
	ret.SafetyLevel = 1 // safe
	return ret
}

// Type representing a successful upload response from the api
type UploadResponse struct {
	BasicResponse
	Id string `xml:"photoid"`
}

// Set client query arguments based on the contents of the UploadParams struct
func fillArgsWithParams(client *FlickrClient, params *UploadParams) {
	if params.Title != "" {
		client.Args.Set("title", params.Title)
	}

	if params.Description != "" {
		client.Args.Set("description", params.Description)
	}

	if len(params.Tags) > 0 {
		client.Args.Set("tags", strings.Join(params.Tags, " "))
	}

	var boolString = func(b bool) string {
		if b {
			return "1"
		}
		return "0"
	}
	client.Args.Set("is_public", boolString(params.IsPublic))
	client.Args.Set("is_friend", boolString(params.IsFriend))
	client.Args.Set("is_family", boolString(params.IsFamily))

	if params.ContentType >= 1 && params.ContentType <= 3 {
		client.Args.Set("content_type", strconv.Itoa(params.ContentType))
	}

	if params.Hidden >= 1 && params.Hidden <= 2 {
		client.Args.Set("hidden", strconv.Itoa(params.Hidden))
	}

	if params.SafetyLevel >= 1 && params.SafetyLevel <= 3 {
		client.Args.Set("safety_level", strconv.Itoa(params.SafetyLevel))
	}
}

// Perform a file upload using the Flickr API. If optionalParams is nil,
// no parameters will be added to the request and Flickr will set User's
// default preferences.
// This call must be signed with write permissions
func UploadPhoto(client *FlickrClient, path string, optionalParams *UploadParams) (*UploadResponse, error) {
	client.EndpointUrl = UPLOAD_ENDPOINT
	client.HTTPVerb = "POST"
	client.SetOAuthDefaults()

	if optionalParams != nil {
		fillArgsWithParams(client, optionalParams)
	}

	client.OAuthSign()

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body, ctype, err := getUploadBody(client, file, file.Name())
	if err != nil {
		return nil, err
	}

	resp := &UploadResponse{}
	err = DoPostBody(client, body, ctype, resp)
	return resp, err
}

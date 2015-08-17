package flickr

import (
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// generate a random multipart boundary string,
// shamelessly copypasted from the std library
func randomBoundary() string {
	var buf [30]byte
	_, err := io.ReadFull(rand.Reader, buf[:])
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%x", buf[:])
}

// Encode the file and request parameters in a multipart body.
// File contents are streamed into the request using an io.Pipe in a separated goroutine
func streamUploadBody(client *FlickrClient, photo io.Reader, body *io.PipeWriter, fileName string, boundary string) {
	// multipart writer to fill the body
	defer body.Close()
	writer := multipart.NewWriter(body)
	writer.SetBoundary(boundary)

	// create the "photo" field
	part, err := writer.CreateFormFile("photo", filepath.Base(fileName))
	if err != nil {
		log.Fatal(err)
		return
	}

	// fill the photo field
	_, err = io.Copy(part, photo)
	if err != nil {
		log.Fatal(err)
		return
	}

	// dump other params
	for key, val := range client.Args {
		_ = writer.WriteField(key, val[0])
	}

	// close the form writer
	err = writer.Close()
	if err != nil {
		log.Fatal(err)
		return
	}
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
func UploadFile(client *FlickrClient, path string, optionalParams *UploadParams) (*UploadResponse, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return UploadReader(client, file, file.Name(), optionalParams)
}

// Same as UploadFile but the photo file is passed as an io.Reader instead of a file path
func UploadReader(client *FlickrClient, photoReader io.Reader, name string, optionalParams *UploadParams) (*UploadResponse, error) {
	client.Init()
	client.EndpointUrl = UPLOAD_ENDPOINT
	client.HTTPVerb = "POST"

	if optionalParams != nil {
		fillArgsWithParams(client, optionalParams)
	}

	client.OAuthSign()

	// write request body in a Pipe
	boundary := randomBoundary()
	r, w := io.Pipe()
	go streamUploadBody(client, photoReader, w, name, boundary)

	// create an HTTP Request
	req, err := http.NewRequest("POST", client.EndpointUrl, r)
	if err != nil {
		return nil, err
	}

	// set content-type
	req.Header.Set("Content-Type", "multipart/form-data; boundary="+boundary)

	// instance an HTTP client
	http_client := &http.Client{}
	// perform upload request streaming the file
	resp, err := http_client.Do(req)
	if err != nil {
		return nil, err
	}

	api_resp := &UploadResponse{}
	err = parseApiResponse(resp, api_resp)
	return api_resp, err
}

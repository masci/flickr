package flickr

import (
	"bytes"
	"encoding/xml"
	flickErr "github.com/masci/flickr.go/flickr/error"
	"io"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path/filepath"
)

// TODO docs
type UploadParams struct {
	Title, Description, Tags     string
	IsPublic, IsFamily, IsFriend bool
	ContentType                  int
	Hidden                       bool
	SafetyLevel                  bool
}

// TODO docs
func NewUploadParams() *UploadParams {
	ret := &UploadParams{}
	return ret
}

// TODO docs
func getUploadBody(client *FlickrClient, file *os.File) (*bytes.Buffer, string, error) {
	// instance an empty request body
	body := &bytes.Buffer{}
	// multipart writer to fill the body
	writer := multipart.NewWriter(body)
	// dump the file in the "photo" field
	part, err := writer.CreateFormFile("photo", filepath.Base(file.Name()))
	if err != nil {
		return nil, "", err
	}
	_, err = io.Copy(part, file)
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

type UploadResponse struct {
	FlickrResponse
	Id int `xml:"photoid"`
}

// TODO docs
func UploadPhoto(client *FlickrClient, path string, optionalParams *UploadParams) (*UploadResponse, error) {
	client.EndpointUrl = UPLOAD_ENDPOINT
	client.HTTPVerb = "POST"
	client.SetDefaultArgs()
	client.Args.Set("oauth_token", client.OAuthToken)
	client.Args.Set("oauth_consumer_key", client.ApiKey)

	if optionalParams == nil {
		optionalParams = NewUploadParams()
	}

	//client.Args.Set("title", optionalParams.Title)
	// TODO finish filling args with optional params
	// ...
	client.Sign(client.OAuthTokenSecret)

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body, ctype, err := getUploadBody(client, file)
	if err != nil {
		return nil, err
	}

	res, err := client.HTTPClient.Post(client.EndpointUrl, ctype, body)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	bodyResponse, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	resp := &UploadResponse{}
	err = xml.Unmarshal(bodyResponse, resp)
	if err != nil {
		return nil, err
	}

	if resp.HasErrors() {
		return resp, flickErr.NewError(10)
	}

	return resp, nil
}

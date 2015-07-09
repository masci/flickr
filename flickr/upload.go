package flickr

import (
	"bytes"
	"encoding/xml"
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
func UploadPhoto(client *FlickrClient, path string, optionalParams *UploadParams) (int, error) {
	client.EndpointUrl = UPLOAD_ENDPOINT
	client.ClearArgs()
	client.Args.Set("api_key", client.ApiKey)

	if optionalParams == nil {
		optionalParams = NewUploadParams()
	}

	client.Args.Set("title", optionalParams.Title)
	// TODO finish filling args with optional params
	// ...
	client.ApiSign(client.ApiSecret)

	file, err := os.Open(path)
	if err != nil {
		return -1, err
	}
	defer file.Close()

	body, ctype, err := getUploadBody(client, file)
	if err != nil {
		return -1, err
	}

	res, err := client.HTTPClient.Post(client.GetUrl(), ctype, body)
	if err != nil {
		return -1, err
	}

	defer res.Body.Close()
	bodyResponse, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return -1, err
	}

	resp := UploadResponse{}
	err = xml.Unmarshal([]byte(bodyResponse), resp)
	if err != nil {
		return -1, err
	}

	return resp.Id, nil
}

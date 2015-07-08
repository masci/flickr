package flickr

import (
	"bytes"
	"mime/multipart"
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
func getUploadBody(client *FlickrClient, file *File) (*bytes.Buffer, string, error) {
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
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, "", err
	}
	// evaluate the content type and the boundary
	contentType := writer.FormDataContentType()

	return body, contentType, nil
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

	body := getUploadBody(client, file)

	res, err := client.HTTP.Post(client.GetUrl(), "", body)

}

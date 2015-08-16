package flickr

import (
	"io/ioutil"
	"os"
	"testing"

	flickErr "github.com/masci/flickr/error"
)

func TestNewUploadParams(t *testing.T) {
	params := NewUploadParams()
	Expect(t, params.Title, "")
	Expect(t, params.Description, "")
	Expect(t, len(params.Tags), 0)
	Expect(t, params.IsPublic, false)
	Expect(t, params.IsFamily, false)
	Expect(t, params.IsFriend, false)
	Expect(t, params.ContentType, 1)
	Expect(t, params.Hidden, 2)
	Expect(t, params.SafetyLevel, 1)
}

func TestFillArgsWithParams(t *testing.T) {
	client := GetTestClient()
	params := NewUploadParams()
	fillArgsWithParams(client, params)

	Expect(t, client.Args.Get("title"), "")
	Expect(t, client.Args.Get("description"), "")
	Expect(t, client.Args.Get("tags"), "")
	Expect(t, client.Args.Get("is_public"), "0")
	Expect(t, client.Args.Get("is_friend"), "0")
	Expect(t, client.Args.Get("is_family"), "0")
	Expect(t, client.Args.Get("content_type"), "1")
	Expect(t, client.Args.Get("hidden"), "2")
	Expect(t, client.Args.Get("safety_level"), "1")

	params.Title = "foo"
	params.Description = "a long description"
	params.Tags = []string{"a", "b", "c"}
	params.IsPublic = true
	params.IsFamily = true
	params.IsFriend = true
	params.ContentType = 100
	params.Hidden = 100
	params.SafetyLevel = 100
	client.ClearArgs()
	fillArgsWithParams(client, params)
	Expect(t, client.Args.Get("title"), "foo")
	Expect(t, client.Args.Get("description"), "a long description")
	Expect(t, client.Args.Get("tags"), "a b c")
	Expect(t, client.Args.Get("is_public"), "1")
	Expect(t, client.Args.Get("is_friend"), "1")
	Expect(t, client.Args.Get("is_family"), "1")
	Expect(t, client.Args.Get("content_type"), "")
	Expect(t, client.Args.Get("hidden"), "")
	Expect(t, client.Args.Get("safety_level"), "")
}

func TestUploadFile(t *testing.T) {
	fclient := GetTestClient()
	server, client := FlickrMock(200, `<?xml version="1.0" encoding="utf-8" ?><rsp stat="ok"></rsp>`, "")
	defer server.Close()
	fclient.HTTPClient = client
	params := NewUploadParams()

	resp, err := UploadFile(fclient, "", params)
	Expect(t, resp == nil, true) // comparing nil interfaces would fail
	_, ok := err.(*os.PathError)
	Expect(t, ok, true)
}

func TestUploadFileKo(t *testing.T) {
	fclient := GetTestClient()
	server, client := FlickrMock(200, `<?xml version="1.0" encoding="utf-8" ?><rsp stat="fail"></rsp>`, "")
	defer server.Close()
	fclient.HTTPClient = client

	fooFile, err := ioutil.TempFile("", "flickr.go")
	defer fooFile.Close()
	resp, err := UploadFile(fclient, fooFile.Name(), nil)
	_, ok := err.(*flickErr.Error)
	Expect(t, ok, true)
	Expect(t, resp.HasErrors(), true)
}

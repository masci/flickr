package flickr

import (
	"testing"
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
}

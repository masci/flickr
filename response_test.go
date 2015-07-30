package flickr

import (
	"encoding/xml"
	"net/http"
	"testing"

	flickErr "github.com/masci/flickr/error"
)

func TestFlickrResponse(t *testing.T) {
	failure := `<?xml version="1.0" encoding="utf-8" ?>
<rsp stat="fail">
  <err code="99" msg="Insufficient permissions. Method requires read privileges; none granted." />
</rsp>
`
	resp := FooResponse{}
	err := xml.Unmarshal([]byte(failure), &resp)
	if err != nil {
		t.Error("Error unmarsshalling", failure)
	}

	Expect(t, resp.HasErrors(), true)
	Expect(t, resp.ErrorCode(), 99)
	Expect(t, resp.ErrorMsg(), "Insufficient permissions. Method requires read privileges; none granted.")

	ok := `<?xml version="1.0" encoding="utf-8" ?>
<rsp stat="ok">
  <user id="23148015@N00">
    <username>Massimiliano Pippi</username>
  </user>
  <foo>Foo!</foo>
</rsp>`

	resp = FooResponse{}
	err = xml.Unmarshal([]byte(ok), &resp)
	if err != nil {
		t.Error("Error unmarsshalling", ok)
	}

	Expect(t, resp.HasErrors(), false)
	Expect(t, resp.Foo, "Foo!")
	Expect(t, resp.ErrorCode(), 0)
	Expect(t, resp.ErrorMsg(), "")

	resp = FooResponse{}
	resp.SetErrorStatus(true)
	resp.SetErrorMsg("a message")
	resp.SetErrorCode(999)
	Expect(t, resp.HasErrors(), true)
	Expect(t, resp.ErrorMsg(), "a message")
	Expect(t, resp.ErrorCode(), 999)
	resp.SetErrorStatus(false)
	Expect(t, resp.HasErrors(), false)
}

type FooResponse struct {
	BasicResponse
	Foo string `xml:"foo"`
}

func TestParseResponse(t *testing.T) {
	bodyStr := `<?xml version="1.0" encoding="utf-8" ?>
<rsp stat="ok">
  <user id="23148015@N00">
    <username>Massimiliano Pippi</username>
  </user>
  <foo>Foo!</foo>
</rsp>`

	flickrResp := &FooResponse{}
	response := &http.Response{}
	response.Body = NewFakeBody(bodyStr)

	err := parseApiResponse(response, flickrResp)

	Expect(t, err, nil)
	Expect(t, flickrResp.Foo, "Foo!")

	response = &http.Response{}
	response.Body = NewFakeBody("a_non_rest_format_error")

	err = parseApiResponse(response, flickrResp)
	ferr, ok := err.(*flickErr.Error)
	Expect(t, ok, true)
	Expect(t, ferr.ErrorCode, 10)

	response = &http.Response{}
	response.Body = NewFakeBody(`<?xml version="1.0" encoding="utf-8" ?><rsp stat="fail"></rsp>`)
	err = parseApiResponse(response, flickrResp)
	//ferr, ok := err.(*flickErr.Error)
	//Expect(t, ok, true)
	//Expect(t, ferr.ErrorCode, 10)
}

func TestExtra(t *testing.T) {
	bodyStr := `<?xml version="1.0" encoding="utf-8" ?>
<rsp stat="ok">
  <user id="23148015@N00">
    <username>Massimiliano Pippi</username>
  </user>
  <foo>Foo!</foo>
  <brands>
    <brand id="canon">Canon</brand>
    <brand id="nikon">Nikon</brand>
    <brand id="apple">Apple</brand>
  </brands>
</rsp>`

	flickrResp := &BasicResponse{}
	response := &http.Response{}
	response.Body = NewFakeBody(bodyStr)

	err := parseApiResponse(response, flickrResp)

	Expect(t, err, nil)
	Expect(t, flickrResp.Extra != "", true)
}

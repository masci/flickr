package photos

import (
	"testing"

	"gopkg.in/masci/flickr.v3"
	flickErr "gopkg.in/masci/flickr.v3/error"
)

const photoInfo = `<?xml version="1.0" encoding="utf-8" ?>
<rsp stat="ok">
  <photo id="52435165562" secret="abc" server="65535" farm="66" dateuploaded="1666047672" isfavorite="0" license="0" safety_level="1" rotation="0" originalsecret="9" originalformat="jpg" views="284" media="photo">
    <owner nsid="1687112@N06" username="pankaj.anand" realname="Pankaj Anand" location="Seattle, United States" iconserver="65535" iconfarm="66" path_alias="pankajanand18" />
    <title>Nikki !!</title>
    <description>Seattle, September, 2022</description>
    <visibility ispublic="1" isfriend="0" isfamily="0" />
    <dates posted="1666047672" taken="2022-09-24 08:07:22" takengranularity="0" takenunknown="0" lastupdate="1666073201" />
    <permissions permcomment="3" permaddmeta="2" />
    <editability cancomment="1" canaddmeta="1" />
    <publiceditability cancomment="1" canaddmeta="0" />
    <usage candownload="1" canblog="1" canprint="1" canshare="1" />
    <comments>0</comments>
    <notes />
    <people haspeople="0" />
    <tags>
      <tag id="41641790-52435165562-7257133" author="1687112@N06" authorname="pankaj.anand" raw="body positive" machine_tag="">bodypositive</tag>
      <tag id="41641790-52435165562-404278" author="1687112@N06" authorname="pankaj.anand" raw="female portrait" machine_tag="">femaleportrait</tag>
      <tag id="41641790-52435165562-1567695" author="1687112@N06" authorname="pankaj.anand" raw="outdoor nude" machine_tag="">outdoornude</tag>
      <tag id="41641790-52435165562-2073283" author="1687112@N06" authorname="pankaj.anand" raw="outdoor portrait" machine_tag="">outdoorportrait</tag>
      <tag id="41641790-52435165562-5644201" author="1687112@N06" authorname="pankaj.anand" raw="outdoor shoot" machine_tag="">outdoorshoot</tag>
      <tag id="41641790-52435165562-41272" author="1687112@N06" authorname="pankaj.anand" raw="pink hair" machine_tag="">pinkhair</tag>
      <tag id="41641790-52435165562-2021078" author="1687112@N06" authorname="pankaj.anand" raw="plus size model" machine_tag="">plussizemodel</tag>
      <tag id="41641790-52435165562-1054996" author="1687112@N06" authorname="pankaj.anand" raw="seattle model" machine_tag="">seattlemodel</tag>
      <tag id="41641790-52435165562-1936113" author="1687112@N06" authorname="pankaj.anand" raw="seattle photographer" machine_tag="">seattlephotographer</tag>
      <tag id="41641790-52435165562-216444454" author="1687112@N06" authorname="pankaj.anand" raw="sigma 50mm art" machine_tag="">sigma50mmart</tag>
      <tag id="41641790-52435165562-235861767" author="1687112@N06" authorname="pankaj.anand" raw="sony portraits" machine_tag="">sonyportraits</tag>
      <tag id="41641790-52435165562-501971677" author="1687112@N06" authorname="pankaj.anand" raw="tamron 70-180" machine_tag="">tamron70180</tag>
      <tag id="41641790-52435165562-271371" author="1687112@N06" authorname="pankaj.anand" raw="2022" machine_tag="">2022</tag>
      <tag id="41641790-52435165562-1986" author="1687112@N06" authorname="pankaj.anand" raw="50mm" machine_tag="">50mm</tag>
      <tag id="41641790-52435165562-69" author="1687112@N06" authorname="pankaj.anand" raw="Seattle" machine_tag="">seattle</tag>
      <tag id="41641790-52435165562-5969" author="1687112@N06" authorname="pankaj.anand" raw="Sony" machine_tag="">sony</tag>
      <tag id="41641790-52435165562-486432753" author="1687112@N06" authorname="pankaj.anand" raw="a7iv" machine_tag="">a7iv</tag>
      <tag id="41641790-52435165562-1186" author="1687112@N06" authorname="pankaj.anand" raw="female" machine_tag="">female</tag>
      <tag id="41641790-52435165562-4890" author="1687112@N06" authorname="pankaj.anand" raw="outdoor" machine_tag="">outdoor</tag>
      <tag id="41641790-52435165562-73" author="1687112@N06" authorname="pankaj.anand" raw="park" machine_tag="">park</tag>
      <tag id="41641790-52435165562-36726" author="1687112@N06" authorname="pankaj.anand" raw="tamron" machine_tag="">tamron</tag>
    </tags>
    <urls>
      <url type="photopage">https://www.flickr.com/photos/pankajanand18/52435165562/</url>
    </urls>
  </photo>
</rsp>`

func TestDelete(t *testing.T) {
	fclient := flickr.GetTestClient()
	server, client := flickr.FlickrMock(200, `<?xml version="1.0" encoding="utf-8" ?><rsp stat="ok"></rsp>`, "")
	defer server.Close()
	fclient.HTTPClient = client
	resp, err := Delete(fclient, "123456")
	flickr.Expect(t, err, nil)
	flickr.Expect(t, resp.HasErrors(), false)
}

func TestDeleteKo(t *testing.T) {
	fclient := flickr.GetTestClient()
	server, client := flickr.FlickrMock(200, `<?xml version="1.0" encoding="utf-8" ?><rsp stat="fail"></rsp>`, "")
	defer server.Close()
	fclient.HTTPClient = client
	resp, err := Delete(fclient, "123456")
	_, ok := err.(*flickErr.Error)
	flickr.Expect(t, ok, true)
	flickr.Expect(t, resp.HasErrors(), true)
}

func TestGetTags(t *testing.T) {
	fclient := flickr.GetTestClient()
	server, client := flickr.FlickrMock(200, photoInfo, "")
	defer server.Close()
	fclient.HTTPClient = client
	resp, err := GetInfo(fclient, "123", "")
	_, ok := err.(*flickErr.Error)
	flickr.Expect(t, ok, false)
	if len(resp.Photo.Tags) <= 0 {
		t.Error("Error in parsing.. size of tags should be greater than zero")
	}
	if resp.Photo.Tags[0].Raw != "body positive" {

		t.Error("Error in parsing.. first value should be body positive")
	}
	flickr.Expect(t, resp.HasErrors(), false)
}

package photosets

var (
	resp = `<?xml version="1.0" encoding="utf-8" ?>
		<rsp stat="ok">
			<photosets page="1" pages="1" perpage="10" total="10">
				<photoset id="1234567890" primary="123456" secret="abcdef" server="1234" farm="1" photos="1" videos="0" needs_interstitial="0" visibility_can_see_set="1" count_views="0" count_comments="0" can_comment="0" date_create="1361132046" date_update="1376079704">
					<title>A photoset</title>
					<description />
				</photoset>
				<photoset id="1234567890" primary="123456" secret="abcdef" server="1234" farm="1" photos="5" videos="0" needs_interstitial="0" visibility_can_see_set="1" count_views="17" count_comments="0" can_comment="0" date_create="1135438501" date_update="1375623695">
					<title>Portraits</title>
					<description>Another cool photosets with some pics inside</description>
				</photoset>
			</photosets>
		</rsp>`

	respKo = `<?xml version="1.0" encoding="utf-8" ?>
		<rsp stat="fail">
		  <err code="1" msg="User not found" />
		</rsp>`
)

func TestGetOwnList(t *testing.T) {

}

func TestGetList(t *testing.T) {

}

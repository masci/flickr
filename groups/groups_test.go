package groups

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/masci/flickr.v2"
	flickErr "gopkg.in/masci/flickr.v2/error"
)

func TestGetGroups(t *testing.T) {
	fclient := flickr.GetTestClient()
	server, client := flickr.FlickrMock(200, getGroupsSamplePayload, "")
	defer server.Close()
	fclient.HTTPClient = client
	resp, err := GetGroups(fclient, 0, 0)
	_, ok := err.(*flickErr.Error)
	flickr.Expect(t, ok, false)
	assert.Greater(t, len(resp.Groups), 0, "The size of groups should be greater than zero")
	assert.Equal(t, "ART", resp.Groups[0].Name, "First param should be Name")
}

func TestGetInfo(t *testing.T) {

	fclient := flickr.GetTestClient()
	server, client := flickr.FlickrMock(200, GroupInfoSample, "")
	defer server.Close()
	fclient.HTTPClient = client
	resp, err := GetInfo(fclient, "123")
	_, ok := err.(*flickErr.Error)
	flickr.Expect(t, ok, false)
	assert.NotNil(t, resp.Group.Restriction, "Restrictions can not be nil")
	assert.Equal(t, "4", resp.Group.Throttle.Count, "The size of groups should be greater than zero")
	assert.Equal(t, "1", resp.Group.Restriction.SafeOk, "First param should be Name")
}

func TestAddPhotoGroup(t *testing.T) {

	fclient := flickr.GetTestClient()
	server, client := flickr.FlickrMock(200, GroupInfoSample, "")
	defer server.Close()
	fclient.HTTPClient = client
	_, err := AddPhoto(fclient, "123", "234")
	flickr.Expect(t, err, nil)

	server, client = flickr.FlickrMock(200, `<rsp stat="fail"></rsp>`, "text/xml")
	defer server.Close()
	fclient.HTTPClient = client
	resp, err := AddPhoto(fclient, "123456", "123")
	_, ok := err.(*flickErr.Error)
	flickr.Expect(t, ok, true)
	flickr.Expect(t, resp.HasErrors(), true)

	fclient = flickr.GetTestClient()
	AddPhoto(fclient, "123456", "123")
	params := []string{"photo_id", "group_id"}
	flickr.AssertParamsInBody(t, fclient, params)

}

func TestCanAddPhotos(t *testing.T) {
	groupInfoResponse := &GroupInfoResponse{}
	groupInfoResponse.Group.Throttle.Remaining = "4"
	flickr.Expect(t, groupInfoResponse.CanAddPhotos(), true)
	groupInfoResponse.Group.Throttle.Remaining = "0"
	flickr.Expect(t, groupInfoResponse.CanAddPhotos(), false)
	groupInfoResponse.Group.Throttle.Remaining = ""
	flickr.Expect(t, groupInfoResponse.CanAddPhotos(), false)
}

const GroupInfoSample = `<?xml version="1.0" encoding="utf-8" ?>
<rsp stat="ok">
  <group id="14743297@N22" nsid="14743297@N22" path_alias="_we_need_beauty_and_poetry" iconserver="65535" iconfarm="66" lang="" ispoolmoderated="1" photo_limit_opt_out="1" eighteenplus="1" invitation_only="0" is_member="1" is_moderator="0" is_admin="0" is_founder="0">
    <name>we need poetry and beauty</name>
    <description>La bellezza è un'esigenza del cuore, la poesia è la forza amorosa che la trasforma in qualcosa di immortale. Nelle foto che entreranno in questo gruppo si deve capire che il pulsante che ha scattato la foto non è stato controllato solo da un occhio, ma anche da un soffio d'amore. La tecnica conta poco, ma accetterò solo foto di profonda partecipazione del fotografo e della modella.</description>
    <rules>1. Only female portraits and no faces of children 2. Remember to post, comment and invite when you would like to. 3. No foul language, insults etc. Please: respect! 4. Only submit FOUR PICTURES PER DAY. You can always come back the next day to post more. 5. No pornography, good erotic pictures are suitable. 6. The admin may block you, delete your pictures, etc without any warning. 7. The admin has the right to change the rules at any time without warning. 8 no pictures of violence or death. They will be removed. 9. Inclusion in the pool is very selective, subjective, offbeat, eclectic, and personal. 10. If you want comment what you see you can, but only with your words. &lt;b&gt;&lt;/b&gt;</rules>
    <members>1104</members>
    <pool_count>26047</pool_count>
    <topic_count>6</topic_count>
    <privacy>3</privacy>
    <roles member="membri" moderator="marcomarchetto956" admin="amministratore" />
    <datecreate>2020-04-04 08:23:27</datecreate>
    <dateactivity>1667875715</dateactivity>
    <blast date_blast_added="1666305761" user_id="161229264">Looking for so many beautiful pictures! The extraordinary feminine beauty that surrounds us is painted with poetry and with the thousand confetti of sensuality. This group needs you and your best photos! ....keep on posting...and don't be afraid of commenting and of being commented!  
25577 25658</blast>
    <throttle count="4" mode="day" remaining="4" />
    <restrictions photos_ok="1" videos_ok="0" images_ok="1" screens_ok="0" art_ok="0" virtual_ok="0" safe_ok="1" moderate_ok="1" restricted_ok="1" has_geo="0" />
  </group>
</rsp>`

const getGroupsSamplePayload = `<?xml version="1.0" encoding="utf-8" ?>
<rsp stat="ok">
  <groups page="1" pages="1" per_page="400" total="268">
    <group nsid="58107094@N00" id="58107094@N00" name="ART" member="1" moderator="0" admin="0" privacy="3" photos="2573832" iconserver="105" iconfarm="1" member_count="34486" topic_count="103" pool_count="2573832" />
    <group nsid="2981389@N23" id="2981389@N23" name="!  Elegant Woman HD®" member="1" moderator="0" admin="0" privacy="3" photos="163739" iconserver="65535" iconfarm="66" member_count="19861" topic_count="24" pool_count="163739" />
    <group nsid="3864394@N24" id="3864394@N24" name="! Human - The Greatest Work of Art" member="1" moderator="0" admin="0" privacy="3" photos="62247" iconserver="65535" iconfarm="66" member_count="6516" topic_count="3" pool_count="62247" />
    <group nsid="2820457@N24" id="2820457@N24" name="! Images" member="1" moderator="0" admin="0" privacy="3" photos="123293" iconserver="594" iconfarm="1" member_count="2649" topic_count="2" pool_count="123293" />
    <group nsid="14752479@N23" id="14752479@N23" name="! JL - Naughty PIN-UP" member="1" moderator="0" admin="0" privacy="2" photos="613" iconserver="65535" iconfarm="66" member_count="37" topic_count="0" pool_count="613" />
    <group nsid="1813440@N22" id="1813440@N22" name="! Love Cute Fashion, Lingerie &amp; Swimwear" member="1" moderator="0" admin="0" privacy="3" photos="291213" iconserver="2872" iconfarm="3" member_count="52290" topic_count="1" pool_count="291213" />
    <group nsid="3946646@N21" id="3946646@N21" name="! show your best portraits !" member="1" moderator="0" admin="0" privacy="3" photos="62426" iconserver="65535" iconfarm="66" member_count="8172" topic_count="8" pool_count="62426" />
    <group nsid="888673@N24" id="888673@N24" name="The Wonderful World of Portraits - (only invited photos)" member="1" moderator="0" admin="0" privacy="3" photos="122127" iconserver="4450" iconfarm="5" member_count="11760" topic_count="35" pool_count="122127" />
    <group nsid="65639919@N00" id="65639919@N00" name="!!!!Gorgeous Cuties!!!" member="1" moderator="0" admin="0" privacy="3" photos="113964" iconserver="72" iconfarm="1" member_count="8162" topic_count="15" pool_count="113964" />
    <group nsid="831954@N24" id="831954@N24" name="!!Perfect Composition!!" member="1" moderator="0" admin="0" privacy="3" photos="2436853" iconserver="7398" iconfarm="8" member_count="33870" topic_count="61" pool_count="2436853" />
    <group nsid="2106120@N20" id="2106120@N20" name="Beautiful women ( Beautiful - Attractive - sexy )" member="1" moderator="0" admin="0" privacy="3" photos="259798" iconserver="0" iconfarm="0" member_count="24407" topic_count="18" pool_count="259798" />
    <group nsid="2692194@N22" id="2692194@N22" name="!*ALL ARE WELCOME*!" member="1" moderator="0" admin="0" privacy="3" photos="3478932" iconserver="2904" iconfarm="3" member_count="15213" topic_count="31" pool_count="3478932" />
    <group nsid="926189@N24" id="926189@N24" name="!Art and Photography (seriously18+) INVITE ONLY" member="1" moderator="0" admin="0" privacy="2" photos="624926" iconserver="7891" iconfarm="8" member_count="17000" topic_count="1" pool_count="624926" />
    <group nsid="14753388@N25" id="14753388@N25" name="!Beautiful Women in Socks!" member="1" moderator="0" admin="0" privacy="3" photos="1480" iconserver="65535" iconfarm="66" member_count="1609" topic_count="2" pool_count="1480" />
    <group nsid="14809535@N21" id="14809535@N21" name="!The Most Beautiful and Unique Breasts on Flikr!" member="1" moderator="0" admin="0" privacy="3" photos="13434" iconserver="65535" iconfarm="66" member_count="4707" topic_count="4" pool_count="13434" />
    <group nsid="1029149@N20" id="1029149@N20" name="&quot;Flickr  Women .......Flickr Mujeres.......Flickr Femmes...&quot;" member="1" moderator="0" admin="0" privacy="3" photos="152174" iconserver="3453" iconfarm="4" member_count="9480" topic_count="19" pool_count="152174" />
    <group nsid="615270@N22" id="615270@N22" name="&quot;Unlimited Photos: No Rules&quot;" member="1" moderator="0" admin="0" privacy="3" photos="7078167" iconserver="2068" iconfarm="3" member_count="33052" topic_count="73" pool_count="7078167" />
    <group nsid="608571@N22" id="608571@N22" name="&quot;Women Portraits&quot;" member="1" moderator="0" admin="0" privacy="3" photos="446037" iconserver="7597" iconfarm="8" member_count="14600" topic_count="21" pool_count="446037" />
    <group nsid="14805334@N23" id="14805334@N23" name="* no limits, just quality photography!" member="1" moderator="0" admin="0" privacy="3" photos="302821" iconserver="65535" iconfarm="66" member_count="3404" topic_count="32" pool_count="302821" />
    <group nsid="1071722@N23" id="1071722@N23" name="*** Flickr Sexy and Nude-Art Photography ***" member="1" moderator="0" admin="0" privacy="1" photos="126951" iconserver="2944" iconfarm="3" member_count="800" topic_count="48" pool_count="126951" />
    <group nsid="1136489@N22" id="1136489@N22" name="***Flickr Global" member="1" moderator="0" admin="0" privacy="3" photos="5310082" iconserver="3642" iconfarm="4" member_count="55925" topic_count="176" pool_count="5310082" />
    <group nsid="940119@N23" id="940119@N23" name="La Febbra 2" member="1" moderator="0" admin="0" privacy="1" photos="118583" iconserver="4765" iconfarm="5" member_count="2935" topic_count="5" pool_count="118583" />
    <group nsid="715137@N23" id="715137@N23" name="***SCREAM OF THE PHOTOGRAPHER***" member="1" moderator="0" admin="0" privacy="3" photos="3313390" iconserver="5334" iconfarm="6" member_count="54316" topic_count="151" pool_count="3313390" />
    <group nsid="1414356@N24" id="1414356@N24" name="*Nude Not Rude-Lingerie*" member="1" moderator="0" admin="0" privacy="3" photos="44499" iconserver="3756" iconfarm="4" member_count="11698" topic_count="21" pool_count="44499" />
    <group nsid="479584@N25" id="479584@N25" name="*Sensual Art*" member="1" moderator="0" admin="0" privacy="3" photos="162392" iconserver="4804" iconfarm="5" member_count="17897" topic_count="21" pool_count="162392" />
    <group nsid="342582@N20" id="342582@N20" name="100 Strangers" member="1" moderator="0" admin="0" privacy="3" photos="50993" iconserver="7426" iconfarm="8" member_count="14228" topic_count="991" pool_count="50993" />
    <group nsid="14771020@N24" id="14771020@N24" name="50 Shades Of Gray" member="1" moderator="0" admin="0" privacy="3" photos="5043" iconserver="65535" iconfarm="66" member_count="638" topic_count="1" pool_count="5043" />
    <group nsid="3542605@N22" id="3542605@N22" name="85mm Prime Portraits" member="1" moderator="0" admin="0" privacy="3" photos="5325" iconserver="4307" iconfarm="5" member_count="372" topic_count="2" pool_count="5325" />
    <group nsid="2781755@N23" id="2781755@N23" name="[][][] The best pictures [][][]" member="1" moderator="0" admin="0" privacy="3" photos="680347" iconserver="7474" iconfarm="8" member_count="7674" topic_count="45" pool_count="680347" />
    <group nsid="3125076@N23" id="3125076@N23" name="_ Women _" member="1" moderator="0" admin="0" privacy="3" photos="37667" iconserver="583" iconfarm="1" member_count="1890" topic_count="2" pool_count="37667" />
    <group nsid="1755214@N23" id="1755214@N23" name=" FLICKR " member="1" moderator="0" admin="0" privacy="3" photos="3223430" iconserver="7287" iconfarm="8" member_count="35516" topic_count="66" pool_count="3223430" />
    <group nsid="74963390@N00" id="74963390@N00" name="A Lifetime of Self-Portraits" member="1" moderator="0" admin="0" privacy="3" photos="40824" iconserver="135" iconfarm="1" member_count="2557" topic_count="4" pool_count="40824" />
    <group nsid="1437812@N22" id="1437812@N22" name="A Portrait Group" member="1" moderator="0" admin="0" privacy="2" photos="551614" iconserver="65535" iconfarm="66" member_count="9066" topic_count="0" pool_count="551614" />
    <group nsid="3814553@N23" id="3814553@N23" name="A UNESCO SITE" member="1" moderator="0" admin="0" privacy="2" photos="63556" iconserver="65535" iconfarm="66" member_count="1230" topic_count="17" pool_count="63556" />
    <group nsid="14603619@N24" id="14603619@N24" name="A. Dangerous Perfume" member="1" moderator="0" admin="0" privacy="3" photos="86560" iconserver="65535" iconfarm="66" member_count="4583" topic_count="2" pool_count="86560" />
    <group nsid="2405724@N25" id="2405724@N25" name="A7 | A7R | A7S" member="1" moderator="0" admin="0" privacy="3" photos="142259" iconserver="3736" iconfarm="4" member_count="3466" topic_count="14" pool_count="142259" />
    <group nsid="92938852@N00" id="92938852@N00" name="Adobe Lightroom" member="1" moderator="0" admin="0" privacy="3" photos="79160" iconserver="3679" iconfarm="4" member_count="26557" topic_count="6784" pool_count="79160" />
    <group nsid="43501458@N00" id="43501458@N00" name="Amateurs" member="1" moderator="0" admin="0" privacy="3" photos="8143839" iconserver="7040" iconfarm="8" member_count="94132" topic_count="774" pool_count="8143839" />
    <group nsid="70108624@N00" id="70108624@N00" name="Amatuer Photography" member="1" moderator="0" admin="0" privacy="3" photos="275980" iconserver="0" iconfarm="0" member_count="2858" topic_count="113" pool_count="275980" />
    <group nsid="903702@N21" id="903702@N21" name="Ambitious Photographers" member="1" moderator="0" admin="0" privacy="3" photos="350363" iconserver="3164" iconfarm="4" member_count="6318" topic_count="29" pool_count="350363" />
    <group nsid="2113006@N20" id="2113006@N20" name="Friends of Zanskar, Ladakh, himalayans countries and Asia !" member="1" moderator="0" admin="0" privacy="2" photos="3417" iconserver="8483" iconfarm="9" member_count="82" topic_count="1" pool_count="3417" />
    <group nsid="76649513@N00" id="76649513@N00" name="ANYTHING ALLOWED!!!!!!" member="1" moderator="0" admin="0" privacy="3" photos="4916803" iconserver="26" iconfarm="1" member_count="32233" topic_count="108" pool_count="4916803" />
    <group nsid="1877553@N25" id="1877553@N25" name="!Art Book (18+)" member="1" moderator="0" admin="0" privacy="2" photos="189850" iconserver="65535" iconfarm="66" member_count="4695" topic_count="4" pool_count="189850" />
    <group nsid="892682@N22" id="892682@N22" name="Art of Images ~ AOI L1~ ( Post 1, Award 3)  NEW AWARD..!" member="1" moderator="0" admin="0" privacy="2" photos="1149453" iconserver="65535" iconfarm="66" member_count="38170" topic_count="92" pool_count="1149453" />
    <group nsid="2622211@N25" id="2622211@N25" name="Artistic Nudes From Behind" member="1" moderator="0" admin="0" privacy="2" photos="30043" iconserver="661" iconfarm="1" member_count="1601" topic_count="3" pool_count="30043" />
    <group nsid="58898522@N00" id="58898522@N00" name="Artistic Photography" member="1" moderator="0" admin="0" privacy="3" photos="6187891" iconserver="65535" iconfarm="66" member_count="142002" topic_count="389" pool_count="6187891" />
    <group nsid="86193961@N00" id="86193961@N00" name="Artistic Photography of Nude Women" member="1" moderator="0" admin="0" privacy="3" photos="114654" iconserver="47" iconfarm="1" member_count="11837" topic_count="31" pool_count="114654" />
    <group nsid="1046032@N20" id="1046032@N20" name="as beautiful as you want" member="1" moderator="0" admin="0" privacy="3" photos="2604060" iconserver="7896" iconfarm="8" member_count="24511" topic_count="136" pool_count="2604060" />
    <group nsid="3678915@N22" id="3678915@N22" name="asian provate pix by requests" member="1" moderator="0" admin="0" privacy="3" photos="595" iconserver="4339" iconfarm="5" member_count="975" topic_count="7" pool_count="595" />
    <group nsid="442480@N23" id="442480@N23" name="Available Light / Existing Light" member="1" moderator="0" admin="0" privacy="3" photos="1666599" iconserver="5329" iconfarm="6" member_count="25691" topic_count="33" pool_count="1666599" />
    <group nsid="90122202@N00" id="90122202@N00" name="Bangalore Shutterbugs Group" member="1" moderator="0" admin="0" privacy="3" photos="51625" iconserver="32" iconfarm="1" member_count="3365" topic_count="1135" pool_count="51625" />
    <group nsid="352933@N20" id="352933@N20" name="Bangalore Travel Photographers" member="1" moderator="0" admin="0" privacy="3" photos="50262" iconserver="205" iconfarm="1" member_count="3003" topic_count="468" pool_count="50262" />
    <group nsid="14809765@N21" id="14809765@N21" name="Beating Around The Bush - Beautiful Natural Women" member="1" moderator="0" admin="0" privacy="3" photos="2280" iconserver="65535" iconfarm="66" member_count="1945" topic_count="0" pool_count="2280" />
    <group nsid="14747783@N25" id="14747783@N25" name="Beautiful Images -- all kinds" member="1" moderator="0" admin="0" privacy="3" photos="20549" iconserver="0" iconfarm="0" member_count="929" topic_count="2" pool_count="20549" />
    <group nsid="887742@N25" id="887742@N25" name="*•●♥  Beautiful  Women  ♥●•*" member="1" moderator="0" admin="0" privacy="3" photos="239444" iconserver="3010" iconfarm="4" member_count="9529" topic_count="3" pool_count="239444" />
    <group nsid="1760479@N25" id="1760479@N25" name="Beautiful Woman (high quality photos)" member="1" moderator="0" admin="0" privacy="3" photos="649372" iconserver="6043" iconfarm="7" member_count="3833" topic_count="2" pool_count="649372" />
    <group nsid="14755381@N25" id="14755381@N25" name="Beautiful Woman Portraits with Pretty Faces (NO 30/60)" member="1" moderator="0" admin="0" privacy="3" photos="16587" iconserver="65535" iconfarm="66" member_count="352" topic_count="2" pool_count="16587" />
    <group nsid="79331465@N00" id="79331465@N00" name="Beautiful Women (no 2nd life and no CD/TV !)" member="1" moderator="0" admin="0" privacy="3" photos="607707" iconserver="65535" iconfarm="66" member_count="33736" topic_count="207" pool_count="607707" />
    <group nsid="1893995@N24" id="1893995@N24" name="Beautiful Women in Tight Outfits and Dresses" member="1" moderator="0" admin="0" privacy="3" photos="155249" iconserver="5116" iconfarm="6" member_count="9100" topic_count="3" pool_count="155249" />
    <group nsid="1869549@N21" id="1869549@N21" name="Beautiful Young Women and Their Faces" member="1" moderator="0" admin="0" privacy="3" photos="87153" iconserver="702" iconfarm="1" member_count="3814" topic_count="3" pool_count="87153" />
    <group nsid="1956087@N22" id="1956087@N22" name="Feminine Beauty 2" member="1" moderator="0" admin="0" privacy="3" photos="786160" iconserver="65535" iconfarm="66" member_count="20951" topic_count="0" pool_count="786160" />
    <group nsid="94899009@N00" id="94899009@N00" name="Beginner's Digital Photography" member="1" moderator="0" admin="0" privacy="3" photos="1078508" iconserver="4036" iconfarm="5" member_count="39679" topic_count="3412" pool_count="1078508" />
    <group nsid="2920545@N23" id="2920545@N23" name="Bellas Ellas" member="1" moderator="0" admin="0" privacy="3" photos="436133" iconserver="65535" iconfarm="66" member_count="15391" topic_count="0" pool_count="436133" />
    <group nsid="14802474@N21" id="14802474@N21" name="Ben's absolute erotic favs" member="1" moderator="0" admin="0" privacy="3" photos="18807" iconserver="65535" iconfarm="66" member_count="2918" topic_count="3" pool_count="18807" />
    <group nsid="4117803@N25" id="4117803@N25" name="Best nude art" member="1" moderator="0" admin="0" privacy="3" photos="114784" iconserver="7828" iconfarm="8" member_count="6551" topic_count="35" pool_count="114784" />
    <group nsid="14661711@N22" id="14661711@N22" name="Beyond the lens..." member="1" moderator="0" admin="0" privacy="1" photos="16110" iconserver="65535" iconfarm="66" member_count="89" topic_count="0" pool_count="16110" />
    <group nsid="1343753@N23" id="1343753@N23" name="Birth of Venus" member="1" moderator="0" admin="0" privacy="3" photos="76545" iconserver="4016" iconfarm="5" member_count="8625" topic_count="4" pool_count="76545" />
    <group nsid="1881703@N25" id="1881703@N25" name="Black &amp; White Beauties" member="1" moderator="0" admin="0" privacy="3" photos="81064" iconserver="2818" iconfarm="3" member_count="5558" topic_count="13" pool_count="81064" />
    <group nsid="29184631@N00" id="29184631@N00" name="Black &amp; White Portraits" member="1" moderator="0" admin="0" privacy="3" photos="433781" iconserver="65535" iconfarm="66" member_count="34091" topic_count="100" pool_count="433781" />
    <group nsid="16978849@N00" id="16978849@N00" name="Black and White" member="1" moderator="0" admin="0" privacy="3" photos="5257525" iconserver="65535" iconfarm="66" member_count="352654" topic_count="2789" pool_count="5257525" />
    <group nsid="55296031@N00" id="55296031@N00" name="BOKEH - for the common folk" member="1" moderator="0" admin="0" privacy="3" photos="354296" iconserver="66" iconfarm="1" member_count="9944" topic_count="48" pool_count="354296" />
    <group nsid="646487@N23" id="646487@N23" name="Bokeh addiction" member="1" moderator="0" admin="0" privacy="3" photos="512209" iconserver="2166" iconfarm="3" member_count="14245" topic_count="14" pool_count="512209" />
    <group nsid="692285@N25" id="692285@N25" name="Bokeh Of The Day" member="1" moderator="0" admin="0" privacy="3" photos="229242" iconserver="2034" iconfarm="3" member_count="7485" topic_count="31" pool_count="229242" />
    <group nsid="38343303@N00" id="38343303@N00" name="Bokeh: Smooth &amp; Silky" member="1" moderator="0" admin="0" privacy="3" photos="1630484" iconserver="36" iconfarm="1" member_count="74767" topic_count="485" pool_count="1630484" />
    <group nsid="4214058@N22" id="4214058@N22" name="Boudoir and sensuality" member="1" moderator="0" admin="0" privacy="3" photos="20874" iconserver="65535" iconfarm="66" member_count="2496" topic_count="6" pool_count="20874" />
    <group nsid="2441147@N23" id="2441147@N23" name="Breasts Beautiful  (No pornography)" member="1" moderator="0" admin="0" privacy="3" photos="146005" iconserver="1625" iconfarm="2" member_count="12757" topic_count="0" pool_count="146005" />
    <group nsid="4029513@N25" id="4029513@N25" name="BREASTS: topless, cleavage, downblouse, clothed very attractive" member="1" moderator="0" admin="0" privacy="0" photos="74314" iconserver="4440" iconfarm="5" member_count="3633" topic_count="6" pool_count="74314" />
    <group nsid="84249994@N00" id="84249994@N00" name="Canon 350D/400D/450D/500D/550D/600D/650D &amp; 700D" member="1" moderator="0" admin="0" privacy="3" photos="624554" iconserver="71" iconfarm="1" member_count="15754" topic_count="456" pool_count="624554" />
    <group nsid="2928429@N21" id="2928429@N21" name="Canon 50mm f/1.8 f/1.4 f/1.2 f/1.0" member="1" moderator="0" admin="0" privacy="3" photos="127879" iconserver="5708" iconfarm="6" member_count="6839" topic_count="35" pool_count="127879" />
    <group nsid="35034359018@N01" id="35034359018@N01" name="Canon DSLR User Group" member="1" moderator="0" admin="0" privacy="3" photos="5535256" iconserver="5" iconfarm="1" member_count="189453" topic_count="19216" pool_count="5535256" />
    <group nsid="35353093@N00" id="35353093@N00" name="Canon EOS Digital" member="1" moderator="0" admin="0" privacy="3" photos="420717" iconserver="31" iconfarm="1" member_count="3754" topic_count="40" pool_count="420717" />
    <group nsid="52240293230@N01" id="52240293230@N01" name="canon photography" member="1" moderator="0" admin="0" privacy="3" photos="7868396" iconserver="9" iconfarm="1" member_count="113033" topic_count="1429" pool_count="7868396" />
    <group nsid="476233@N23" id="476233@N23" name="Canon Users - Open group" member="1" moderator="0" admin="0" privacy="3" photos="270437" iconserver="1176" iconfarm="2" member_count="3703" topic_count="30" pool_count="270437" />
    <group nsid="1490049@N20" id="1490049@N20" name="closed ---" member="1" moderator="0" admin="0" privacy="3" photos="135949" iconserver="4147" iconfarm="5" member_count="10525" topic_count="13" pool_count="135949" />
    <group nsid="2937771@N20" id="2937771@N20" name="DAMN that's HOT!" member="1" moderator="0" admin="0" privacy="0" photos="68545" iconserver="1595" iconfarm="2" member_count="2177" topic_count="12" pool_count="68545" />
    <group nsid="2543356@N24" id="2543356@N24" name="Eternal Womanhood" member="1" moderator="0" admin="0" privacy="3" photos="208010" iconserver="1511" iconfarm="2" member_count="9255" topic_count="7" pool_count="208010" />
    <group nsid="311760@N20" id="311760@N20" name="Decisive Moments... Classics in B&amp;W" member="1" moderator="0" admin="0" privacy="3" photos="15444" iconserver="175" iconfarm="1" member_count="3841" topic_count="29" pool_count="15444" />
    <group nsid="33456739@N00" id="33456739@N00" name="Digital Photography School" member="1" moderator="0" admin="0" privacy="3" photos="825849" iconserver="78" iconfarm="1" member_count="21018" topic_count="725" pool_count="825849" />
    <group nsid="1567384@N23" id="1567384@N23" name="Discovery  India" member="1" moderator="0" admin="0" privacy="3" photos="143605" iconserver="4023" iconfarm="5" member_count="3156" topic_count="18" pool_count="143605" />
    <group nsid="1445159@N25" id="1445159@N25" name="Earth Air Atmosphere Nude (pure nudity )" member="1" moderator="0" admin="0" privacy="2" photos="36992" iconserver="855" iconfarm="1" member_count="1628" topic_count="15" pool_count="36992" />
    <group nsid="66433063@N00" id="66433063@N00" name="EF 50mm 1:1.8 II" member="1" moderator="0" admin="0" privacy="3" photos="107338" iconserver="53" iconfarm="1" member_count="7762" topic_count="124" pool_count="107338" />
    <group nsid="14767721@N24" id="14767721@N24" name="Enchanting, Mesmerizing, Spellbinding Women" member="1" moderator="0" admin="0" privacy="0" photos="33679" iconserver="65535" iconfarm="66" member_count="460" topic_count="1" pool_count="33679" />
    <group nsid="14833125@N22" id="14833125@N22" name="Erotic Art - really the best" member="1" moderator="0" admin="0" privacy="3" photos="6272" iconserver="65535" iconfarm="66" member_count="1254" topic_count="1" pool_count="6272" />
    <group nsid="2846476@N22" id="2846476@N22" name="EROTIC ART, EROTIC NUDE (NO PORN !!!)" member="1" moderator="0" admin="0" privacy="3" photos="55722" iconserver="65535" iconfarm="66" member_count="2430" topic_count="13" pool_count="55722" />
    <group nsid="2621856@N23" id="2621856@N23" name="Erotic Travel (No porn !!! No men !!!)" member="1" moderator="0" admin="0" privacy="3" photos="116959" iconserver="7697" iconfarm="8" member_count="6915" topic_count="1" pool_count="116959" />
    <group nsid="1494581@N21" id="1494581@N21" name="Exotic Bodies - Hi Grade" member="1" moderator="0" admin="0" privacy="3" photos="97396" iconserver="4625" iconfarm="5" member_count="12913" topic_count="10" pool_count="97396" />
    <group nsid="2817051@N20" id="2817051@N20" name="Exquisite Creatures" member="1" moderator="0" admin="0" privacy="2" photos="92977" iconserver="7761" iconfarm="8" member_count="1365" topic_count="0" pool_count="92977" />
    <group nsid="3448596@N23" id="3448596@N23" name="Exquisite Erotism" member="1" moderator="0" admin="0" privacy="3" photos="37673" iconserver="7807" iconfarm="8" member_count="16174" topic_count="93" pool_count="37673" />
    <group nsid="938655@N20" id="938655@N20" name="Exquisite Erotism Outdoor" member="1" moderator="0" admin="0" privacy="3" photos="19017" iconserver="2932" iconfarm="3" member_count="4779" topic_count="5" pool_count="19017" />
    <group nsid="3190486@N22" id="3190486@N22" name="Exquisite Lingerie" member="1" moderator="0" admin="0" privacy="3" photos="22860" iconserver="7896" iconfarm="8" member_count="4103" topic_count="8" pool_count="22860" />
    <group nsid="2971446@N23" id="2971446@N23" name="Exquisite Skirts and shorts" member="1" moderator="0" admin="0" privacy="3" photos="13735" iconserver="3938" iconfarm="4" member_count="2960" topic_count="4" pool_count="13735" />
    <group nsid="2704186@N24" id="2704186@N24" name="Faces that can launch 1000 ships - No 30/60" member="1" moderator="0" admin="0" privacy="3" photos="88465" iconserver="65535" iconfarm="66" member_count="2455" topic_count="1" pool_count="88465" />
    <group nsid="1439799@N21" id="1439799@N21" name="Fascinating Women in Black &amp; White Art" member="1" moderator="0" admin="0" privacy="3" photos="17236" iconserver="3752" iconfarm="4" member_count="1584" topic_count="4" pool_count="17236" />
    <group nsid="64332025@N00" id="64332025@N00" name="FASHION PHOTOGRAPHY" member="1" moderator="0" admin="0" privacy="3" photos="843149" iconserver="65535" iconfarm="66" member_count="54353" topic_count="1628" pool_count="843149" />
    <group nsid="1552090@N23" id="1552090@N23" name="Female Loveliness" member="1" moderator="0" admin="0" privacy="3" photos="62118" iconserver="7371" iconfarm="8" member_count="3746" topic_count="0" pool_count="62118" />
    <group nsid="61237115@N00" id="61237115@N00" name="Female Models" member="1" moderator="0" admin="0" privacy="3" photos="568975" iconserver="65535" iconfarm="66" member_count="24455" topic_count="119" pool_count="568975" />
    <group nsid="32177968@N00" id="32177968@N00" name="Feminine Mystique" member="1" moderator="0" admin="0" privacy="3" photos="309419" iconserver="65535" iconfarm="66" member_count="16136" topic_count="84" pool_count="309419" />
    <group nsid="1506929@N20" id="1506929@N20" name="feminine petite women" member="1" moderator="0" admin="0" privacy="2" photos="333518" iconserver="7319" iconfarm="8" member_count="19106" topic_count="6" pool_count="333518" />
    <group nsid="76535076@N00" id="76535076@N00" name="Flickr Addicts" member="1" moderator="0" admin="0" privacy="3" photos="10096107" iconserver="65535" iconfarm="66" member_count="139206" topic_count="206" pool_count="10096107" />
    <group nsid="2641619@N21" id="2641619@N21" name="Flickr Alpha Group" member="1" moderator="0" admin="0" privacy="1" photos="79" iconserver="3826" iconfarm="4" member_count="3757" topic_count="498" pool_count="79" />
    <group nsid="1148171@N20" id="1148171@N20" name="Flickr Community (INVITE ALL CONTACTS)" member="1" moderator="0" admin="0" privacy="3" photos="1741180" iconserver="2553" iconfarm="3" member_count="44231" topic_count="216" pool_count="1741180" />
    <group nsid="2641368@N20" id="2641368@N20" name="Flickr Sensual, Beautiful and Artistic Scented Woman" member="1" moderator="0" admin="0" privacy="3" photos="145124" iconserver="1514" iconfarm="2" member_count="6222" topic_count="4" pool_count="145124" />
    <group nsid="11405564@N00" id="11405564@N00" name="Flickr's Finest Female Thongs" member="1" moderator="0" admin="0" privacy="3" photos="1630" iconserver="56" iconfarm="1" member_count="2665" topic_count="19" pool_count="1630" />
    <group nsid="34427469792@N01" id="34427469792@N01" name="FlickrCentral" member="1" moderator="0" admin="0" privacy="3" photos="8594732" iconserver="2871" iconfarm="3" member_count="314331" topic_count="12854" pool_count="8594732" />
    <group nsid="95309787@N00" id="95309787@N00" name="Flickritis: Where 86,000 people are still happily infected!" member="1" moderator="0" admin="0" privacy="3" photos="10100000" iconserver="18" iconfarm="1" member_count="87711" topic_count="1490" pool_count="10100000" />
    <group nsid="38436807@N00" id="38436807@N00" name="FlickrToday (only 1 pic per day)" member="1" moderator="0" admin="0" privacy="3" photos="10100000" iconserver="32" iconfarm="1" member_count="167005" topic_count="890" pool_count="10100000" />
    <group nsid="735470@N20" id="735470@N20" name="Flowers Planet Photography" member="1" moderator="0" admin="0" privacy="3" photos="262674" iconserver="3152" iconfarm="4" member_count="18615" topic_count="14" pool_count="262674" />
    <group nsid="2126006@N25" id="2126006@N25" name="Fresh People" member="1" moderator="0" admin="0" privacy="2" photos="159321" iconserver="7059" iconfarm="8" member_count="2471" topic_count="0" pool_count="159321" />
    <group nsid="64909868@N00" id="64909868@N00" name="From the Airplane Window" member="1" moderator="0" admin="0" privacy="3" photos="35390" iconserver="98" iconfarm="1" member_count="4003" topic_count="4" pool_count="35390" />
    <group nsid="2744319@N23" id="2744319@N23" name="Garden of Eden (Female nude)" member="1" moderator="0" admin="0" privacy="2" photos="25939" iconserver="4675" iconfarm="5" member_count="1236" topic_count="6" pool_count="25939" />
    <group nsid="517670@N25" id="517670@N25" name="Girls Wearing Tights" member="1" moderator="0" admin="0" privacy="3" photos="1487" iconserver="0" iconfarm="0" member_count="1276" topic_count="4" pool_count="1487" />
    <group nsid="71308295@N00" id="71308295@N00" name="Glamour Photography [ No Serials, Plz Read ]" member="1" moderator="0" admin="0" privacy="3" photos="90790" iconserver="15" iconfarm="1" member_count="8498" topic_count="59" pool_count="90790" />
    <group nsid="38132155@N00" id="38132155@N00" name="Gorgeous" member="1" moderator="0" admin="0" privacy="3" photos="163649" iconserver="7" iconfarm="1" member_count="7986" topic_count="35" pool_count="163649" />
    <group nsid="14750720@N25" id="14750720@N25" name="Gorgeous Faces and Breasts" member="1" moderator="0" admin="0" privacy="0" photos="18988" iconserver="65535" iconfarm="66" member_count="488" topic_count="2" pool_count="18988" />
    <group nsid="14674002@N25" id="14674002@N25" name="Gorgeous Women / Celebrating Beautiful Ladies" member="1" moderator="0" admin="0" privacy="1" photos="41239" iconserver="65535" iconfarm="66" member_count="556" topic_count="1" pool_count="41239" />
    <group nsid="14703237@N25" id="14703237@N25" name="Gorgeous, Sensual, Enticing, Alluring Fashion Showcasing Women" member="1" moderator="0" admin="0" privacy="0" photos="67606" iconserver="65535" iconfarm="66" member_count="1645" topic_count="3" pool_count="67606" />
    <group nsid="1676686@N24" id="1676686@N24" name="Great Photography in Action - photoExplode.com" member="1" moderator="0" admin="0" privacy="3" photos="103416" iconserver="3023" iconfarm="4" member_count="1357" topic_count="12" pool_count="103416" />
    <group nsid="2575977@N21" id="2575977@N21" name="Great Portraits" member="1" moderator="0" admin="0" privacy="3" photos="332084" iconserver="65535" iconfarm="66" member_count="6907" topic_count="18" pool_count="332084" />
    <group nsid="94761711@N00" id="94761711@N00" name="HCSP (Hardcore Street Photography)" member="1" moderator="0" admin="0" privacy="3" photos="3583" iconserver="7334" iconfarm="8" member_count="88927" topic_count="3195" pool_count="3583" />
    <group nsid="569830@N25" id="569830@N25" name="HDR-Imaging" member="1" moderator="0" admin="0" privacy="3" photos="709103" iconserver="8638" iconfarm="9" member_count="28478" topic_count="175" pool_count="709103" />
    <group nsid="1304786@N23" id="1304786@N23" name="HDR Photo" member="1" moderator="0" admin="0" privacy="3" photos="210479" iconserver="2587" iconfarm="3" member_count="7581" topic_count="12" pool_count="210479" />
    <group nsid="1604876@N25" id="1604876@N25" name="HDR WORLDWIDE" member="1" moderator="0" admin="0" privacy="3" photos="229767" iconserver="5052" iconfarm="6" member_count="7864" topic_count="8" pool_count="229767" />
    <group nsid="80429107@N00" id="80429107@N00" name="Historic India" member="1" moderator="0" admin="0" privacy="3" photos="22040" iconserver="4877" iconfarm="5" member_count="2898" topic_count="53" pool_count="22040" />
    <group nsid="25794387@N00" id="25794387@N00" name="Home Studio" member="1" moderator="0" admin="0" privacy="3" photos="49189" iconserver="50" iconfarm="1" member_count="4483" topic_count="5" pool_count="49189" />
    <group nsid="12173949@N00" id="12173949@N00" name="I Love My Canon" member="1" moderator="0" admin="0" privacy="3" photos="1298107" iconserver="0" iconfarm="0" member_count="10551" topic_count="35" pool_count="1298107" />
    <group nsid="865297@N20" id="865297@N20" name="Oriental Land" member="1" moderator="0" admin="0" privacy="3" photos="126608" iconserver="3771" iconfarm="4" member_count="2583" topic_count="24" pool_count="126608" />
    <group nsid="2389839@N23" id="2389839@N23" name="in explore" member="1" moderator="0" admin="0" privacy="3" photos="1572334" iconserver="3744" iconfarm="4" member_count="100169" topic_count="721" pool_count="1572334" />
    <group nsid="79356248@N00" id="79356248@N00" name="Incredible India !" member="1" moderator="0" admin="0" privacy="3" photos="51073" iconserver="1078" iconfarm="2" member_count="4536" topic_count="46" pool_count="51073" />
    <group nsid="90544387@N00" id="90544387@N00" name="India and Indians" member="1" moderator="0" admin="0" privacy="3" photos="324259" iconserver="8683" iconfarm="9" member_count="9411" topic_count="117" pool_count="324259" />
    <group nsid="52992818@N00" id="52992818@N00" name="India as we see it" member="1" moderator="0" admin="0" privacy="3" photos="184393" iconserver="30" iconfarm="1" member_count="3702" topic_count="23" pool_count="184393" />
    <group nsid="52242377700@N01" id="52242377700@N01" name="India Images" member="1" moderator="0" admin="0" privacy="3" photos="516395" iconserver="1" iconfarm="1" member_count="20148" topic_count="151" pool_count="516395" />
    <group nsid="2901184@N23" id="2901184@N23" name="INDO INDIAN RELICS" member="1" moderator="0" admin="0" privacy="3" photos="12751" iconserver="725" iconfarm="1" member_count="970" topic_count="3" pool_count="12751" />
    <group nsid="2657214@N25" id="2657214@N25" name="INSPIRATIO (WOMEN ONLY)" member="1" moderator="0" admin="0" privacy="2" photos="55530" iconserver="7845" iconfarm="8" member_count="1488" topic_count="11" pool_count="55530" />
    <group nsid="76984205@N00" id="76984205@N00" name="Interesting Portraits ☺ [people]" member="1" moderator="0" admin="0" privacy="3" photos="851193" iconserver="28" iconfarm="1" member_count="28539" topic_count="61" pool_count="851193" />
    <group nsid="4047524@N23" id="4047524@N23" name="INVITATION: Seductive, Alluring, Flirtatious Women" member="1" moderator="0" admin="0" privacy="0" photos="118713" iconserver="7802" iconfarm="8" member_count="3366" topic_count="5" pool_count="118713" />
    <group nsid="796830@N22" id="796830@N22" name="Just Clouds P1/ A2 ** INVITE YOUR FRIENDS!!**Post 1 Award 2." member="1" moderator="0" admin="0" privacy="3" photos="26220" iconserver="2076" iconfarm="3" member_count="3090" topic_count="32" pool_count="26220" />
    <group nsid="1619957@N24" id="1619957@N24" name="Le Poet" member="1" moderator="0" admin="0" privacy="1" photos="78853" iconserver="5309" iconfarm="6" member_count="2815" topic_count="23" pool_count="78853" />
    <group nsid="86285131@N00" id="86285131@N00" name="Life Through My Viewfinder" member="1" moderator="0" admin="0" privacy="3" photos="232217" iconserver="150" iconfarm="1" member_count="3669" topic_count="19" pool_count="232217" />
    <group nsid="988940@N24" id="988940@N24" name="Lightroom Showcase" member="1" moderator="0" admin="0" privacy="3" photos="106717" iconserver="3387" iconfarm="4" member_count="3502" topic_count="0" pool_count="106717" />
    <group nsid="93318005@N00" id="93318005@N00" name="Lingerie" member="1" moderator="0" admin="0" privacy="3" photos="23761" iconserver="23" iconfarm="1" member_count="8405" topic_count="85" pool_count="23761" />
    <group nsid="1725800@N23" id="1725800@N23" name="_The Best Asses (female)" member="1" moderator="0" admin="0" privacy="3" photos="28272" iconserver="5772" iconfarm="6" member_count="7042" topic_count="7" pool_count="28272" />
    <group nsid="12901672@N00" id="12901672@N00" name="Love and only love" member="1" moderator="0" admin="0" privacy="1" photos="55170" iconserver="4721" iconfarm="5" member_count="3328" topic_count="42" pool_count="55170" />
    <group nsid="52240880306@N01" id="52240880306@N01" name="Magic Moments" member="1" moderator="0" admin="0" privacy="3" photos="3084479" iconserver="65535" iconfarm="66" member_count="63212" topic_count="110" pool_count="3084479" />
    <group nsid="2841338@N20" id="2841338@N20" name="MasterClub" member="1" moderator="0" admin="0" privacy="2" photos="82167" iconserver="5705" iconfarm="6" member_count="1201" topic_count="42" pool_count="82167" />
    <group nsid="4215978@N22" id="4215978@N22" name="Maximum Female Beauty" member="1" moderator="0" admin="0" privacy="2" photos="99293" iconserver="65535" iconfarm="66" member_count="1380" topic_count="2" pool_count="99293" />
    <group nsid="35468145766@N01" id="35468145766@N01" name="Models" member="1" moderator="0" admin="0" privacy="3" photos="622470" iconserver="883" iconfarm="1" member_count="24175" topic_count="174" pool_count="622470" />
    <group nsid="1208049@N24" id="1208049@N24" name="Models - Only high quality photos" member="1" moderator="0" admin="0" privacy="1" photos="124098" iconserver="65535" iconfarm="66" member_count="5087" topic_count="1" pool_count="124098" />
    <group nsid="52240914900@N01" id="52240914900@N01" name="MOON Shots" member="1" moderator="0" admin="0" privacy="3" photos="165988" iconserver="1" iconfarm="1" member_count="27457" topic_count="159" pool_count="165988" />
    <group nsid="1583974@N22" id="1583974@N22" name="Mordus de photos Studio" member="1" moderator="0" admin="0" privacy="3" photos="188750" iconserver="3665" iconfarm="4" member_count="7262" topic_count="23" pool_count="188750" />
    <group nsid="2369838@N20" id="2369838@N20" name="My Private Favorite Photo" member="1" moderator="0" admin="0" privacy="1" photos="12962" iconserver="7354" iconfarm="8" member_count="482" topic_count="4" pool_count="12962" />
    <group nsid="14816895@N22" id="14816895@N22" name="Nature désirée" member="1" moderator="0" admin="0" privacy="3" photos="11225" iconserver="65535" iconfarm="66" member_count="5089" topic_count="82" pool_count="11225" />
    <group nsid="81127001@N00" id="81127001@N00" name="Nice pictures" member="1" moderator="0" admin="0" privacy="3" photos="4562950" iconserver="2891" iconfarm="3" member_count="65222" topic_count="185" pool_count="4562950" />
    <group nsid="44266388@N00" id="44266388@N00" name="Nude Amateur Models" member="1" moderator="0" admin="0" privacy="3" photos="125593" iconserver="0" iconfarm="0" member_count="25043" topic_count="248" pool_count="125593" />
    <group nsid="2141464@N23" id="2141464@N23" name="Nude Photographers" member="1" moderator="0" admin="0" privacy="2" photos="97281" iconserver="8049" iconfarm="9" member_count="3749" topic_count="13" pool_count="97281" />
    <group nsid="3014196@N22" id="3014196@N22" name="One rule : WOMEN" member="1" moderator="0" admin="0" privacy="3" photos="528718" iconserver="7108" iconfarm="8" member_count="8461" topic_count="31" pool_count="528718" />
    <group nsid="76071845@N00" id="76071845@N00" name="One Strobe Pony" member="1" moderator="0" admin="0" privacy="3" photos="84388" iconserver="3050" iconfarm="4" member_count="9101" topic_count="148" pool_count="84388" />
    <group nsid="1927900@N22" id="1927900@N22" name="non-nude erotica" member="1" moderator="0" admin="0" privacy="3" photos="22863" iconserver="8840" iconfarm="9" member_count="2460" topic_count="2" pool_count="22863" />
    <group nsid="894842@N25" id="894842@N25" name="Paradise by the photographers light" member="1" moderator="0" admin="0" privacy="3" photos="534360" iconserver="65535" iconfarm="66" member_count="16276" topic_count="3" pool_count="534360" />
    <group nsid="12938526@N00" id="12938526@N00" name="People / Persons - + Photo Titles" member="1" moderator="0" admin="0" privacy="3" photos="955239" iconserver="65535" iconfarm="66" member_count="32275" topic_count="39" pool_count="955239" />
    <group nsid="21819324@N00" id="21819324@N00" name="People of the World" member="1" moderator="0" admin="0" privacy="3" photos="1043458" iconserver="4" iconfarm="1" member_count="24747" topic_count="77" pool_count="1043458" />
    <group nsid="89483931@N00" id="89483931@N00" name="People pictures" member="1" moderator="0" admin="0" privacy="3" photos="327518" iconserver="134" iconfarm="1" member_count="6259" topic_count="4" pool_count="327518" />
    <group nsid="1812051@N23" id="1812051@N23" name="People Portraits World" member="1" moderator="0" admin="0" privacy="3" photos="231862" iconserver="5828" iconfarm="6" member_count="9091" topic_count="12" pool_count="231862" />
    <group nsid="14687468@N25" id="14687468@N25" name="People's Choice: Women - Choose Five per Day - All Accepted" member="1" moderator="0" admin="0" privacy="0" photos="35482" iconserver="65535" iconfarm="66" member_count="425" topic_count="1" pool_count="35482" />
    <group nsid="1584087@N20" id="1584087@N20" name="Photographer Universe" member="1" moderator="0" admin="0" privacy="3" photos="1056035" iconserver="2924" iconfarm="3" member_count="8991" topic_count="56" pool_count="1056035" />
    <group nsid="2755427@N24" id="2755427@N24" name="Photographer's Fine Arts" member="1" moderator="0" admin="0" privacy="1" photos="34764" iconserver="8792" iconfarm="9" member_count="183" topic_count="2" pool_count="34764" />
    <group nsid="40732569271@N01" id="40732569271@N01" name="Photography" member="1" moderator="0" admin="0" privacy="3" photos="6702680" iconserver="65535" iconfarm="66" member_count="77324" topic_count="332" pool_count="6702680" />
    <group nsid="58184182@N00" id="58184182@N00" name="Photoshop Support Group" member="1" moderator="0" admin="0" privacy="2" photos="353468" iconserver="5448" iconfarm="6" member_count="34978" topic_count="8722" pool_count="353468" />
    <group nsid="52239745968@N01" id="52239745968@N01" name="Portrait" member="1" moderator="0" admin="0" privacy="3" photos="2315606" iconserver="65535" iconfarm="66" member_count="106115" topic_count="766" pool_count="2315606" />
    <group nsid="1359595@N22" id="1359595@N22" name="Portrait  World (人像世界)" member="1" moderator="0" admin="0" privacy="3" photos="399874" iconserver="65535" iconfarm="66" member_count="9443" topic_count="8" pool_count="399874" />
    <group nsid="3945515@N23" id="3945515@N23" name="Portrait (Fine Art, Cinematic and Conceptual) -- 肖像（藝術" member="1" moderator="0" admin="0" privacy="3" photos="12255" iconserver="65535" iconfarm="66" member_count="1737" topic_count="0" pool_count="12255" />
    <group nsid="1267409@N20" id="1267409@N20" name="Portrait Beauty !" member="1" moderator="0" admin="0" privacy="2" photos="248012" iconserver="5479" iconfarm="6" member_count="9364" topic_count="43" pool_count="248012" />
    <group nsid="3599274@N21" id="3599274@N21" name="Portrait by Invitation" member="1" moderator="0" admin="0" privacy="3" photos="7797" iconserver="65535" iconfarm="66" member_count="173" topic_count="6" pool_count="7797" />
    <group nsid="37981909@N00" id="37981909@N00" name="Portrait of FACE &lt; Head to Shoulder Shots Only &gt;" member="1" moderator="0" admin="0" privacy="3" photos="187132" iconserver="65535" iconfarm="66" member_count="11028" topic_count="37" pool_count="187132" />
    <group nsid="51772208@N00" id="51772208@N00" name="Portrait Portrait Portrait (Read THE RULE).....yes because there" member="1" moderator="0" admin="0" privacy="3" photos="184452" iconserver="91" iconfarm="1" member_count="9235" topic_count="20" pool_count="184452" />
    <group nsid="991969@N20" id="991969@N20" name="Portrait Strobist" member="1" moderator="0" admin="0" privacy="3" photos="75911" iconserver="3517" iconfarm="4" member_count="4089" topic_count="16" pool_count="75911" />
    <group nsid="388485@N22" id="388485@N22" name="Portraits 'n Fashion" member="1" moderator="0" admin="0" privacy="3" photos="190745" iconserver="177" iconfarm="1" member_count="8561" topic_count="28" pool_count="190745" />
    <group nsid="1528452@N23" id="1528452@N23" name="Portraits (Individuality, Imagination, Originality)" member="1" moderator="0" admin="0" privacy="3" photos="435397" iconserver="4127" iconfarm="5" member_count="52233" topic_count="71" pool_count="435397" />
    <group nsid="2172681@N22" id="2172681@N22" name="PORTRAITS ***" member="1" moderator="0" admin="0" privacy="3" photos="503001" iconserver="7386" iconfarm="8" member_count="19724" topic_count="0" pool_count="503001" />
    <group nsid="33945696@N00" id="33945696@N00" name="PORTRAITS - A Gallery for your Best Portraiture" member="1" moderator="0" admin="0" privacy="3" photos="1195519" iconserver="5485" iconfarm="6" member_count="35924" topic_count="76" pool_count="1195519" />
    <group nsid="2750200@N22" id="2750200@N22" name="Portraits / Retratos" member="1" moderator="0" admin="0" privacy="3" photos="512156" iconserver="7578" iconfarm="8" member_count="33325" topic_count="0" pool_count="512156" />
    <group nsid="1430965@N20" id="1430965@N20" name="Portraits and Faces" member="1" moderator="0" admin="0" privacy="3" photos="2268626" iconserver="65535" iconfarm="66" member_count="100738" topic_count="158" pool_count="2268626" />
    <group nsid="924761@N22" id="924761@N22" name="Portraits du monde" member="1" moderator="0" admin="0" privacy="3" photos="253324" iconserver="5324" iconfarm="6" member_count="11798" topic_count="34" pool_count="253324" />
    <group nsid="1171283@N22" id="1171283@N22" name="Portraits have Souls" member="1" moderator="0" admin="0" privacy="3" photos="399144" iconserver="2491" iconfarm="3" member_count="13504" topic_count="11" pool_count="399144" />
    <group nsid="1642531@N23" id="1642531@N23" name="Portraits of women (woman only)" member="1" moderator="0" admin="0" privacy="3" photos="516767" iconserver="5297" iconfarm="6" member_count="15749" topic_count="13" pool_count="516767" />
    <group nsid="1597446@N21" id="1597446@N21" name="Portraits Only" member="1" moderator="0" admin="0" privacy="1" photos="777110" iconserver="65535" iconfarm="66" member_count="27973" topic_count="25" pool_count="777110" />
    <group nsid="1994388@N22" id="1994388@N22" name="PORTRAITS POUR FEMMES DU MONDE" member="1" moderator="0" admin="0" privacy="3" photos="152279" iconserver="5458" iconfarm="6" member_count="6092" topic_count="5" pool_count="152279" />
    <group nsid="1727343@N24" id="1727343@N24" name="Portraits Super  (Only High Quality Images)" member="1" moderator="0" admin="0" privacy="3" photos="207381" iconserver="6207" iconfarm="7" member_count="6613" topic_count="14" pool_count="207381" />
    <group nsid="363008@N20" id="363008@N20" name="Portraits Unlimited" member="1" moderator="0" admin="0" privacy="3" photos="713838" iconserver="204" iconfarm="1" member_count="12441" topic_count="29" pool_count="713838" />
    <group nsid="1311364@N25" id="1311364@N25" name="PORTRAITS with credits of 'GOOD PHOTO'" member="1" moderator="0" admin="0" privacy="3" photos="331174" iconserver="4008" iconfarm="5" member_count="18499" topic_count="5" pool_count="331174" />
    <group nsid="1202963@N23" id="1202963@N23" name="Portraits with natural light" member="1" moderator="0" admin="0" privacy="3" photos="715707" iconserver="2598" iconfarm="3" member_count="34090" topic_count="36" pool_count="715707" />
    <group nsid="44001204@N00" id="44001204@N00" name="Portraiture" member="1" moderator="0" admin="0" privacy="3" photos="1219815" iconserver="65535" iconfarm="66" member_count="38078" topic_count="86" pool_count="1219815" />
    <group nsid="29496069@N00" id="29496069@N00" name="Portraiture Photography" member="1" moderator="0" admin="0" privacy="3" photos="2099085" iconserver="65535" iconfarm="66" member_count="79160" topic_count="316" pool_count="2099085" />
    <group nsid="58146428@N00" id="58146428@N00" name="Portrait★Faces" member="1" moderator="0" admin="0" privacy="3" photos="149622" iconserver="172" iconfarm="1" member_count="52619" topic_count="58" pool_count="149622" />
    <group nsid="26964488@N00" id="26964488@N00" name="Powerful Portraits" member="1" moderator="0" admin="0" privacy="3" photos="357208" iconserver="3088" iconfarm="4" member_count="10611" topic_count="11" pool_count="357208" />
    <group nsid="90935949@N00" id="90935949@N00" name="PRIME 35mm LENS" member="1" moderator="0" admin="0" privacy="3" photos="175512" iconserver="45" iconfarm="1" member_count="6617" topic_count="12" pool_count="175512" />
    <group nsid="910949@N22" id="910949@N22" name="Quality Black &amp; White Images" member="1" moderator="0" admin="0" privacy="2" photos="248827" iconserver="3198" iconfarm="4" member_count="8897" topic_count="356" pool_count="248827" />
    <group nsid="1801554@N22" id="1801554@N22" name="ReallyGoodNudes" member="1" moderator="0" admin="0" privacy="3" photos="26260" iconserver="2844" iconfarm="3" member_count="4707" topic_count="15" pool_count="26260" />
    <group nsid="14823827@N22" id="14823827@N22" name="Related to Water in some way" member="1" moderator="0" admin="0" privacy="3" photos="3221" iconserver="0" iconfarm="0" member_count="383" topic_count="1" pool_count="3221" />
    <group nsid="3945454@N24" id="3945454@N24" name="Rubens Beauty of all Ages (PLUMP FEMALE)" member="1" moderator="0" admin="0" privacy="1" photos="12581" iconserver="870" iconfarm="1" member_count="1152" topic_count="11" pool_count="12581" />
    <group nsid="1196701@N25" id="1196701@N25" name="Samyang" member="1" moderator="0" admin="0" privacy="3" photos="82792" iconserver="7440" iconfarm="8" member_count="4745" topic_count="50" pool_count="82792" />
    <group nsid="1415923@N22" id="1415923@N22" name="School of Digital Photography" member="1" moderator="0" admin="0" privacy="3" photos="1215593" iconserver="5501" iconfarm="6" member_count="19278" topic_count="132" pool_count="1215593" />
    <group nsid="35034348455@N01" id="35034348455@N01" name="Seattle" member="1" moderator="0" admin="0" privacy="3" photos="200933" iconserver="5" iconfarm="1" member_count="8701" topic_count="425" pool_count="200933" />
    <group nsid="86249273@N00" id="86249273@N00" name="Seattle Scene" member="1" moderator="0" admin="0" privacy="3" photos="32798" iconserver="34" iconfarm="1" member_count="2815" topic_count="93" pool_count="32798" />
    <group nsid="25779358@N00" id="25779358@N00" name="Seattle Strobist" member="1" moderator="0" admin="0" privacy="3" photos="6610" iconserver="125" iconfarm="1" member_count="1025" topic_count="89" pool_count="6610" />
    <group nsid="391302@N23" id="391302@N23" name="Self Expression (Moody Self-Portraits)" member="1" moderator="0" admin="0" privacy="3" photos="82591" iconserver="208" iconfarm="1" member_count="7984" topic_count="14" pool_count="82591" />
    <group nsid="91514935@N00" id="91514935@N00" name="Self Taught Photographers - LOOK AT DISCUSSIONS FIRST!" member="1" moderator="0" admin="0" privacy="3" photos="1822963" iconserver="118" iconfarm="1" member_count="59426" topic_count="1143" pool_count="1822963" />
    <group nsid="1339103@N21" id="1339103@N21" name="Self-Portrait Café (SPC)" member="1" moderator="0" admin="0" privacy="3" photos="62320" iconserver="2594" iconfarm="3" member_count="5939" topic_count="9" pool_count="62320" />
    <group nsid="89438260@N00" id="89438260@N00" name="Self-Portraits!" member="1" moderator="0" admin="0" privacy="3" photos="583608" iconserver="6" iconfarm="1" member_count="45517" topic_count="212" pool_count="583608" />
    <group nsid="2710168@N24" id="2710168@N24" name="Sensual Naked Art (WOMEN NUDE !!!)" member="1" moderator="0" admin="0" privacy="2" photos="41429" iconserver="7817" iconfarm="8" member_count="1342" topic_count="9" pool_count="41429" />
    <group nsid="18329448@N00" id="18329448@N00" name="Sexy" member="1" moderator="0" admin="0" privacy="3" photos="150102" iconserver="65535" iconfarm="66" member_count="15144" topic_count="46" pool_count="150102" />
    <group nsid="2921800@N21" id="2921800@N21" name="Sexy models posing" member="1" moderator="1" admin="0" privacy="3" photos="481007" iconserver="65535" iconfarm="66" member_count="18835" topic_count="6" pool_count="481007" />
    <group nsid="2010877@N23" id="2010877@N23" name="show us what you have got" member="1" moderator="0" admin="0" privacy="1" photos="136281" iconserver="0" iconfarm="0" member_count="5777" topic_count="81" pool_count="136281" />
    <group nsid="2471391@N20" id="2471391@N20" name="Sigma 50mm f/1.4 DG HSM Art" member="1" moderator="0" admin="0" privacy="3" photos="31473" iconserver="7321" iconfarm="8" member_count="4122" topic_count="34" pool_count="31473" />
    <group nsid="2447289@N22" id="2447289@N22" name="Sony A7 / A7R / A7S / A7II / A7RII / A7SII / A7III / A7RIII Comm" member="1" moderator="0" admin="0" privacy="3" photos="375325" iconserver="3718" iconfarm="4" member_count="11450" topic_count="156" pool_count="375325" />
    <group nsid="14809450@N23" id="14809450@N23" name="Sony A7IV / a7 IV / ILCE-7M4 (33MP) - No 30/60 and # Limits" member="1" moderator="0" admin="0" privacy="3" photos="9147" iconserver="65535" iconfarm="66" member_count="515" topic_count="6" pool_count="9147" />
    <group nsid="2335287@N24" id="2335287@N24" name="SONY Alpha A7, A7r, A7S, A7 II &amp; A7r II" member="1" moderator="0" admin="0" privacy="3" photos="346646" iconserver="7437" iconfarm="8" member_count="7650" topic_count="82" pool_count="346646" />
    <group nsid="925860@N22" id="925860@N22" name="Sony Alpha Community" member="1" moderator="0" admin="0" privacy="3" photos="310466" iconserver="3816" iconfarm="4" member_count="4513" topic_count="15" pool_count="310466" />
    <group nsid="1384461@N24" id="1384461@N24" name="Sony Camera Club" member="1" moderator="0" admin="0" privacy="3" photos="757851" iconserver="5523" iconfarm="6" member_count="25665" topic_count="1839" pool_count="757851" />
    <group nsid="3452827@N21" id="3452827@N21" name="Sony FE 85mm f/1.8 (SEL85F18)" member="1" moderator="0" admin="0" privacy="3" photos="13232" iconserver="3942" iconfarm="4" member_count="2948" topic_count="5" pool_count="13232" />
    <group nsid="30494513@N00" id="30494513@N00" name="strobe" member="1" moderator="0" admin="0" privacy="3" photos="40089" iconserver="3" iconfarm="1" member_count="2468" topic_count="36" pool_count="40089" />
    <group nsid="677112@N23" id="677112@N23" name="Strobist Gear" member="1" moderator="0" admin="0" privacy="3" photos="28331" iconserver="2352" iconfarm="3" member_count="4625" topic_count="208" pool_count="28331" />
    <group nsid="71917374@N00" id="71917374@N00" name="Strobist.com" member="1" moderator="0" admin="0" privacy="2" photos="615524" iconserver="49" iconfarm="1" member_count="122271" topic_count="51927" pool_count="615524" />
    <group nsid="52242317293@N01" id="52242317293@N01" name="Sunsets &amp; Sunrises around the world (We Rock Again!)" member="1" moderator="0" admin="0" privacy="3" photos="1403967" iconserver="65535" iconfarm="66" member_count="131237" topic_count="696" pool_count="1403967" />
    <group nsid="4529769@N21" id="4529769@N21" name="Tamron 28-75mm F2.8 DiIII RXD (A036)" member="1" moderator="0" admin="0" privacy="3" photos="11555" iconserver="887" iconfarm="1" member_count="1488" topic_count="2" pool_count="11555" />
    <group nsid="14709378@N20" id="14709378@N20" name="Tamron 70-180 mm F/2.8 Di iii VXD" member="1" moderator="0" admin="0" privacy="3" photos="2140" iconserver="0" iconfarm="0" member_count="370" topic_count="2" pool_count="2140" />
    <group nsid="2875318@N21" id="2875318@N21" name="Team Sony" member="1" moderator="0" admin="0" privacy="3" photos="709789" iconserver="5784" iconfarm="6" member_count="5285" topic_count="86" pool_count="709789" />
    <group nsid="1400663@N22" id="1400663@N22" name="Thai cute &amp; sexy girls" member="1" moderator="0" admin="0" privacy="3" photos="452" iconserver="4061" iconfarm="5" member_count="612" topic_count="2" pool_count="452" />
    <group nsid="14756222@N22" id="14756222@N22" name="THE BEST NUDE ART PHOTOGRAPHY" member="1" moderator="0" admin="0" privacy="3" photos="45194" iconserver="65535" iconfarm="66" member_count="2578" topic_count="5" pool_count="45194" />
    <group nsid="4488127@N23" id="4488127@N23" name="The best photographs from around the world." member="1" moderator="0" admin="0" privacy="3" photos="3902365" iconserver="65535" iconfarm="66" member_count="45841" topic_count="1563" pool_count="3902365" />
    <group nsid="917745@N23" id="917745@N23" name="The Best Portraits AOI ~AOI L2~(Admin Invite) (P1-Award &amp; Fave3)" member="1" moderator="0" admin="0" privacy="3" photos="17148" iconserver="3677" iconfarm="4" member_count="3237" topic_count="8" pool_count="17148" />
    <group nsid="20421382@N00" id="20421382@N00" name="The Portrait Group" member="1" moderator="0" admin="0" privacy="3" photos="3392117" iconserver="65535" iconfarm="66" member_count="166806" topic_count="2701" pool_count="3392117" />
    <group nsid="1760856@N23" id="1760856@N23" name="Tra un Manifesto e Lo Specchio - High quality images" member="1" moderator="0" admin="0" privacy="3" photos="693951" iconserver="2836" iconfarm="3" member_count="11730" topic_count="3" pool_count="693951" />
    <group nsid="14789763@N21" id="14789763@N21" name="TTArtisan Lenses" member="1" moderator="0" admin="0" privacy="3" photos="1786" iconserver="65535" iconfarm="66" member_count="210" topic_count="1" pool_count="1786" />
    <group nsid="70344726@N00" id="70344726@N00" name="Turning the lens on Flickr (Self Portraits)" member="1" moderator="0" admin="0" privacy="3" photos="47313" iconserver="72" iconfarm="1" member_count="3515" topic_count="5" pool_count="47313" />
    <group nsid="1960966@N23" id="1960966@N23" name="VERAMENTE MAGICHE (GLAMOUR NUDE&amp;PORTRAIT)" member="1" moderator="0" admin="0" privacy="2" photos="36603" iconserver="1892" iconfarm="2" member_count="1374" topic_count="9" pool_count="36603" />
    <group nsid="32266655@N00" id="32266655@N00" name="Views: 500" member="1" moderator="0" admin="0" privacy="3" photos="434996" iconserver="4276" iconfarm="5" member_count="20753" topic_count="217" pool_count="434996" />
    <group nsid="74571383@N00" id="74571383@N00" name="Vignetting" member="1" moderator="0" admin="0" privacy="3" photos="370129" iconserver="16" iconfarm="1" member_count="25255" topic_count="72" pool_count="370129" />
    <group nsid="49287201@N00" id="49287201@N00" name="VirusPhoto" member="1" moderator="0" admin="0" privacy="3" photos="488068" iconserver="7569" iconfarm="8" member_count="5285" topic_count="78" pool_count="488068" />
    <group nsid="52240442914@N01" id="52240442914@N01" name="Water" member="1" moderator="0" admin="0" privacy="3" photos="223520" iconserver="1" iconfarm="1" member_count="7608" topic_count="24" pool_count="223520" />
    <group nsid="14743297@N22" id="14743297@N22" name="we need poetry and beauty" member="1" moderator="0" admin="0" privacy="3" photos="26047" iconserver="65535" iconfarm="66" member_count="1104" topic_count="6" pool_count="26047" />
    <group nsid="1903354@N20" id="1903354@N20" name="We Shoot RAW" member="1" moderator="0" admin="0" privacy="3" photos="516944" iconserver="8424" iconfarm="9" member_count="5683" topic_count="18" pool_count="516944" />
    <group nsid="727975@N20" id="727975@N20" name="woman portrait" member="1" moderator="0" admin="0" privacy="3" photos="113730" iconserver="2062" iconfarm="3" member_count="4456" topic_count="3" pool_count="113730" />
    <group nsid="1533426@N23" id="1533426@N23" name="Woman. Period." member="1" moderator="0" admin="0" privacy="2" photos="146292" iconserver="3725" iconfarm="4" member_count="3261" topic_count="16" pool_count="146292" />
    <group nsid="1638764@N21" id="1638764@N21" name="Women *Simply beautiful*" member="1" moderator="0" admin="0" privacy="3" photos="255915" iconserver="65535" iconfarm="66" member_count="5333" topic_count="14" pool_count="255915" />
    <group nsid="421891@N25" id="421891@N25" name="Women - Portraits of seduction and passion" member="1" moderator="0" admin="0" privacy="3" photos="53483" iconserver="1300" iconfarm="2" member_count="3804" topic_count="6" pool_count="53483" />
    <group nsid="14679529@N23" id="14679529@N23" name="Women who love to expose in public outside" member="1" moderator="0" admin="0" privacy="1" photos="13039" iconserver="65535" iconfarm="66" member_count="2920" topic_count="2" pool_count="13039" />
    <group nsid="768751@N23" id="768751@N23" name="Women Women Women" member="1" moderator="0" admin="0" privacy="3" photos="519077" iconserver="2137" iconfarm="3" member_count="10989" topic_count="25" pool_count="519077" />
    <group nsid="88292168@N00" id="88292168@N00" name="world wide women" member="1" moderator="0" admin="0" privacy="3" photos="212809" iconserver="7" iconfarm="1" member_count="12275" topic_count="48" pool_count="212809" />
    <group nsid="857078@N20" id="857078@N20" name="worldlightning" member="1" moderator="0" admin="0" privacy="2" photos="388809" iconserver="3162" iconfarm="4" member_count="5904" topic_count="35" pool_count="388809" />
    <group nsid="1517082@N22" id="1517082@N22" name="WOW FACTOR.WOMEN. (BEST NUDES &amp; EROTICA)" member="1" moderator="0" admin="0" privacy="1" photos="28975" iconserver="65535" iconfarm="66" member_count="1391" topic_count="10" pool_count="28975" />
    <group nsid="1278445@N21" id="1278445@N21" name="You are woman, you are beautiful, let the world see you" member="1" moderator="0" admin="0" privacy="3" photos="268455" iconserver="3809" iconfarm="4" member_count="8760" topic_count="20" pool_count="268455" />
    <group nsid="25078481@N00" id="25078481@N00" name="you can Canon" member="1" moderator="0" admin="0" privacy="3" photos="1120507" iconserver="30" iconfarm="1" member_count="7396" topic_count="32" pool_count="1120507" />
    <group nsid="43116471@N00" id="43116471@N00" name="Young Photographers" member="1" moderator="0" admin="0" privacy="3" photos="3566011" iconserver="7007" iconfarm="8" member_count="115444" topic_count="1341" pool_count="3566011" />
    <group nsid="1609077@N23" id="1609077@N23" name="zwartwit" member="1" moderator="0" admin="0" privacy="1" photos="588152" iconserver="5250" iconfarm="6" member_count="12427" topic_count="150" pool_count="588152" />
    <group nsid="14812466@N22" id="14812466@N22" name="zwartwit." member="1" moderator="0" admin="0" privacy="1" photos="55358" iconserver="65535" iconfarm="66" member_count="3239" topic_count="16" pool_count="55358" />
    <group nsid="14724474@N24" id="14724474@N24" name="★ PORTRAIT BEAUTY OVER 1000 VIEWS" member="1" moderator="0" admin="0" privacy="3" photos="4996" iconserver="65535" iconfarm="66" member_count="259" topic_count="0" pool_count="4996" />
    <group nsid="1099987@N20" id="1099987@N20" name="★★★Shutter Bug - Tips &amp; Tricks★★★" member="1" moderator="0" admin="0" privacy="3" photos="2140017" iconserver="3361" iconfarm="4" member_count="28593" topic_count="397" pool_count="2140017" />
    <group nsid="888631@N25" id="888631@N25" name="♥♥ People around us ♥♥" member="1" moderator="0" admin="0" privacy="3" photos="960446" iconserver="3062" iconfarm="4" member_count="11546" topic_count="20" pool_count="960446" />
  </groups>
</rsp>`

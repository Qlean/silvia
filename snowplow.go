package main

import(
	// "fmt"
	"time"
	"strings"
	"strconv"
	// "encoding/json"
)

// type SnowplowContentsData struct {
// 	UtmSource   string `json:"utm_source"`
// 	UtmMedium   string `json:"utm_medium"`
// 	UtmCampaign string `json:"utm_campaign"`
// 	UtmContent  string `json:"utm_content"`
// 	UtmTerm     string `json:"utm_term"`
// 	GTestID     string `json:"google_test_id"`
// 	GTestVar    string `json:"google_test_var"`
// }

// type SnowplowContents struct {
// 	Schema string                 `json:"schema"`
// 	Data   []SnowplowContentsData `json:"data"`
// }

type SnowplowTimestamp struct {
	time.Time
}

type ScreenResolution struct {
	Width int
	Height int
}

type SnowplowData struct {
	Event           string            `json:"e"`
	PageURL         string            `json:"url"`
	PageTtile       string            `json:"page"`
	PageReferrer    string            `json:"refr"`
	VTracker        string            `json:"tv"`
	Aid             string            `json:"aid"`
	Platform        string            `json:"p"`
	OsTimezone      string            `json:"tz"`
	BrLang          string            `json:"lang"`
	DocCharset      string            `json:"cs"`
	BrPDF           string            `json:"f_pdf"`
	BrQT            string            `json:"f_qt"`
	BrRealPlayer    string            `json:"f_realp"`
	BrWMA           string            `json:"f_wma"`
	BrDir           string            `json:"f_dir"`
	BrFlash         string            `json:"f_fla"`
	BrJava          string            `json:"f_java"`
	BrGears         string            `json:"f_gears"`
	BrAg            string            `json:"f_ag"`
	DvceScreen      ScreenResolution  `json:"res"`
	BrColorDepth    string            `json:"cd"`
	BrCookies       string            `json:"cookie"`
	EventID         string            `json:"eid"`
	DvceTstamp      SnowplowTimestamp `json:"dtm"`
	Contexts        string            `json:"co"`
	BrView          ScreenResolution  `json:"vp"`
	Doc             ScreenResolution  `json:"ds"`
	VisitNum        string            `json:"vid"`
	DomainUserID    string            `json:"duid"`
	UserFingerprint string            `json:"fp"`
	UserID          string            `json:"uid"`
}


type Snowplow struct {
	Schema string         `json:"schema"`
	Data   []SnowplowData `json:"data"`
}


func (timestamp *SnowplowTimestamp) UnmarshalJSON(b []byte) (err error) {
	if b[0] == '"' && b[len(b)-1] == '"' {
		b = b[1 : len(b)-1]
	}

	strUnix := string(b)
	if len(strUnix) >= 13 { strUnix = strUnix[0 : len(strUnix)-3] }

	i, err := strconv.ParseInt(strUnix, 10, 64)
	timestamp.Time = time.Unix(i, 0)
	return
}

func (res *ScreenResolution) UnmarshalJSON(b []byte) (err error) {
	if b[0] == '"' && b[len(b)-1] == '"' {
		b = b[1 : len(b)-1]
	}
	s := strings.Split(string(b), "x")
	res.Width, err = strconv.Atoi(s[0])
	res.Height, err = strconv.Atoi(s[1])
	return
}


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
	E               string            `json:"e" db:"event"`
	Url             string            `json:"url" db:"page_url"`
	Page            string            `json:"page" db:"page_title"`
	Referer         string            `json:"refr" db:"page_referrer"`
	Tv              string            `json:"tv" db:"v_tracker"`
	AID             string            `json:"aid" db:"app_id"`
	P               string            `json:"p" db:"platform"`
	Timezone        string            `json:"tz" db:"os_timezone"`
	Lang            string            `json:"lang" db:"br_lang"`
	Cs              string            `json:"cs" db:"doc_charset"`
	PDF             string            `json:"f_pdf"`
	QT              string            `json:"f_qt"`
	RealPlayer      string            `json:"f_realp"`
	WMA             string            `json:"f_wma"`
	Dir             string            `json:"f_dir"`
	Flash           string            `json:"f_fla"`
	Java            string            `json:"f_java"`
	Gears           string            `json:"f_gears"`
	Ag              string            `json:"f_ag"`
	Res             ScreenResolution  `json:"res"`
	Cd              string            `json:"cd" db:"br_colordepth"`
	Cookie          string            `json:"cookie"`
	EID             string            `json:"eid" db:"event_id"`
	DeviceTimestamp SnowplowTimestamp `json:"dtm"`
	Contexts        string            `json:"co" db:"contexts"`
	Vp              ScreenResolution  `json:"vp"`
	Ds              ScreenResolution  `json:"ds"`
	Vid             string            `json:"vid" db:"visit_num"`
	DUID            string            `json:"duid" db:"domain_userid"`
	Fp              string            `json:"fp" db:"user_fingerprint"`
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


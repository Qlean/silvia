package main

import(
	"time"
)

type SnowplowContentsData struct {
	UtmSource   string `json:"utm_source"`
	UtmMedium   string `json:"utm_medium"`
	UtmCampaign string `json:"utm_campaign"`
	UtmContent  string `json:"utm_content"`
	UtmTerm     string `json:"utm_term"`
	GTestID     string `json:"google_test_id"`
	GTestVar    string `json:"google_test_var"`
}

type SnowplowContents struct {
	Schema string                       `json:"schema"`
	Data   map[int]SnowplowContentsData `json:"data"`
}

type SnowplowDeviceTimestamp struct {
	time.Time
}

type ScreenResolution struct {
	Width int
	Height int
}

type SnowplowData struct {
	E               string                  `json:"e"`
	Url             string                  `json:"url"`
	Page            string                  `json:"page"`
	Referer         string                  `json:"refr"`
	Tv              string                  `json:"tv"`
	Tna             string                  `json:"tna"`
	AID             string                  `json:"aid"`
	P               string                  `json:"p"`
	Timezone        string                  `json:"tz"`
	Lang            string                  `json:"lang"`
	Cs              string                  `json:"cs"`
	PDF             bool                    `json:"f_pdf"`
	QT              bool                    `json:"f_qt"`
	RealPlayer      bool                    `json:"f_realp"`
	WMA             bool                    `json:"f_wma"`
	Dir             bool                    `json:"f_dir"`
	Flash           bool                    `json:"f_fla"`
	Java            bool                    `json:"f_java"`
	Gears           bool                    `json:"f_gears"`
	Ag              bool                    `json:"f_ag"`
	Res             ScreenResolution        `json:"res"`
	Cd              int                     `json:"cd"`
	Cookie          bool                    `json:"cookie"`
	EID             string                  `json:"eid"`
	DeviceTimestamp SnowplowDeviceTimestamp `json:"dtm"`
	Contents        SnowplowContents        `json:"co"`
	Vp              ScreenResolution        `json:"vp"`
	Ds              ScreenResolution        `json:"ds"`
	Vid             int                     `json:"vid"`
	DUID            string                  `json:"duid"`
	Fp              string                  `json:"fp"`
}

type Snowplow struct {
	Schema string       `json:"schema"`
	Data   map[int]SnowplowData `json:"data"`
}

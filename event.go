package main

import(
	"log"
	"time"
	"strconv"
	"encoding/json"

	"github.com/mssola/user_agent"
)

type SnowplowContextsData struct {
	Data struct {
		UtmSource   string `json:"utm_source"`
		UtmMedium   string `json:"utm_medium"`
		UtmCampaign string `json:"utm_campaign"`
		UtmContent  string `json:"utm_content"`
		UtmTerm     string `json:"utm_term"`
	}
}

type SnowplowContexts struct {
	Schema string                 `json:"schema"`
	Data   []SnowplowContextsData `json:"data"`
}


type Event struct {
	Id              int    `db:"-"`
	Aid             string `db:"app_id"`
	Platform        string `db:"platform"`
	CollectorTstamp time.Time `db:"collector_tstamp"`
	DvceTstamp      time.Time `db:"dvce_tstamp"`
	Event           string `db:"event"`
	EventID         string `db:"event_id"`
	VTracker        string `db:"v_tracker"`
	UserID          int `db:"user_id"`
	UserIP          string `db:"user_ipaddress"`
	UserFingerprint string `db:"user_fingerprint"`
	DomainUserID    string `db:"domain_userid"`
	VisitNum        int `db:"visit_num"`
	SessionID       string `db:"session_id"`
	GeoCountry      string `db:"geo_country"`
	GeoRegion       string `db:"geo_region"`
	GeoCity         string `db:"geo_city"`
	GeoZipcode      string `db:"geo_zipcode"`
	GeoLatitude     float32 `db:"geo_latitude"`
	GeoLongtitude   float32 `db:"geo_longitude"`
	GeoRegionName   string `db:"geo_region_name"`
	PageURL         string `db:"page_url"`
	PageTtile       string `db:"page_title"`
	PageReferrer    string `db:"page_referrer"`
	PageURLHost     string `db:"page_urlhost"`
	PageURLPath     string `db:"page_urlpath"`
	PageURLQuery    string `db:"page_urlquery"`
	RefrURLHost     string `db:"refr_urlhost"`
	RefrURLPath     string `db:"refr_urlpath"`
	RefrURLQuety    string `db:"refr_urlquery"`
	UtmMedium       string `db:"utm_medium"`
	UtmSource       string `db:"utm_source"`
	UtmTerm         string `db:"utm_term"`
	UtmContent      string `db:"utm_content"`
	UtmCampaign     string `db:"utm_campaign"`
	Contexts        string `db:"contexts"`
	SeCategory      string `db:"se_category"`
	SeAction        string `db:"se_action"`
	SeLabel         string `db:"se_label"`
	SeProperty      string `db:"se_property"`
	SeValue         float32 `db:"se_value"`
	UnstructEvent   string `db:"unstruct_event"`
	Useragent       string `db:"useragent"`
	BrName          string `db:"br_name"`
	BrFamily        string `db:"br_family"`
	BrVersion       string `db:"br_version"`
	BrType          string `db:"br_type"`
	BrRndrNgn       string `db:"br_renderengine"`
	BrLang          string `db:"br_lang"`
	BrPDF           bool `db:"br_features_pdf"`
	BrFlash         bool `db:"br_features_flash"`
	BrJava          bool `db:"br_features_java"`
	BrDir           bool `db:"br_features_director"`
	BrQT            bool `db:"br_features_quicktime"`
	BrRealPlayer    bool `db:"br_features_realplayer"`
	BrWMA           bool `db:"br_features_windowsmedia"`
	BrGears         bool `db:"br_features_gears"`
	BrAg            bool `db:"br_features_silverlight"`
	BrCookies       bool `db:"br_cookies"`
	BrColorDepth    int  `db:"br_colordepth"`
	BrViewWidth     int `db:"br_viewwidth"`
	BrViewHeight    int `db:"br_viewheight"`
	OsName          string `db:"os_name"`
	OsFamily        string `db:"os_family"`
	OsManufacturer  string `db:"os_manufacturer"`
	OsTimezone      string `db:"os_timezone"`
	DvceType        string `db:"dvce_type"`
	DvceIsMobile    bool `db:"dvce_ismobile"`
	DvceScreenWidth int `db:"dvce_screenwidth"`
	DvceScreenHeight int `db:"dvce_screenheight"`
	DocCharset      string `db:"doc_charset"`
	DocWidth        int `db:"doc_width"`
	DocHeight       int `db:"doc_height"`
}

func nginxLogToEvent(nginxLog *NginxLog) (event *Event) {
	// Initialize object
	event = &Event{}

	// Binding first level nginx fields
	event.UserIP = nginxLog.IPAddress
	event.CollectorTstamp = time.Now()
	event.PageReferrer = nginxLog.Referer
	event.Useragent = nginxLog.Useragent

	// Binding snowplow data
	snowplowData := nginxLog.RequestBody.Data[0]

	// Bind string values
	event.Event = snowplowData.Event
	event.PageURL = snowplowData.PageURL
	event.PageTtile = snowplowData.PageTtile
	event.PageReferrer = snowplowData.PageReferrer
	event.VTracker = snowplowData.VTracker
	event.Aid = snowplowData.Aid
	event.Platform = snowplowData.Platform
	event.OsTimezone = snowplowData.OsTimezone
	event.BrLang = snowplowData.BrLang
	event.DocCharset = snowplowData.DocCharset
	event.EventID = snowplowData.EventID
	event.DomainUserID = snowplowData.DomainUserID
	event.UserFingerprint = snowplowData.UserFingerprint

	// Bind integer values
	event.BrColorDepth, _ = strconv.Atoi(snowplowData.BrColorDepth)
	event.VisitNum, _ = strconv.Atoi(snowplowData.VisitNum)
	event.UserID, _ = strconv.Atoi(snowplowData.UserID)

	// Bind boolean values
	event.BrPDF, _ = strconv.ParseBool(snowplowData.BrPDF)
	event.BrQT, _ = strconv.ParseBool(snowplowData.BrQT)
	event.BrRealPlayer, _ = strconv.ParseBool(snowplowData.BrRealPlayer)
	event.BrWMA, _ = strconv.ParseBool(snowplowData.BrWMA)
	event.BrDir, _ = strconv.ParseBool(snowplowData.BrDir)
	event.BrFlash, _ = strconv.ParseBool(snowplowData.BrFlash)
	event.BrJava, _ = strconv.ParseBool(snowplowData.BrJava)
	event.BrGears, _ = strconv.ParseBool(snowplowData.BrGears)
	event.BrAg, _ = strconv.ParseBool(snowplowData.BrAg)
	event.BrCookies, _ = strconv.ParseBool(snowplowData.BrCookies)

	// Bind screen resolutions
	event.BrViewWidth = snowplowData.BrView.Width
	event.BrViewHeight = snowplowData.BrView.Height

	event.DvceScreenWidth = snowplowData.DvceScreen.Width
	event.DvceScreenHeight = snowplowData.DvceScreen.Height

	event.DocWidth = snowplowData.Doc.Width
	event.DocHeight = snowplowData.Doc.Width

	// Bind time values
	event.DvceTstamp = snowplowData.DvceTstamp.Time

	// Bind Geo values
	geoIPRecord := geo.GetRecord(nginxLog.IPAddress)
	event.GeoCountry = geoIPRecord.CountryName
	event.GeoRegion = geoIPRecord.Region
	event.GeoCity = geoIPRecord.City
	event.GeoZipcode = geoIPRecord.PostalCode
	event.GeoLatitude = geoIPRecord.Latitude
	event.GeoLongtitude = geoIPRecord.Longitude
	event.GeoRegionName = geoIPRecord.Region

	// Bind useragent values

	ua := user_agent.New(nginxLog.Useragent)
	event.DvceIsMobile = ua.Mobile()

	event.OsName = ua.OS()
	event.DvceType = ua.Platform()
	name, version := ua.Engine()
	event.BrRndrNgn = name + version

	name, version = ua.Browser()
	event.BrName = name
	event.BrVersion = version

	// Bind contexts
	if len(snowplowData.Contexts) == 0 { snowplowData.Contexts = "{}" }
	event.Contexts = snowplowData.Contexts

	contexts := &SnowplowContexts{}
	err := json.Unmarshal([]byte(snowplowData.Contexts), &contexts)
	if err != nil {
		log.Println("Cant Unmarshal contexts: ", err)
	}

	contextsData := contexts.Data[0].Data
	event.UtmCampaign = contextsData.UtmCampaign
	event.UtmContent = contextsData.UtmContent
	event.UtmTerm = contextsData.UtmTerm
	event.UtmSource = contextsData.UtmSource
	event.UtmMedium = contextsData.UtmMedium

	// Bind unstructured event
	if len(snowplowData.UnstructEvent) == 0 { snowplowData.UnstructEvent = "{}" }
	event.UnstructEvent = snowplowData.UnstructEvent

	return
}

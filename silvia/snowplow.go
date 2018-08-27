package silvia

import (
	"database/sql"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/abh/geoip"
	"github.com/mssola/user_agent"
)

type (
	// database structure
	SnowplowContextsData struct {
		Data struct {
			UtmSource   string `json:"utm_source"`
			UtmMedium   string `json:"utm_medium"`
			UtmCampaign string `json:"utm_campaign"`
			UtmContent  string `json:"utm_content"`
			UtmTerm     string `json:"utm_term"`
		}
	}

	SnowplowContexts struct {
		Schema string                 `json:"schema"`
		Data   []SnowplowContextsData `json:"data"`
	}

	SnowplowContextsJsonField struct {
		Schema string            `json:"schema"`
		Data   []json.RawMessage `json:"data"`
	}

	SnowplowEvent struct {
		Body             []byte          `db:"-"`
		Id               int             `db:"-"`
		Aid              sql.NullString  `db:"app_id"`
		Platform         sql.NullString  `db:"platform"`
		CollectorTstamp  time.Time       `db:"collector_tstamp"`
		DvceTstamp       time.Time       `db:"dvce_tstamp"`
		Event            sql.NullString  `db:"event"`
		EventID          sql.NullString  `db:"event_id"`
		VTracker         sql.NullString  `db:"v_tracker"`
		UserID           sql.NullInt64   `db:"user_id"`
		UserIP           sql.NullString  `db:"user_ipaddress"`
		UserFingerprint  sql.NullString  `db:"user_fingerprint"`
		DomainUserID     sql.NullString  `db:"domain_userid"`
		VisitNum         sql.NullInt64   `db:"visit_num"`
		SessionID        sql.NullString  `db:"session_id"`
		GeoCountry       sql.NullString  `db:"geo_country"`
		GeoRegion        sql.NullString  `db:"geo_region"`
		GeoCity          sql.NullString  `db:"geo_city"`
		GeoZipcode       sql.NullString  `db:"geo_zipcode"`
		GeoLatitude      float32         `db:"geo_latitude"`
		GeoLongtitude    float32         `db:"geo_longitude"`
		GeoRegionName    sql.NullString  `db:"geo_region_name"`
		PageURL          sql.NullString  `db:"page_url"`
		PageTtile        sql.NullString  `db:"page_title"`
		PageReferrer     sql.NullString  `db:"page_referrer"`
		PageURLHost      sql.NullString  `db:"page_urlhost"`
		PageURLPath      sql.NullString  `db:"page_urlpath"`
		PageURLQuery     sql.NullString  `db:"page_urlquery"`
		RefrURLHost      sql.NullString  `db:"refr_urlhost"`
		RefrURLPath      sql.NullString  `db:"refr_urlpath"`
		RefrURLQuety     sql.NullString  `db:"refr_urlquery"`
		UtmMedium        sql.NullString  `db:"utm_medium"`
		UtmSource        sql.NullString  `db:"utm_source"`
		UtmTerm          sql.NullString  `db:"utm_term"`
		UtmContent       sql.NullString  `db:"utm_content"`
		UtmCampaign      sql.NullString  `db:"utm_campaign"`
		Contexts         sql.NullString  `db:"contexts"`
		SeCategory       sql.NullString  `db:"se_category"`
		SeAction         sql.NullString  `db:"se_action"`
		SeLabel          sql.NullString  `db:"se_label"`
		SeProperty       sql.NullString  `db:"se_property"`
		SeValue          sql.NullFloat64 `db:"se_value"`
		UnstructEvent    sql.NullString  `db:"unstruct_event"`
		Useragent        sql.NullString  `db:"useragent"`
		BrName           sql.NullString  `db:"br_name"`
		BrFamily         sql.NullString  `db:"br_family"`
		BrVersion        sql.NullString  `db:"br_version"`
		BrType           sql.NullString  `db:"br_type"`
		BrRndrNgn        sql.NullString  `db:"br_renderengine"`
		BrLang           sql.NullString  `db:"br_lang"`
		BrPDF            bool            `db:"br_features_pdf"`
		BrFlash          bool            `db:"br_features_flash"`
		BrJava           bool            `db:"br_features_java"`
		BrDir            bool            `db:"br_features_director"`
		BrQT             bool            `db:"br_features_quicktime"`
		BrRealPlayer     bool            `db:"br_features_realplayer"`
		BrWMA            bool            `db:"br_features_windowsmedia"`
		BrGears          bool            `db:"br_features_gears"`
		BrAg             bool            `db:"br_features_silverlight"`
		BrCookies        bool            `db:"br_cookies"`
		BrColorDepth     sql.NullInt64   `db:"br_colordepth"`
		BrViewWidth      int             `db:"br_viewwidth"`
		BrViewHeight     int             `db:"br_viewheight"`
		OsName           sql.NullString  `db:"os_name"`
		OsFamily         sql.NullString  `db:"os_family"`
		OsManufacturer   sql.NullString  `db:"os_manufacturer"`
		OsTimezone       sql.NullString  `db:"os_timezone"`
		DvceType         sql.NullString  `db:"dvce_type"`
		DvceIsMobile     bool            `db:"dvce_ismobile"`
		DvceScreenWidth  int             `db:"dvce_screenwidth"`
		DvceScreenHeight int             `db:"dvce_screenheight"`
		DocCharset       sql.NullString  `db:"doc_charset"`
		DocWidth         int             `db:"doc_width"`
		DocHeight        int             `db:"doc_height"`
	}

	// request structure
	SnowplowTimestamp struct {
		time.Time
	}

	ScreenResolution struct {
		Width  int
		Height int
	}

	SnowplowData struct {
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
		UnstructEvent   string            `json:"ue_pr"`
		SeCategory      string            `json:"se_ca"`
		SeAction        string            `json:"se_ac"`
		SeLabel         string            `json:"se_la"`
		SeProperty      string            `json:"se_pr"`
		SeValue         string            `json:"se_va"`
	}

	Snowplow struct {
		Schema string         `json:"schema"`
		Data   []SnowplowData `json:"data"`
	}

	NginxTime struct {
		time.Time
	}

	SnowplowRequest struct {
		IPAddress   string    `json:"ip_addr"`
		TimeLocal   NginxTime `json:"time_local"`
		RequestBody Snowplow  `json:"request_body"`
		Referer     string    `json:"http_referer"`
		Useragent   string    `json:"http_user_agent"`
	}
)

func (nginxTime *NginxTime) UnmarshalJSON(b []byte) (err error) {
	if b[0] == '"' && b[len(b)-1] == '"' {
		b = b[1 : len(b)-1]
	}
	layout := "02/Jan/2006:15:04:05 -0700"
	nginxTime.Time, err = time.Parse(layout, string(b))
	return
}

func (timestamp *SnowplowTimestamp) UnmarshalJSON(b []byte) (err error) {
	if b[0] == '"' && b[len(b)-1] == '"' {
		b = b[1 : len(b)-1]
	}

	strUnix := string(b)
	if len(strUnix) >= 13 {
		strUnix = strUnix[0 : len(strUnix)-3]
	}

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

func (event *SnowplowEvent) Transform(request []byte, geo *geoip.GeoIP) error {
	event.Body = request

	normString, err := strconv.Unquote(`"` + string(request) + `"`)
	if err != nil {
		return err
	}

	normString = strings.Replace(normString, "\" ", "", -1)

	snowplowRequest := &SnowplowRequest{}
	err = json.Unmarshal([]byte(normString), snowplowRequest)
	if err != nil {
		return err
	}

	// Binding first level nginx fields
	checkStringForNull(snowplowRequest.IPAddress, &event.UserIP)
	event.CollectorTstamp = time.Now()
	// checkStringForNull(snowplowRequest.Referer, &event.PageReferrer)
	checkStringForNull(snowplowRequest.Useragent, &event.Useragent)

	// Binding snowplow data
	snowplowData := snowplowRequest.RequestBody.Data[0]

	// Bind string values
	checkStringForNull(snowplowData.Event, &event.Event)
	checkStringForNull(snowplowData.PageURL, &event.PageURL)
	checkStringForNull(snowplowData.PageTtile, &event.PageTtile)
	checkStringForNull(snowplowData.PageReferrer, &event.PageReferrer)
	checkStringForNull(snowplowData.VTracker, &event.VTracker)
	checkStringForNull(snowplowData.Aid, &event.Aid)
	checkStringForNull(snowplowData.Platform, &event.Platform)
	checkStringForNull(snowplowData.OsTimezone, &event.OsTimezone)
	checkStringForNull(snowplowData.BrLang, &event.BrLang)
	checkStringForNull(snowplowData.DocCharset, &event.DocCharset)
	checkStringForNull(snowplowData.EventID, &event.EventID)
	checkStringForNull(snowplowData.DomainUserID, &event.DomainUserID)
	checkStringForNull(snowplowData.UserFingerprint, &event.UserFingerprint)
	checkStringForNull(snowplowData.SeCategory, &event.SeCategory)
	checkStringForNull(snowplowData.SeAction, &event.SeAction)
	checkStringForNull(snowplowData.SeLabel, &event.SeLabel)
	checkStringForNull(snowplowData.SeProperty, &event.SeProperty)
	checkFloatForNull(snowplowData.SeValue, &event.SeValue)

	// Bind integer values
	checkIntForNull(snowplowData.BrColorDepth, &event.BrColorDepth)
	checkIntForNull(snowplowData.VisitNum, &event.VisitNum)
	checkIntForNull(snowplowData.UserID, &event.UserID)

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
	geoIPRecord := geo.GetRecord(snowplowRequest.IPAddress)
	if geoIPRecord != nil {
		checkStringForNull(geoIPRecord.CountryName, &event.GeoCountry)
		checkStringForNull(geoIPRecord.Region, &event.GeoRegion)
		checkStringForNull(geoIPRecord.City, &event.GeoCity)
		checkStringForNull(geoIPRecord.PostalCode, &event.GeoZipcode)
		event.GeoLatitude = geoIPRecord.Latitude
		event.GeoLongtitude = geoIPRecord.Longitude
		checkStringForNull(geoIPRecord.Region, &event.GeoRegionName)
	}

	// Bind useragent values

	ua := user_agent.New(snowplowRequest.Useragent)
	event.DvceIsMobile = ua.Mobile()

	checkStringForNull(ua.OS(), &event.OsName)
	checkStringForNull(ua.Platform(), &event.DvceType)

	name, version := ua.Engine()
	checkStringForNull(name+version, &event.BrRndrNgn)

	name, version = ua.Browser()
	checkStringForNull(name, &event.BrName)
	checkStringForNull(version, &event.BrVersion)

	// Bind contexts
	if len(snowplowData.Contexts) > 0 {
		contextsField := &SnowplowContextsJsonField{}
		err := json.Unmarshal([]byte(snowplowData.Contexts), &contextsField)
		if err != nil {
			return err
		} else {
			rawContexts := contextsField.Data[0]
			rawContexts = rawContexts[8 : len(rawContexts)-1]

			// is json?
			var js interface{}
			err := json.Unmarshal([]byte(rawContexts), &js)
			if err != nil {
				rawContexts = contextsField.Data[0]
			}

			checkStringForNull(string(rawContexts), &event.Contexts)
		}

		contexts := &SnowplowContexts{}
		err = json.Unmarshal([]byte(snowplowData.Contexts), &contexts)
		if err != nil {
			return err
		} else {
			contextsData := contexts.Data[0].Data
			checkStringForNull(contextsData.UtmCampaign, &event.UtmCampaign)
			checkStringForNull(contextsData.UtmContent, &event.UtmContent)
			checkStringForNull(contextsData.UtmTerm, &event.UtmTerm)
			checkStringForNull(contextsData.UtmSource, &event.UtmSource)
			checkStringForNull(contextsData.UtmMedium, &event.UtmMedium)
		}
	} else {
		event.Contexts.Valid = false
	}

	// Bind unstructured event
	if len(snowplowData.UnstructEvent) > 0 {
		checkStringForNull(snowplowData.UnstructEvent, &event.UnstructEvent)
	}

	return nil
}

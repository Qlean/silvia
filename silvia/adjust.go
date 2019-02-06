package silvia

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/serenize/snaker"
)

type (
	NullTime struct {
		Time  time.Time
		Valid bool
	}

	AdjustEvent struct {
		Body       []byte         `db:"-"`
		Id         int            `db:"-"`
		Aid        sql.NullString `db:"app_id"`
		An         sql.NullString `db:"app_name"`
		Av         sql.NullString `db:"app_version"`
		St         sql.NullString `db:"store"`
		Tr         sql.NullString `db:"tracker"`
		Ptr        sql.NullString `db:"previous_tracker"`
		TrNetwork  sql.NullString `db:"tracker_network"`
		TrCampaign sql.NullString `db:"tracker_campaign"`
		TrAdgroup  sql.NullString `db:"tracker_adgroup"`
		TrCreative sql.NullString `db:"tracker_creative"`
		Imp        sql.NullInt64  `db:"impression_based"`
		Org        sql.NullInt64  `db:"is_organic"`
		Rej        sql.NullString `db:"rejection_reason"`
		Adid       sql.NullString `db:"adjust_device_id"`
		Idfa       sql.NullString `db:"idfa"`
		Anid       sql.NullString `db:"android_id"`
		Mac        sql.NullString `db:"mac_sha1"`
		Mt         sql.NullString `db:"match_type"`
		Ref        sql.NullString `db:"referrer"`
		Ua         sql.NullString `db:"user_agent"`
		Ip         sql.NullString `db:"ip_address"`
		Clt        NullTime       `db:"clicked_at"`
		It         NullTime       `db:"installed_at"`
		Ct         NullTime       `db:"event_at"`
		Rt         NullTime       `db:"reattributed_at"`
		Reg        sql.NullString `db:"region"`
		C          sql.NullString `db:"country"`
		Lng        sql.NullString `db:"language"`
		Dn         sql.NullString `db:"device_name"`
		Dt         sql.NullString `db:"device_type"`
		Os         sql.NullString `db:"os_name"`
		Osv        sql.NullString `db:"os_version"`
		Env        sql.NullString `db:"environment"`
		Tre        sql.NullInt64  `db:"tracking_enabled"`
		Tz         sql.NullString `db:"device_timezone"`
		Ls         sql.NullInt64  `db:"last_session_sec"`
		Cs         sql.NullInt64  `db:"current_session_sec"`
		Dl         sql.NullString `db:"deeplink"`
		Pp         sql.NullString `db:"partner_parameters"`
		FbCgn      sql.NullString `db:"fb_campaign_group_name"`
		FbCgid     sql.NullString `db:"fb_campaign_group_id"`
		FbCn       sql.NullString `db:"fb_campaign_name"`
		GbCid      sql.NullString `db:"fb_campaign_id"`
		FbAdn      sql.NullString `db:"fb_adgroup_name"`
		FbAdid     sql.NullString `db:"fb_adgroup_id"`
		Pubp       sql.NullString `db:"publisher_parameters"`
		L          sql.NullString `db:"label"`
		Eg         sql.NullString `db:"event_group"`
		Et         sql.NullString `db:"event_token"`
		En         sql.NullString `db:"event_name"`
		ErrType    sql.NullString `db:error_type`
		Error      sql.NullString `db:error_text`
		ErrorEvent sql.NullString `db:error_event`
	}

	AdjustRequest struct {
		IPAddress string `json:"ip_addr"`
		TimeLocal string `json:"time_local"`
		Request   string `json:"request"`
		Referer   string `json:"http_referer"`
		Useragent string `json:"http_user_agent"`
	}
)

func (nt *NullTime) Scan(value interface{}) error {
	nt.Time, nt.Valid = value.(time.Time)
	return nil
}

func (nt NullTime) Value() (driver.Value, error) {
	if !nt.Valid {
		return nil, nil
	}
	return nt.Time, nil
}

func (event *AdjustEvent) Transform(request []byte) error {
	event.Body = request

	normString, err := strconv.Unquote(`"` + string(request) + `"`)
	if err != nil {
		return err
	}

	normString = strings.Replace(normString, "\" ", "", -1)
	adjustRequest := &AdjustRequest{}
	err = json.Unmarshal([]byte(normString), adjustRequest)
	if err == nil {
		request = []byte(adjustRequest.Request)[4 : len(adjustRequest.Request)-9]
	}

	u, err := url.Parse(string(request))
	if err != nil {
		return err
	}

	queryParams := u.Query()

	structType := reflect.TypeOf(*event)
	structValue := reflect.ValueOf(event).Elem()

	nullStringType := reflect.TypeOf(sql.NullString{})
	nullInt64Type := reflect.TypeOf(sql.NullInt64{})
	NullTimeType := reflect.TypeOf(NullTime{})

	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		fieldName := field.Name
		requestParamName := snaker.CamelToSnake(fieldName)
		requestParamValue := queryParams.Get(requestParamName)
		fieldValue := structValue.FieldByName(field.Name)

		switch field.Type {
		case nullStringType:
			structNullString := reflect.Indirect(fieldValue)
			if len(requestParamValue) > 0 {
				structNullString.FieldByName("String").SetString(requestParamValue)
				structNullString.FieldByName("Valid").SetBool(true)
			}
		case nullInt64Type:
			structNullInt64 := reflect.Indirect(fieldValue)
			int64v, err := strconv.ParseInt(requestParamValue, 10, 64)
			if err != nil {
				structNullInt64.FieldByName("Valid").SetBool(false)
			} else {
				structNullInt64.FieldByName("Int64").SetInt(int64v)
				structNullInt64.FieldByName("Valid").SetBool(true)
			}
		case NullTimeType:
			structNullTime := reflect.Indirect(fieldValue)
			i, err := strconv.ParseInt(requestParamValue, 10, 64)
			if err != nil {
				structNullTime.FieldByName("Valid").SetBool(false)
			} else {
				tm := time.Unix(i, 0).UTC()
				structNullTime.FieldByName("Time").Set(reflect.ValueOf(tm))
				structNullTime.FieldByName("Valid").SetBool(true)
			}
		}
	}

	return nil
}

package silvia

import(
	"time"
	"reflect"
	"net/url"
	"strconv"
	"strings"
	"database/sql"
	"encoding/json"

	"github.com/serenize/snaker"
)

type(
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
		Clt        time.Time      `db:"clicked_at"`
		It         time.Time      `db:"installed_at"`
		Ct         time.Time      `db:"event_at"`
		Rt         time.Time      `db:"reattributed_at"`
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
	}

	AdjustRequest struct {
		IPAddress   string    `json:"ip_addr"`
		TimeLocal   string    `json:"time_local"`
		Request     string    `json:"request"`
		Referer     string    `json:"http_referer"`
		Useragent   string    `json:"http_user_agent"`
	}

)

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
		request = []byte(adjustRequest.Request)[4:len(adjustRequest.Request)-9]
	}

	u, err := url.Parse(string(request))
	if err != nil {
		return err
	}

	queryParams := u.Query()

	structType := reflect.TypeOf(*event)
	structValue := reflect.ValueOf(event).Elem()

	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		fieldName := field.Name
		requestParamName := snaker.CamelToSnake(fieldName)
		requestParamValue := queryParams.Get(requestParamName)
		fieldValue := structValue.FieldByName(field.Name)

		nullStringType := reflect.TypeOf(sql.NullString{})
		nullInt64Type := reflect.TypeOf(sql.NullInt64{})
		timeType := reflect.TypeOf(time.Time{})

		switch field.Type {
		case nullStringType:
			structNullString := reflect.Indirect(fieldValue)
			structNullString.FieldByName("String").SetString(requestParamValue)
			if len(requestParamValue) > 0 { structNullString.FieldByName("Valid").SetBool(true) }
		case nullInt64Type:
			structNullString := reflect.Indirect(fieldValue)
			int64v, err := strconv.ParseInt(requestParamValue, 10, 64)
			if err != nil {
				structNullString.FieldByName("Valid").SetBool(false)
			} else {
				structNullString.FieldByName("Int64").SetInt(int64v)
				structNullString.FieldByName("Valid").SetBool(true)
			}
		case timeType:
			i, err := strconv.ParseInt(requestParamValue, 10, 64)
			if err == nil {
				tm := time.Unix(i, 0)
				fieldValue.Set(reflect.ValueOf(tm))
			}
		}
	}

    return nil
}

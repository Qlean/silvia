package main

import(
	"time"
)

type NginxTime struct {
	time.Time
}

func (nginxTime *NginxTime) UnmarshalJSON(b []byte) (err error) {
	if b[0] == '"' && b[len(b)-1] == '"' {
		b = b[1 : len(b)-1]
	}
	layout := "02/Jan/2006:15:04:05 -0700"
	nginxTime.Time, err = time.Parse(layout, string(b))
	return
}

type NginxLog struct {
	IPAddress   string    `json:"ip_addr"`
	TimeLocal   NginxTime `json:"time_local"`
	RequestBody Snowplow  `json:"request_body"`
	Referer     string    `json:"http_referer"`
	UserAgent   string    `json:"http_user_agent"`
}

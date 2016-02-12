package main

import()

type(
	AdjustEvent struct {
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
		Imp        impression_based	smallint
		Org        is_organic	smallint
		Rej        sql.NullString `db:"rejection_reason"`
		Adid       sql.NullString `db:"adjust_device_id"`
		Idfa       sql.NullString `db:"idfa"`
		Anid       sql.NullString `db:"android_id"`
		Mac        sql.NullString `db:"mac_sha1"`
		Mt         sql.NullString `db:"match_type"`
		Ref        sql.NullString `db:"referrer"`
		Ua         sql.NullString `db:"user_agent"`
		Ip         sql.NullString `db:"ip_address"`
		Clt        clicked_at	timestamp without time zone
		It         installed_at	timestamp without time zone
		Ct         event_at	timestamp without time zone
		Rt         reattributed_at	timestamp without time zone
		Reg        sql.NullString `db:"region"`
		C          sql.NullString `db:"country"`
		Lng        sql.NullString `db:"language"`
		Dn         sql.NullString `db:"device_name"`
		Dt         sql.NullString `db:"device_type"`
		Os         sql.NullString `db:"os_name"`
		Osv        sql.NullString `db:"os_version"`
		Env        sql.NullString `db:"environment"`
		Tre        tracking_enabled	smallint
		Tz         sql.NullString `db:"device_timezone"`
		Ls         last_session_sec	integer
		Cs         current_session_sec	integer
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
)

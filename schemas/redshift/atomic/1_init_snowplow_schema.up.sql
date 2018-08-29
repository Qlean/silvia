create schema if not exists atomic;

create table if not exists atomic.events
(
	app_id varchar(32) encode bytedict,
	platform varchar(32),
	etl_tstamp timestamp,
	collector_tstamp timestamp not null,
	dvce_tstamp timestamp,
	event varchar(128),
	event_id char(36) not null,
	v_tracker varchar(32) encode bytedict,
	user_id integer,
	user_ipaddress varchar(45),
	user_fingerprint varchar(50),
	domain_userid varchar(36),
	visit_num smallint,
	session_id char(36),
	geo_country varchar(120) encode runlength,
	geo_region char(64),
	geo_city varchar(75),
	geo_zipcode varchar(15),
	geo_latitude double precision encode runlength,
	geo_longitude double precision encode runlength,
	geo_region_name varchar(100),
	ip_isp varchar(100),
	ip_organization varchar(100),
	ip_domain varchar(100),
	ip_netspeed varchar(100),
	page_url varchar(3000),
	page_title varchar(2000),
	page_referrer varchar(65535),
	page_urlhost varchar(255),
	page_urlport integer,
	page_urlpath varchar(3000),
	page_urlquery varchar(6000),
	page_urlfragment varchar(3000),
	refr_urlhost varchar(255),
	refr_urlport integer,
	refr_urlpath varchar(6000),
	refr_urlquery varchar(6000),
	refr_urlfragment varchar(3000),
	utm_medium varchar(255),
	utm_source varchar(255),
	utm_term varchar(255),
	utm_content varchar(500),
	utm_campaign varchar(255),
	contexts varchar(65535),
	se_category varchar(1000),
	se_action varchar(1000),
	se_label varchar(1000),
	se_property varchar(1000),
	se_value double precision encode runlength,
	unstruct_event varchar(65535),
	useragent varchar(1000),
	br_name varchar(50),
	br_family varchar(50),
	br_version varchar(50),
	br_type varchar(50),
	br_renderengine varchar(50),
	br_lang varchar(32) encode bytedict,
	br_features_pdf boolean,
	br_features_flash boolean,
	br_features_java boolean encode runlength,
	br_features_director boolean encode runlength,
	br_features_quicktime boolean encode runlength,
	br_features_realplayer boolean encode runlength,
	br_features_windowsmedia boolean encode runlength,
	br_features_gears boolean encode runlength,
	br_features_silverlight boolean encode runlength,
	br_cookies boolean,
	br_colordepth varchar(12),
	br_viewwidth integer,
	br_viewheight integer,
	os_name varchar(50),
	os_family varchar(50),
	os_manufacturer varchar(50),
	os_timezone varchar(50),
	dvce_type varchar(128),
	dvce_ismobile boolean,
	dvce_screenwidth integer,
	dvce_screenheight integer,
	doc_charset varchar(32),
	doc_width integer,
	doc_height integer,
	geo_timezone varchar(64)
)
;

create or replace view atomic.v_space_used_per_tbl as
SELECT info.dbase_name, info.schemaname, info.tablename, info.tbl_oid, info.megabytes, info.rowcount, info.unsorted_rowcount, info.pct_unsorted, CASE WHEN (info.rowcount = 0) THEN 'n/a'::character varying WHEN (info.pct_unsorted >= (((20)::numeric)::numeric(18,0))::numeric(20,2)) THEN 'VACUUM SORT recommended'::character varying ELSE 'n/a'::character varying END AS recommendation FROM (SELECT btrim(((pgdb.datname)::character varying)::text) AS dbase_name, btrim(((pgn.nspname)::character varying)::text) AS schemaname, btrim(((a.name)::character varying)::text) AS tablename, a.id AS tbl_oid, b.mbytes AS megabytes, CASE WHEN (pgc.reldiststyle = 8) THEN a.rows_all_dist ELSE a."rows" END AS rowcount, CASE WHEN (pgc.reldiststyle = 8) THEN a.unsorted_rows_all_dist ELSE a.unsorted_rows END AS unsorted_rowcount, CASE WHEN (pgc.reldiststyle = 8) THEN (CASE WHEN ((det.n_sortkeys = 0) OR ((det.n_sortkeys IS NULL) AND (0 IS NULL))) THEN (NULL::numeric)::numeric(18,0) ELSE CASE WHEN ((a.rows_all_dist = 0) OR ((a.rows_all_dist IS NULL) AND (0 IS NULL))) THEN ((0)::numeric)::numeric(18,0) ELSE (((((a.unsorted_rows_all_dist)::numeric)::numeric(18,0))::numeric(32,0) / ((a.rows_all_dist)::numeric)::numeric(18,0)) * ((100)::numeric)::numeric(18,0)) END END)::numeric(20,2) ELSE (CASE WHEN ((det.n_sortkeys = 0) OR ((det.n_sortkeys IS NULL) AND (0 IS NULL))) THEN (NULL::numeric)::numeric(18,0) ELSE CASE WHEN ((a."rows" = 0) OR ((a."rows" IS NULL) AND (0 IS NULL))) THEN ((0)::numeric)::numeric(18,0) ELSE (((((a.unsorted_rows)::numeric)::numeric(18,0))::numeric(32,0) / ((a."rows")::numeric)::numeric(18,0)) * ((100)::numeric)::numeric(18,0)) END END)::numeric(20,2) END AS pct_unsorted FROM ((((((SELECT stv_tbl_perm.db_id, stv_tbl_perm.id, stv_tbl_perm.name, "max"(stv_tbl_perm."rows") AS rows_all_dist, ("max"(stv_tbl_perm."rows") - "max"(stv_tbl_perm.sorted_rows)) AS unsorted_rows_all_dist, sum(stv_tbl_perm."rows") AS "rows", (sum(stv_tbl_perm."rows") - sum(stv_tbl_perm.sorted_rows)) AS unsorted_rows FROM stv_tbl_perm GROUP BY stv_tbl_perm.db_id, stv_tbl_perm.id, stv_tbl_perm.name) a JOIN pg_class pgc ON ((pgc.oid = (a.id)::oid))) JOIN pg_namespace pgn ON ((pgn.oid = pgc.relnamespace))) JOIN pg_database pgdb ON ((pgdb.oid = (a.db_id)::oid))) JOIN (SELECT pg_attribute.attrelid, min(((CASE WHEN (pg_attribute.attisdistkey = true) THEN pg_attribute.attname ELSE NULL::name END)::character varying)::text) AS "distkey", min(((CASE WHEN (pg_attribute.attsortkeyord = 1) THEN pg_attribute.attname ELSE NULL::name END)::character varying)::text) AS head_sort, "max"(pg_attribute.attsortkeyord) AS n_sortkeys, "max"(pg_attribute.attencodingtype) AS max_enc, (((((sum(CASE WHEN (pg_attribute.attencodingtype <> 0) THEN 1 ELSE 0 END))::numeric)::numeric(18,0))::numeric(20,3) / (((count(pg_attribute.attencodingtype))::numeric)::numeric(18,0))::numeric(20,3)) * 100.00) AS pct_enc FROM pg_attribute GROUP BY pg_attribute.attrelid) det ON ((det.attrelid = (a.id)::oid))) LEFT JOIN (SELECT stv_blocklist.tbl, count(*) AS mbytes FROM stv_blocklist GROUP BY stv_blocklist.tbl) b ON ((a.id = b.tbl))) WHERE (pgc.relowner > 1)) info;
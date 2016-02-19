# Silvia
Transforms analytics trackers data (Snowplow and Adjust) and writes it to PostgreSQL.

## Dependencies

- Nginx version > 1.7.1 (with syslog support)
- [Mozilla's Heka](https://github.com/mozilla-services/heka) (With SysLog inputs, decoders and AMQP Outputs)
- RabbitMQ (with "snowplow" and "adjust" queues and exchanges)
- Consul (for configuration, health checking and service discovery)
- PostgreSQL (events strorage)
- `GeoLiteCity.dat` database file in same folder with application

## Deployment
![Infrastructure visualistaion](images/deployment.png)

## Configuration

### Consul

You're should provide following keys in consul Key-Value storage:

- `silvia/pg_connect` - PostgreSQL connection configuration (e.g. `dbname=analytics host=postgres user=postgres password=secretpass`)
- `silvia/port` - API port binding
- `silvia/rabbit_addr` - RabbitMQ address
- `silvia/rabbit_addr` - RabbitMQ port
- `silvia/ring_size` - Ring buffers size (e.g. `10`)

### Nginx

```
log_format snowplow "{\x22ip_addr\x22:\x22$remote_addr\x22,\x22time_local\x22:\x22$time_local\x22,\x22request_body\x22:$request_body,\x22http_referer\x22:\x22$http_referer\x22,\x22http_user_agent\x22:\x22$http_user_agent\x22}";
log_format adjust "$request_uri";

server {
	listen 80

	location /empty_proxy {
		access_log off;
		return 200;
	}

	location /stats/snowplow {
		proxy_set_header ClientAddress $proxy_add_x_forwarded_for;
		proxy_redirect off;
		proxy_pass http://127.0.0.1/empty_proxy;
		add_header Access-Control-Allow-Origin "*";
		client_max_body_size 4k;
		client_body_buffer_size 4k;
		client_body_in_single_buffer on;
		access_log syslog:server=localhost:8514,facility=local0 snowplow;
	}

	location /stats/adjust {
		proxy_set_header ClientAddress $proxy_add_x_forwarded_for;
		proxy_redirect off;
		proxy_pass http://127.0.0.1/empty_proxy;
		add_header Access-Control-Allow-Origin "*";
		client_max_body_size 4k;
		client_body_buffer_size 4k;
		client_body_in_single_buffer on;
		access_log syslog:server=localhost:8515,facility=local0 adjust;
	}
}
```

### Heka

```
[PayloadEncoder]
append_newlines = false

[AdjustSyslogInput]
type="UdpInput"
address = "127.0.0.1:8515"
splitter = "NullSplitter"
decoder = "AdjustSyslogDecoder"

[AdjustSyslogDecoder]
type = "SandboxDecoder"
filename = "lua_decoders/rsyslog.lua"

[AdjustSyslogDecoder.config]
type = "RSYSLOG_TraditionalForwardFormat"
template = '<%PRI%>%TIMESTAMP% %HOSTNAME% %syslogtag:1:32%%msg:::sp-if-no-1st-sp%%msg%'

[AdjustAMQPOutput]
type="AMQPOutput"
url = "amqp://guest:guest@rabbitmq.service.consul/"
message_matcher = "Logger == 'AdjustSyslogInput' && Fields[programname] == 'nginx'"
exchange = "adjust"
routing_key = "adjust"
exchange_type = "fanout"
exchange_auto_delete = false
exchange_durability = true
encoder = "PayloadEncoder"
use_framing = false

[SnowplowSyslogInput]
type="UdpInput"
address = "127.0.0.1:8514"
splitter = "NullSplitter"
decoder = "SnowplowSyslogDecoder"

[SnowplowSyslogDecoder]
type = "SandboxDecoder"
filename = "lua_decoders/rsyslog.lua"

[SnowplowSyslogDecoder.config]
type = "RSYSLOG_TraditionalForwardFormat"
template = '<%PRI%>%TIMESTAMP% %HOSTNAME% %syslogtag:1:32%%msg:::sp-if-no-1st-sp%%msg%'

[SnowplowAMQPOutput]
type="AMQPOutput"
url = "amqp://guest:guest@rabbitmq.service.consul/"
message_matcher = "Logger == 'SnowplowSyslogInput' && Fields[programname] == 'nginx'"
exchange = "snowplow"
routing_key = "snowplow"
exchange_type = "fanout"
exchange_auto_delete = false
exchange_durability = true
encoder = "PayloadEncoder"
use_framing = false
```

## App structure
![Goroutines communication](images/goroutines.png)

## API endpoints

`/v1/status` - Application summary statistics. Returns `200` if healthy, `429` if not.

`/v1/ring?tracker=TRACKER_NAME&ring=RINGNAME` - Display particular ring buffer. Parameters: tracker (`snowplow` or `adjust`), ring (`success` or `failed`)

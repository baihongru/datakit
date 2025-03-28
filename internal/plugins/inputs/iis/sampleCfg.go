// Unless explicitly stated otherwise all files in this repository are licensed
// under the MIT License.
// This product includes software developed at Guance Cloud (https://www.guance.com/).
// Copyright 2021-present Guance, Inc.

package iis

//nolint:lll
const (
	sampleConfig = `
[[inputs.iis]]
  ## (optional) collect interval, default is 15 seconds
  interval = '15s'
  ##

  [inputs.iis.log]
    files = []
    ## grok pipeline script path
    pipeline = "iis.p"

  [inputs.iis.tags]
    ## tag1 = "v1"
    ## tag2 = "v2"
`

	//nolint:deadcode
	pipelineCfg = `
grok(_, "%{TIMESTAMP_ISO8601:time} %{IP:server_ip} %{DATA:http_method} %{DATA:http_url} %{DATA:url_param} %{NUMBER:port} %{DATA:username} %{IP:client_ip} %{DATA:user_agent} %{DATA:referer} %{NUMBER:status_code} %{NUMBER:sub_status} %{NUMBER:win32_status} %{NUMBER:time_taken}")

cast(port, "int")
cast(status_code, "int")
cast(sub_status, "int")
cast(win32_status, "int")
cast(time_taken, "int")

group_between(status_code, [200,299], "ok", status)
group_between(status_code, [300,399], "notice", status)
group_between(status_code, [400,499], "warning", status)
group_between(status_code, [500,599], "error", status)

nullif(url_param, "-")
nullif(username, "-")
nullif(referer, "-")

default_time(time, "UTC")`
)

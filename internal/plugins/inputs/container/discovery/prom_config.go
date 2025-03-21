// Unless explicitly stated otherwise all files in this repository are licensed
// under the MIT License.
// This product includes software developed at Guance Cloud (https://www.guance.com/).
// Copyright 2021-present Guance, Inc.

package discovery

import (
	"fmt"
	"net/url"
	"time"

	bstoml "github.com/BurntSushi/toml"
	iprom "gitlab.jiagouyun.com/cloudcare-tools/datakit/internal/prom"
)

const (
	annotationPromExport = "datakit/prom.instances"
	defaultInterval      = time.Second * 30
)

type promConfig struct {
	Source   string        `toml:"source" json:"source"`
	Interval time.Duration `toml:"interval"`
	Timeout  time.Duration `toml:"timeout"`

	URL  string   `toml:"url" json:"url"` // deprecated
	URLs []string `toml:"urls" json:"urls"`

	IgnoreReqErr           bool         `toml:"ignore_req_err" json:"ignore_req_err"`
	MetricTypes            []string     `toml:"metric_types" json:"metric_types"`
	MetricNameFilter       []string     `toml:"metric_name_filter" json:"metric_name_filter"`
	MetricNameFilterIgnore []string     `toml:"metric_name_filter_ignore" json:"metric_name_filter_ignore"`
	MeasurementPrefix      string       `toml:"measurement_prefix" json:"measurement_prefix"`
	MeasurementName        string       `toml:"measurement_name" json:"measurement_name"`
	Measurements           []iprom.Rule `toml:"measurements" json:"measurements"`

	TLSOpen            bool   `toml:"tls_open" json:"tls_open"`
	UDSPath            string `toml:"uds_path" json:"uds_path"`
	CacertFile         string `toml:"tls_ca"`
	CertFile           string `toml:"tls_cert"`
	KeyFile            string `toml:"tls_key"`
	InsecureSkipVerify bool   `toml:"insecure_skip_verify" json:"insecure_skip_verify"`
	BearerTokenFile    string `toml:"bearer_token_file" json:"bearer_token_file"`

	TagsIgnore  []string            `toml:"tags_ignore" json:"tags_ignore"`
	TagsRename  *iprom.RenameTags   `toml:"tags_rename" json:"tags_rename"`
	AsLogging   *iprom.AsLogging    `toml:"as_logging" json:"as_logging"`
	IgnoreTagKV map[string][]string `toml:"ignore_tag_kv_match" json:"ignore_tag_kv_match"`
	HTTPHeaders map[string]string   `toml:"http_headers" json:"http_headers"`

	Tags           map[string]string
	DisableInfoTag bool `toml:"disable_info_tag" json:"disable_info_tag"`

	Auth map[string]string `toml:"auth" json:"auth"`
}

type promOption func(c *promConfig)

func withLabelAsTags(m map[string]string, keys []string) promOption {
	return func(c *promConfig) {
		if len(keys) == 0 || len(m) == 0 {
			return
		}
		for _, key := range keys {
			if v, ok := m[key]; ok {
				withTag(replaceLabelKey(key), v)(c)
			}
		}
	}
}

func withTag(key, value string) promOption {
	return func(c *promConfig) {
		if c.Tags == nil {
			c.Tags = make(map[string]string)
		}
		if _, ok := c.Tags[key]; !ok {
			c.Tags[key] = value
		}
	}
}

type wrapPromConfig struct {
	Inputs struct {
		Prom []*promConfig `toml:"prom"`
	} `toml:"inputs"`
}

func parseURLHost(cfg *promConfig) (map[string]string, error) {
	res := make(map[string]string)
	for _, urlstr := range cfg.URLs {
		u, err := url.Parse(urlstr)
		if err != nil {
			return nil, fmt.Errorf("invalid url %s, err: %w", urlstr, err)
		}
		res[urlstr] = u.Host
	}
	return res, nil
}

func parsePromConfigs(str string) ([]*promConfig, error) {
	c := wrapPromConfig{}
	if err := bstoml.Unmarshal([]byte(str), &c); err != nil {
		return nil, fmt.Errorf("unable to parse toml: %w", err)
	}
	for _, cfg := range c.Inputs.Prom {
		if cfg.URL != "" {
			cfg.URLs = append(cfg.URLs, cfg.URL)
		}
		if cfg.Interval <= 0 {
			cfg.Interval = defaultInterval
		}
	}
	return c.Inputs.Prom, nil
}

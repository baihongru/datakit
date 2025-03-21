// Unless explicitly stated otherwise all files in this repository are licensed
// under the MIT License.
// This product includes software developed at Guance Cloud (https://www.guance.com/).
// Copyright 2021-present Guance, Inc.

package postgresql

import (
	"github.com/GuanceCloud/cliutils/point"
	"gitlab.jiagouyun.com/cloudcare-tools/datakit/internal/plugins/inputs"
)

type customerObjectMeasurement struct {
	name   string
	tags   map[string]string
	fields map[string]interface{}
	ipt    *Input
}

// Point implement MeasurementV2.
func (m *customerObjectMeasurement) Point() *point.Point {
	return point.NewPointV2(m.name,
		append(point.NewTags(m.tags), point.NewKVs(m.fields)...),
		append(point.DefaultObjectOptions(), point.WithExtraTags(m.ipt.mergedTags))...)
}

//nolint:lll
func (m *customerObjectMeasurement) Info() *inputs.MeasurementInfo { //nolint:funlen
	return &inputs.MeasurementInfo{
		Name: "database",
		Type: "custom_object",
		Fields: map[string]interface{}{
			"uptime": &inputs.FieldInfo{
				DataType: inputs.Int,
				Type:     inputs.Gauge,
				Unit:     inputs.DurationSecond,
				Desc:     "Current PostgreSQL uptime",
			},

			"display_name": &inputs.FieldInfo{
				DataType: inputs.String,
				Type:     inputs.Gauge,
				Unit:     inputs.UnknownUnit,
				Desc:     "Displayed name in UI",
			},

			"version": &inputs.FieldInfo{
				DataType: inputs.String,
				Type:     inputs.Gauge,
				Unit:     inputs.UnknownUnit,
				Desc:     "Current version of PostgreSQL",
			},
		},
		Tags: map[string]interface{}{
			"name": &inputs.TagInfo{
				Desc: "Object uniq ID",
			},

			"col_co_status": &inputs.TagInfo{
				Desc: "Current status of collector on PostgreSQL(`OK/NotOK`)",
			},

			"ip": &inputs.TagInfo{
				Desc: "",
			},

			"host": &inputs.TagInfo{
				Desc: "The server host address",
			},
			"reason": &inputs.TagInfo{
				Desc: "If status not ok, we'll get some reasons about the status",
			},
		},
	}
}

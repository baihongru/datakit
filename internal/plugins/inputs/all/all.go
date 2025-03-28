// Unless explicitly stated otherwise all files in this repository are licensed
// under the MIT License.
// This product includes software developed at Guance Cloud (https://www.guance.com/).
// Copyright 2021-present Guance, Inc.

//go:build !datakit_elinker && !datakit_lite && !datakit_aws_lambda && with_inputs
// +build !datakit_elinker,!datakit_lite,!datakit_aws_lambda,with_inputs

// Package inputs wraps all inputs implements
package inputs

import (
	// default inputs.
	_ "gitlab.jiagouyun.com/cloudcare-tools/datakit/internal/plugins/inputs/cpu"
	_ "gitlab.jiagouyun.com/cloudcare-tools/datakit/internal/plugins/inputs/dialtesting"
	_ "gitlab.jiagouyun.com/cloudcare-tools/datakit/internal/plugins/inputs/disk"
	_ "gitlab.jiagouyun.com/cloudcare-tools/datakit/internal/plugins/inputs/diskio"
	_ "gitlab.jiagouyun.com/cloudcare-tools/datakit/internal/plugins/inputs/dk"
	_ "gitlab.jiagouyun.com/cloudcare-tools/datakit/internal/plugins/inputs/hostobject"
	_ "gitlab.jiagouyun.com/cloudcare-tools/datakit/internal/plugins/inputs/mem"
	_ "gitlab.jiagouyun.com/cloudcare-tools/datakit/internal/plugins/inputs/net"
	_ "gitlab.jiagouyun.com/cloudcare-tools/datakit/internal/plugins/inputs/swap"
	_ "gitlab.jiagouyun.com/cloudcare-tools/datakit/internal/plugins/inputs/system"

	_ "gitlab.jiagouyun.com/cloudcare-tools/datakit/internal/plugins/inputs/apache"
	_ "gitlab.jiagouyun.com/cloudcare-tools/datakit/internal/plugins/inputs/awslambda"
	_ "gitlab.jiagouyun.com/cloudcare-tools/datakit/internal/plugins/inputs/beats_output"
	_ "gitlab.jiagouyun.com/cloudcare-tools/datakit/internal/plugins/inputs/cassandra"
	_ "gitlab.jiagouyun.com/cloudcare-tools/datakit/internal/plugins/inputs/cat"
	_ "gitlab.jiagouyun.com/cloudcare-tools/datakit/internal/plugins/inputs/chrony"
	_ "gitlab.jiagouyun.com/cloudcare-tools/datakit/internal/plugins/inputs/clickhousev1"
	_ "gitlab.jiagouyun.com/cloudcare-tools/datakit/internal/plugins/inputs/cloudprober"
	_ "gitlab.jiagouyun.com/cloudcare-tools/datakit/internal/plugins/inputs/cockroachdb"
	_ "gitlab.jiagouyun.com/cloudcare-tools/datakit/internal/plugins/inputs/consul"
	_ "gitlab.jiagouyun.com/cloudcare-tools/datakit/internal/plugins/inputs/container"
	_ "gitlab.jiagouyun.com/cloudcare-tools/datakit/internal/plugins/inputs/coredns"
	_ "gitlab.jiagouyun.com/cloudcare-tools/datakit/internal/plugins/inputs/couchbase"
	_ "gitlab.jiagouyun.com/cloudcare-tools/datakit/internal/plugins/inputs/couchdb"
	_ "gitlab.jiagouyun.com/cloudcare-tools/datakit/internal/plugins/inputs/db2"
	_ "gitlab.jiagouyun.com/cloudcare-tools/datakit/internal/plugins/inputs/ddtrace"
	_ "gitlab.jiagouyun.com/cloudcare-tools/datakit/internal/plugins/inputs/demo"
	_ "gitlab.jiagouyun.com/cloudcare-tools/datakit/internal/plugins/inputs/doris"
	_ "gitlab.jiagouyun.com/cloudcare-tools/datakit/internal/plugins/inputs/ebpf"
	_ "gitlab.jiagouyun.com/cloudcare-tools/datakit/internal/plugins/inputs/ebpftrace"
	_ "gitlab.jiagouyun.com/cloudcare-tools/datakit/internal/plugins/inputs/elasticsearch"
	_ "gitlab.jiagouyun.com/cloudcare-tools/datakit/internal/plugins/inputs/etcd"
	_ "gitlab.jiagouyun.com/cloudcare-tools/datakit/internal/plugins/inputs/external"
	_ "gitlab.jiagouyun.com/cloudcare-tools/datakit/internal/plugins/inputs/flinkv1"
	_ "gitlab.jiagouyun.com/cloudcare-tools/datakit/internal/plugins/inputs/gitlab"
	_ "gitlab.jiagouyun.com/cloudcare-tools/datakit/internal/plugins/inputs/graphite"
	_ "gitlab.jiagouyun.com/cloudcare-tools/datakit/internal/plugins/inputs/host_healthcheck"
	_ "gitlab.jiagouyun.com/cloudcare-tools/datakit/internal/plugins/inputs/hostdir"
	_ "gitlab.jiagouyun.com/cloudcare-tools/datakit/internal/plugins/inputs/influxdb"
	_ "gitlab.jiagouyun.com/cloudcare-tools/datakit/internal/plugins/inputs/ipmi"
	_ "gitlab.jiagouyun.com/cloudcare-tools/datakit/internal/plugins/inputs/jaeger"
	_ "gitlab.jiagouyun.com/cloudcare-tools/datakit/internal/plugins/inputs/jenkins"
	_ "gitlab.jiagouyun.com/cloudcare-tools/datakit/internal/plugins/inputs/jvm"
	_ "gitlab.jiagouyun.com/cloudcare-tools/datakit/internal/plugins/inputs/kafka"
	_ "gitlab.jiagouyun.com/cloudcare-tools/datakit/internal/plugins/inputs/kafkamq"
	_ "gitlab.jiagouyun.com/cloudcare-tools/datakit/internal/plugins/inputs/kubernetesprometheus"
	_ "gitlab.jiagouyun.com/cloudcare-tools/datakit/internal/plugins/inputs/logfwdserver"
	_ "gitlab.jiagouyun.com/cloudcare-tools/datakit/internal/plugins/inputs/logging"
	_ "gitlab.jiagouyun.com/cloudcare-tools/datakit/internal/plugins/inputs/logstreaming"
	_ "gitlab.jiagouyun.com/cloudcare-tools/datakit/internal/plugins/inputs/lsblk"
	_ "gitlab.jiagouyun.com/cloudcare-tools/datakit/internal/plugins/inputs/memcached"
	_ "gitlab.jiagouyun.com/cloudcare-tools/datakit/internal/plugins/inputs/mongodb"
	_ "gitlab.jiagouyun.com/cloudcare-tools/datakit/internal/plugins/inputs/mysql"
	_ "gitlab.jiagouyun.com/cloudcare-tools/datakit/internal/plugins/inputs/neo4j"
	_ "gitlab.jiagouyun.com/cloudcare-tools/datakit/internal/plugins/inputs/netflow"
	_ "gitlab.jiagouyun.com/cloudcare-tools/datakit/internal/plugins/inputs/netstat"
	_ "gitlab.jiagouyun.com/cloudcare-tools/datakit/internal/plugins/inputs/newrelic"
	_ "gitlab.jiagouyun.com/cloudcare-tools/datakit/internal/plugins/inputs/nfs"
	_ "gitlab.jiagouyun.com/cloudcare-tools/datakit/internal/plugins/inputs/nginx"
	_ "gitlab.jiagouyun.com/cloudcare-tools/datakit/internal/plugins/inputs/nsq"
	_ "gitlab.jiagouyun.com/cloudcare-tools/datakit/internal/plugins/inputs/nvidiasmi"
	_ "gitlab.jiagouyun.com/cloudcare-tools/datakit/internal/plugins/inputs/oceanbase"
	_ "gitlab.jiagouyun.com/cloudcare-tools/datakit/internal/plugins/inputs/opentelemetry"
	_ "gitlab.jiagouyun.com/cloudcare-tools/datakit/internal/plugins/inputs/oracle"
	_ "gitlab.jiagouyun.com/cloudcare-tools/datakit/internal/plugins/inputs/pinpoint"
	_ "gitlab.jiagouyun.com/cloudcare-tools/datakit/internal/plugins/inputs/ploffload"
	_ "gitlab.jiagouyun.com/cloudcare-tools/datakit/internal/plugins/inputs/postgresql"
	_ "gitlab.jiagouyun.com/cloudcare-tools/datakit/internal/plugins/inputs/process"
	_ "gitlab.jiagouyun.com/cloudcare-tools/datakit/internal/plugins/inputs/profile"
	_ "gitlab.jiagouyun.com/cloudcare-tools/datakit/internal/plugins/inputs/prom"
	_ "gitlab.jiagouyun.com/cloudcare-tools/datakit/internal/plugins/inputs/promremote"
	_ "gitlab.jiagouyun.com/cloudcare-tools/datakit/internal/plugins/inputs/promtail"
	_ "gitlab.jiagouyun.com/cloudcare-tools/datakit/internal/plugins/inputs/promv2"
	_ "gitlab.jiagouyun.com/cloudcare-tools/datakit/internal/plugins/inputs/proxy"
	_ "gitlab.jiagouyun.com/cloudcare-tools/datakit/internal/plugins/inputs/pushgateway"
	_ "gitlab.jiagouyun.com/cloudcare-tools/datakit/internal/plugins/inputs/pythond"
	_ "gitlab.jiagouyun.com/cloudcare-tools/datakit/internal/plugins/inputs/rabbitmq"
	_ "gitlab.jiagouyun.com/cloudcare-tools/datakit/internal/plugins/inputs/redis"
	_ "gitlab.jiagouyun.com/cloudcare-tools/datakit/internal/plugins/inputs/rum"
	_ "gitlab.jiagouyun.com/cloudcare-tools/datakit/internal/plugins/inputs/sensors"
	_ "gitlab.jiagouyun.com/cloudcare-tools/datakit/internal/plugins/inputs/skywalking"
	_ "gitlab.jiagouyun.com/cloudcare-tools/datakit/internal/plugins/inputs/smart"
	_ "gitlab.jiagouyun.com/cloudcare-tools/datakit/internal/plugins/inputs/snmp"
	_ "gitlab.jiagouyun.com/cloudcare-tools/datakit/internal/plugins/inputs/socket"
	_ "gitlab.jiagouyun.com/cloudcare-tools/datakit/internal/plugins/inputs/solr"
	_ "gitlab.jiagouyun.com/cloudcare-tools/datakit/internal/plugins/inputs/sqlserver"
	_ "gitlab.jiagouyun.com/cloudcare-tools/datakit/internal/plugins/inputs/ssh"
	_ "gitlab.jiagouyun.com/cloudcare-tools/datakit/internal/plugins/inputs/statsd"
	_ "gitlab.jiagouyun.com/cloudcare-tools/datakit/internal/plugins/inputs/tdengine"
	_ "gitlab.jiagouyun.com/cloudcare-tools/datakit/internal/plugins/inputs/tomcat"
	_ "gitlab.jiagouyun.com/cloudcare-tools/datakit/internal/plugins/inputs/vsphere"
	_ "gitlab.jiagouyun.com/cloudcare-tools/datakit/internal/plugins/inputs/zabbix_exporter"
	_ "gitlab.jiagouyun.com/cloudcare-tools/datakit/internal/plugins/inputs/zipkin"

	// only windows.
	_ "gitlab.jiagouyun.com/cloudcare-tools/datakit/internal/plugins/inputs/iis"
	_ "gitlab.jiagouyun.com/cloudcare-tools/datakit/internal/plugins/inputs/winevent"
)

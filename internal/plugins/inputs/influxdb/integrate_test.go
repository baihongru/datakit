// Unless explicitly stated otherwise all files in this repository are licensed
// under the MIT License.
// This product includes software developed at Guance Cloud (https://www.guance.com/).
// Copyright 2021-present Guance, Inc.

package influxdb

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/GuanceCloud/cliutils/point"
	dockertest "github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/require"
	"gitlab.jiagouyun.com/cloudcare-tools/datakit/internal/io"
	"gitlab.jiagouyun.com/cloudcare-tools/datakit/internal/plugins/inputs"
	"gitlab.jiagouyun.com/cloudcare-tools/datakit/internal/testutils"
)

// ATTENTION: Docker version should use v20.10.18 in integrate tests. Other versions are not tested.

var mExpect = map[string]struct{}{
	influxdbCq:            {},
	influxdbHttpd:         {},
	influxdbMemstats:      {},
	influxdbQueryExecutor: {},
	influxdbRuntime:       {},
	influxdbSubscriber:    {},
	influxdbWrite:         {},
}

func TestIntegrate(t *testing.T) {
	if !testutils.CheckIntegrationTestingRunning() {
		t.Skip()
	}

	testutils.PurgeRemoteByName(inputName)       // purge at first.
	defer testutils.PurgeRemoteByName(inputName) // purge at last.

	start := time.Now()
	cases, err := buildCases(t)
	if err != nil {
		cr := &testutils.CaseResult{
			Name:          t.Name(),
			Status:        testutils.TestPassed,
			FailedMessage: err.Error(),
			Cost:          time.Since(start),
		}

		_ = testutils.Flush(cr)
		return
	}

	t.Logf("testing %d cases...", len(cases))

	for _, tc := range cases {
		func(tc *caseSpec) {
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()
				caseStart := time.Now()

				t.Logf("testing %s...", tc.name)

				if err := testutils.RetryTestRun(tc.run); err != nil {
					tc.cr.Status = testutils.TestFailed
					tc.cr.FailedMessage = err.Error()

					panic(err)
				} else {
					tc.cr.Status = testutils.TestPassed
				}

				tc.cr.Cost = time.Since(caseStart)

				require.NoError(t, testutils.Flush(tc.cr))

				t.Cleanup(func() {
					// clean remote docker resources
					if tc.resource == nil {
						return
					}

					tc.pool.Purge(tc.resource)
				})
			})
		}(tc)
	}
}

func getConfAccessPoint(scheme, host, port string) string {
	return fmt.Sprintf("%s://%s/debug/vars", scheme, net.JoinHostPort(host, port))
}

func buildCases(t *testing.T) ([]*caseSpec, error) {
	t.Helper()

	remote := testutils.GetRemote()

	bases := []struct {
		name                      string // Also used as build image name:tag.
		conf                      string
		dockerFileText            string // Empty if not build image.
		exposedPorts              []string
		cmd                       []string
		optsInfluxdbCq            []inputs.PointCheckOption
		optsInfluxdbHttpd         []inputs.PointCheckOption
		optsInfluxdbMemstats      []inputs.PointCheckOption
		optsInfluxdbQueryExecutor []inputs.PointCheckOption
		optsInfluxdbRuntime       []inputs.PointCheckOption
		optsInfluxdbSubscriber    []inputs.PointCheckOption
		optsInfluxdbWrite         []inputs.PointCheckOption
	}{
		////////////////////////////////////////////////////////////////////////
		// influxdb:1.8.10 TLS base64
		////////////////////////////////////////////////////////////////////////
		{
			name: "pubrepo.jiagouyun.com/image-repo-for-testing/influxdb:1.8.10-alpine-tls",
			conf: `url = ""
			interval = '2s'
			timeout = "5s"
			election = false
			[tlsconf]
			  insecure_skip_verify = true
			  cert_base64 = "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSURJVENDQWdtZ0F3SUJBZ0lSQU1Id0txODJuN1NBWTNzdE1ieERxNlF3RFFZSktvWklodmNOQVFFTEJRQXcKS3pFU01CQUdBMVVFQ2hNSlEyOWphM0p2WVdOb01SVXdFd1lEVlFRREV3eERiMk5yY205aFkyZ2dRMEV3SGhjTgpNalF3TlRFd01ETXpORFUzV2hjTk1qa3dOVEUxTURNek5EVTNXakFqTVJJd0VBWURWUVFLRXdsRGIyTnJjbTloClkyZ3hEVEFMQmdOVkJBTVRCSEp2YjNRd2dnRWlNQTBHQ1NxR1NJYjNEUUVCQVFVQUE0SUJEd0F3Z2dFS0FvSUIKQVFETmpIOWJ2aitLL2wzMkxNRENDVUo2cTh2NEhGdTZZd3ZDdUZuU1NVWHVsbWZ3ZVA2M3RseFRNY0RGT09LNQpCQWdIWitUREJXQjcvK0tCS1cxVElhY1B0aFRoV0sxSlRtUUJOd3dPRmgrbnMzaWRxaHh0bHNsK3N3enJLWC83Ckl5MG8wdUs5QW8wVzZyTEl0ZlpjN0l0ak1ZbEpOczBmUUFKOWhDYmRFRllWSi9iODh0ZTJLcUkxU3dqa21VZ0sKZEtobTJBbjlYMmpRY2pIRml6TlA1V0R4RlUzUzJBUWdMdEN5elhDMFFvSlJoMW4zUFVtaldaQ2dtUEJrNnM4SgpJSE5oYTNkZWN5ZnZhVmhBVWJObkcwbmxHcDdJNXNRcHkydkpwenp1Y25xd1F2ZWp6c1NaTi81djFTZUFlbmQyCmpnalIrc1RoektKVDRrOXlDREdtb1kwVEFnTUJBQUdqU0RCR01BNEdBMVVkRHdFQi93UUVBd0lGb0RBVEJnTlYKSFNVRUREQUtCZ2dyQmdFRkJRY0RBakFmQmdOVkhTTUVHREFXZ0JSVXN6U0tacmYrbGt5Sm12aUxHTXMzRkVveAovREFOQmdrcWhraUc5dzBCQVFzRkFBT0NBUUVBR2hOVDhJS2RnRlpMUkNKbE93SHlHbC9LL0lYeFRkdUpNUTBICndtZ1FhTFgzRmlBK1g4d2N0NjNIMXExaXNTQjdsbVhkVFdBQXdzQThMTTBtWlJOOHUzWUhjYXVsT29xLzhBSVkKNVNHdXZ4QmEvdkJic2pCbndLVHU0SUZtRGtVYld3UThaamdwOUpoUU5NTmdwS2JjL2EwL3ZsWFVuOS8rTmJsegpiZGxxc0lYd05qZ2hpeGxsZDNWSXRlOXgreFdRb0Q4cldTUFVEcldmMFVNbGtVTDI4UG44WU41V0I4SHhJVmQwCkJsanorSjJqL0VBOC9HNUZnUkFMVVE2V1ROdmJRMkJzdyszTVhiOFE1Q1pQL3FCRUMrTys0RnZmb2NGL2VvRlUKWnZFbEsvSHBPV3VVQXF2blZVcUpFTkNKUnZkeTNOSHpRcTdpcnhaRHY0dDY3NTBidWc9PQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg=="

              cert_key_base64 = "LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlFcEFJQkFBS0NBUUVBell4L1c3NC9pdjVkOWl6QXdnbENlcXZMK0J4YnVtTUx3cmhaMGtsRjdwWm44SGorCnQ3WmNVekhBeFRqaXVRUUlCMmZrd3dWZ2UvL2lnU2x0VXlHbkQ3WVU0Vml0U1U1a0FUY01EaFlmcDdONG5hb2MKYlpiSmZyTU02eWwvK3lNdEtOTGl2UUtORnVxeXlMWDJYT3lMWXpHSlNUYk5IMEFDZllRbTNSQldGU2YyL1BMWAp0aXFpTlVzSTVKbElDblNvWnRnSi9WOW8wSEl4eFlzelQrVmc4UlZOMHRnRUlDN1FzczF3dEVLQ1VZZFo5ejFKCm8xbVFvSmp3Wk9yUENTQnpZV3QzWG5NbjcybFlRRkd6Wnh0SjVScWV5T2JFS2N0cnlhYzg3bko2c0VMM284N0UKbVRmK2I5VW5nSHAzZG80STBmckU0Y3lpVStKUGNnZ3hwcUdORXdJREFRQUJBb0lCQVFDSjNOWCsvcGMzN211dgpGVTBqMTNvVE5PN1ZOby8vYnpjUUh2MS9vVTJhUEo3eUZ2VWcydHNKb2JFZGxvM2FjZTNBcWRveFE0WDNKU1VTClpHckMreXRGeW1ZdXpuOUxUNXliaEFROTNuRFUxZmJzS0pCd29GWDgrTEtOZDRRek9PQ3RKT1NXeVFOQWY2SHkKSkxsY2tmcmJTUG8vZE5ZWFE2Tm45QjdzM213ZU96N2VML29xTXNHZTNDdHFVRDEycUFLK3pqVDFSL3pJL1FVSwpvUFFTcnllbm5saTg4eEM2UXpTenZ4UmRNSStwZGk1aWNENVJGMjRqdWk5UlZyYlk5VnlCVzNOVmRKci9MUzRJClRCWnphUmVTcGhiOFNKQ2ExMWRJblRUTlZGZWNEdkhNM1ZVb3RQYW9DRlZEQ2hFaG5ST1FJc0xxcXhLY2JmcmMKZks1bGswWEJBb0dCQU5GOUFGcUliOGltMHgwUnZ2OVdRUkVTRE1NWG5hTlIzTVZQaG5PdU5seDZUd3BHY2RMdwpkanZBenhXbzltQ1ZXUDBWMWQrY1BlVkxodmhEbmMwcWJWZUR3SW56dW8rR0F2UTEvbXFGYXVTbkdPRE95eVh0Ck9xRGNFRlQ2UmdqcmNoQXBnZG4xY3N5MVpLVUtQUkFSU0hSZStLRFkreERmQUlOZExKUzJ2MFk1QW9HQkFQc3YKbFNtK0x1TnNzN3RzWUI3MnYrc1FaREdEeExIWnFSN0xTeTRVZzZaZzQwbmptT3ZObVRzOUNBWFpXR2RnMzQyaQpyUDYzOENJcmsvMXJ5aWNIYTMwRWpJdWRhZUFlUHZlRTA2Y3YwNW9JVXJXbGY5VWFLRjRXM2Y5T0VQa1IzVWx3ClNjaDUwempmMWVMcE82SGt6ZjZhSjhyck9KbGJXSnh0UHVKR3NzMnJBb0dBRlVUNXlqZGNFaVZOL2YrVlF0dUIKRTdpZmJ4ZHd1K3BOM2dLckJnZkVJVE9SM3RzMEoxU2V6SVpSQUVQOWIrVDUrZ2hEaE1hYVNqT1c2cElDN1pmSApMa0dFUlAxb0RiWnZpbGdKRXN1bEJMNHFlbmpFaTM5QW1xQjlVQU54Sk9xeTFBMUN6OXhwNFhyeFV3aHRGcnFLCmZyWTl6Q2I3cHNUZGluamxVOXdTSTVFQ2dZRUFsRlNrQkROLzR1TTFPKytpejRZdERUWHZ4T0dvVE5KWk1Zc2gKaVVPcC9wMW1leUxCRWphbVR6b2FPOEgrbDRXNFhoNTdoQ3ZBelp6b1ZwWEptY1NpNy8rNHMxV3d5UjF6VjUyRAprMDRGNmdjU09KeFQ0ZGNCa1paMVlDZU1sRmk5VVhuU3lHVlFtMXhySlFWUUpxbEVFQjZlY3hEMnFuRXI0YXdOCm4zZmFiT01DZ1lCSDM5MmM3MVVCOUVxak9EV0J3bDNHeEhIaUJUZCtQOWdIODdDTExGUlNKZk9nR2ppV2JXVTYKb2M0WVg4ZEFoSTlwdkJUazhOM3NBalg3ano1N1lOblpDUmwxOURFc3RJV2dEaUhWZWN6d0VyV0ltV0ZEZ3Z3cAozRk9KMnVUWHIrdlBGeXBVcU81MVBPUnFrUzFqN3RGb2V2a3BEYXBMKzJsdGRRNEcvQi8wQnc9PQotLS0tLUVORCBSU0EgUFJJVkFURSBLRVktLS0tLQo="


			`, // set conf URL later.
			exposedPorts: []string{"8086/tcp"},
		},
		////////////////////////////////////////////////////////////////////////
		// influxdb:1.8.10 TLS
		////////////////////////////////////////////////////////////////////////
		{
			name: "pubrepo.jiagouyun.com/image-repo-for-testing/influxdb:1.8.10-alpine-tls",
			conf: `url = ""
			interval = '2s'
			timeout = "5s"
			election = false
			[tlsconf]
			  cert = "./tls/peer.crt"
			  cert_key = "./tls/peer.key"
			  insecure_skip_verify = true
			`, // set conf URL later.
			exposedPorts: []string{"8086/tcp"},
		},
		////////////////////////////////////////////////////////////////////////
		// influxdb:1.8.10 TLS no tls file
		////////////////////////////////////////////////////////////////////////
		{
			name: "pubrepo.jiagouyun.com/image-repo-for-testing/influxdb:1.8.10-alpine-tls",
			conf: `url = ""
			interval = '2s'
			timeout = "5s"
			election = false
			[tlsconf]
			  insecure_skip_verify = true
			`, // set conf URL later.
			exposedPorts: []string{"8086/tcp"},
		},
		////////////////////////////////////////////////////////////////////////
		// influxdb:1.8.10
		////////////////////////////////////////////////////////////////////////
		{
			name: "influxdb:1.8.10-alpine",
			conf: `url = ""
				interval = '2s'
				timeout = "5s"
				election = true`, // set conf URL later.
			exposedPorts: []string{"8086/tcp"},
			optsInfluxdbCq: []inputs.PointCheckOption{
				inputs.WithExtraTags(map[string]string{
					"election": "1",
				}),
			},
			optsInfluxdbHttpd: []inputs.PointCheckOption{
				inputs.WithExtraTags(map[string]string{
					"election": "1",
				}),
			},
			optsInfluxdbMemstats: []inputs.PointCheckOption{
				inputs.WithExtraTags(map[string]string{
					"election": "1",
				}),
			},
			optsInfluxdbQueryExecutor: []inputs.PointCheckOption{
				inputs.WithExtraTags(map[string]string{
					"election": "1",
				}),
			},
			optsInfluxdbRuntime: []inputs.PointCheckOption{
				inputs.WithExtraTags(map[string]string{
					"election": "1",
				}),
			},
			optsInfluxdbSubscriber: []inputs.PointCheckOption{
				inputs.WithExtraTags(map[string]string{
					"election": "1",
				}),
			},
			optsInfluxdbWrite: []inputs.PointCheckOption{
				inputs.WithExtraTags(map[string]string{
					"election": "1",
				}),
			},
		},
		{
			name: "influxdb:1.8.10-alpine",
			conf: `url = ""
				interval = '2s'
				timeout = "5s"
				election = false`, // set conf URL later.
			exposedPorts: []string{"8086/tcp"},
		},
	}

	var cases []*caseSpec

	// compose cases
	for _, base := range bases {
		feeder := io.NewMockedFeeder()

		ipt := defaultInput()
		ipt.feeder = feeder

		_, err := toml.Decode(base.conf, ipt)
		require.NoError(t, err)

		if ipt.Election {
			ipt.Tagger = testutils.NewTaggerElection()
		} else {
			ipt.Tagger = testutils.NewTaggerHost()
		}

		repoTag := strings.Split(base.name, ":")

		cases = append(cases, &caseSpec{
			t:       t,
			ipt:     ipt,
			name:    base.name,
			feeder:  feeder,
			repo:    repoTag[0],
			repoTag: repoTag[1],

			dockerFileText: base.dockerFileText,
			exposedPorts:   base.exposedPorts,
			cmd:            base.cmd,

			optsInfluxdbCq:            base.optsInfluxdbCq,
			optsInfluxdbHttpd:         base.optsInfluxdbHttpd,
			optsInfluxdbMemstats:      base.optsInfluxdbMemstats,
			optsInfluxdbQueryExecutor: base.optsInfluxdbQueryExecutor,
			optsInfluxdbRuntime:       base.optsInfluxdbRuntime,
			optsInfluxdbSubscriber:    base.optsInfluxdbSubscriber,
			optsInfluxdbWrite:         base.optsInfluxdbWrite,

			cr: &testutils.CaseResult{
				Name:        t.Name(),
				Case:        base.name,
				ExtraFields: map[string]any{},
				ExtraTags: map[string]string{
					"image":       repoTag[0],
					"image_tag":   repoTag[1],
					"docker_host": remote.Host,
					"docker_port": remote.Port,
				},
			},
		})
	}

	return cases, nil
}

////////////////////////////////////////////////////////////////////////////////

// caseSpec.

type caseSpec struct {
	t *testing.T

	name                      string
	repo                      string
	repoTag                   string
	dockerFileText            string
	exposedPorts              []string
	serverPorts               []string
	optsInfluxdbCq            []inputs.PointCheckOption
	optsInfluxdbHttpd         []inputs.PointCheckOption
	optsInfluxdbMemstats      []inputs.PointCheckOption
	optsInfluxdbQueryExecutor []inputs.PointCheckOption
	optsInfluxdbRuntime       []inputs.PointCheckOption
	optsInfluxdbSubscriber    []inputs.PointCheckOption
	optsInfluxdbWrite         []inputs.PointCheckOption
	cmd                       []string
	mCount                    map[string]struct{}

	ipt    *Input
	feeder *io.MockedFeeder

	pool     *dockertest.Pool
	resource *dockertest.Resource

	cr *testutils.CaseResult
}

func (cs *caseSpec) checkPoint(pts []*point.Point) error {
	for _, pt := range pts {
		var opts []inputs.PointCheckOption
		opts = append(opts, inputs.WithExtraTags(cs.ipt.Tags))

		measurement := pt.Name()

		switch measurement {
		case influxdbCq:
			opts = append(opts, cs.optsInfluxdbCq...)
			opts = append(opts, inputs.WithDoc(&InfluxdbCqM{}))

			msgs := inputs.CheckPoint(pt, opts...)

			for _, msg := range msgs {
				cs.t.Logf("check measurement %s failed: %+#v", measurement, msg)
			}

			if len(msgs) > 0 {
				return fmt.Errorf("check measurement %s failed: %+#v", measurement, msgs)
			}

			cs.mCount[influxdbCq] = struct{}{}

		case influxdbHttpd:
			opts = append(opts, cs.optsInfluxdbHttpd...)
			opts = append(opts, inputs.WithDoc(&InfluxdbHttpdM{}))

			msgs := inputs.CheckPoint(pt, opts...)

			for _, msg := range msgs {
				cs.t.Logf("check measurement %s failed: %+#v", measurement, msg)
			}

			if len(msgs) > 0 {
				return fmt.Errorf("check measurement %s failed: %+#v", measurement, msgs)
			}

			cs.mCount[influxdbHttpd] = struct{}{}

		case influxdbMemstats:
			opts = append(opts, cs.optsInfluxdbMemstats...)
			opts = append(opts, inputs.WithDoc(&InfluxdbMemstatsM{}))

			msgs := inputs.CheckPoint(pt, opts...)

			for _, msg := range msgs {
				cs.t.Logf("check measurement %s failed: %+#v", measurement, msg)
			}

			if len(msgs) > 0 {
				return fmt.Errorf("check measurement %s failed: %+#v", measurement, msgs)
			}

			cs.mCount[influxdbMemstats] = struct{}{}

		case influxdbQueryExecutor:
			opts = append(opts, cs.optsInfluxdbQueryExecutor...)
			opts = append(opts, inputs.WithDoc(&InfluxdbQueryExecutorM{}))

			msgs := inputs.CheckPoint(pt, opts...)

			for _, msg := range msgs {
				cs.t.Logf("check measurement %s failed: %+#v", measurement, msg)
			}

			if len(msgs) > 0 {
				return fmt.Errorf("check measurement %s failed: %+#v", measurement, msgs)
			}

			cs.mCount[influxdbQueryExecutor] = struct{}{}

		case influxdbRuntime:
			opts = append(opts, cs.optsInfluxdbRuntime...)
			opts = append(opts, inputs.WithDoc(&InfluxdbRuntimeM{}))

			msgs := inputs.CheckPoint(pt, opts...)

			for _, msg := range msgs {
				cs.t.Logf("check measurement %s failed: %+#v", measurement, msg)
			}

			if len(msgs) > 0 {
				return fmt.Errorf("check measurement %s failed: %+#v", measurement, msgs)
			}

			cs.mCount[influxdbRuntime] = struct{}{}

		case influxdbSubscriber:
			opts = append(opts, cs.optsInfluxdbSubscriber...)
			opts = append(opts, inputs.WithDoc(&InfluxdbSubscriberM{}))

			msgs := inputs.CheckPoint(pt, opts...)

			for _, msg := range msgs {
				cs.t.Logf("check measurement %s failed: %+#v", measurement, msg)
			}

			if len(msgs) > 0 {
				return fmt.Errorf("check measurement %s failed: %+#v", measurement, msgs)
			}

			cs.mCount[influxdbSubscriber] = struct{}{}

		case influxdbWrite:
			opts = append(opts, cs.optsInfluxdbWrite...)
			opts = append(opts, inputs.WithDoc(&InfluxdbWriteM{}))

			msgs := inputs.CheckPoint(pt, opts...)

			for _, msg := range msgs {
				cs.t.Logf("check measurement %s failed: %+#v", measurement, msg)
			}

			if len(msgs) > 0 {
				return fmt.Errorf("check measurement %s failed: %+#v", measurement, msgs)
			}

			cs.mCount[influxdbWrite] = struct{}{}

		default: // TODO: check other measurement
			panic("unknown measurement: " + measurement)
		}

		// check if tag appended
		if len(cs.ipt.Tags) != 0 {
			cs.t.Logf("checking tags %+#v...", cs.ipt.Tags)

			tags := pt.Tags()
			for k, expect := range cs.ipt.Tags {
				if v := tags.Get(k); v != nil {
					got := v.GetS()
					if got != expect {
						return fmt.Errorf("expect tag value %s, got %s", expect, got)
					}
				} else {
					return fmt.Errorf("tag %s not found, got %v", k, tags)
				}
			}
		}
	}

	// TODO: some other checking on @pts, such as `if some required measurements exist'...

	return nil
}

func (cs *caseSpec) run() error {
	r := testutils.GetRemote()
	dockerTCP := r.TCPURL()

	cs.t.Logf("get remote: %+#v, TCP: %s", r, dockerTCP)

	start := time.Now()

	p, err := cs.getPool(dockerTCP)
	if err != nil {
		return err
	}

	dockerFileDir, dockerFilePath, err := cs.getDockerFilePath()
	if err != nil {
		return err
	}
	defer os.RemoveAll(dockerFileDir)

	uniqueContainerName := testutils.GetUniqueContainerName(inputName)

	var resource *dockertest.Resource

	if len(cs.dockerFileText) == 0 {
		// Just run a container from existing docker image.
		resource, err = p.RunWithOptions(
			&dockertest.RunOptions{
				Name: uniqueContainerName, // ATTENTION: not cs.name.

				Repository: cs.repo,
				Tag:        cs.repoTag,
				Cmd:        cs.cmd,

				ExposedPorts: cs.exposedPorts,
			},

			func(c *docker.HostConfig) {
				c.RestartPolicy = docker.RestartPolicy{Name: "no"}
				c.AutoRemove = true
			},
		)
	} else {
		// Build docker image from Dockerfile and run a container from it.
		resource, err = p.BuildAndRunWithOptions(
			dockerFilePath,

			&dockertest.RunOptions{
				ContainerName: uniqueContainerName,
				Name:          cs.name, // ATTENTION: not uniqueContainerName.

				Repository: cs.repo,
				Tag:        cs.repoTag,
				Cmd:        cs.cmd,

				ExposedPorts: cs.exposedPorts,
			},

			func(c *docker.HostConfig) {
				c.RestartPolicy = docker.RestartPolicy{Name: "no"}
				c.AutoRemove = true
			},
		)
	}

	if err != nil {
		return err
	}

	cs.pool = p
	cs.resource = resource

	if err := cs.getMappingPorts(); err != nil {
		return err
	}
	if cs.ipt.TLSConf != nil {
		cs.ipt.URL = getConfAccessPoint("https", r.Host, cs.serverPorts[0]) // set conf URL here.
	} else {
		cs.ipt.URL = getConfAccessPoint("http", r.Host, cs.serverPorts[0]) // set conf URL here.
	}

	cs.t.Logf("check service(%s:%v)...", r.Host, cs.serverPorts)

	if err := cs.portsOK(r); err != nil {
		return err
	}

	cs.cr.AddField("container_ready_cost", int64(time.Since(start)))

	var wg sync.WaitGroup

	// start input
	cs.t.Logf("start input...")
	wg.Add(1)
	go func() {
		defer wg.Done()
		cs.ipt.Run()
	}()

	// wait data
	start = time.Now()
	cs.t.Logf("wait points...")
	pts, err := cs.feeder.NPoints(10, 5*time.Minute)
	if err != nil {
		return err
	}

	cs.cr.AddField("point_latency", int64(time.Since(start)))
	cs.cr.AddField("point_count", len(pts))

	// for k, v := range pts {
	// 	cs.t.Logf("pt [%02d] = %s", k, v.LineProto())
	// }

	cs.t.Logf("get %d points", len(pts))
	cs.mCount = make(map[string]struct{})
	if err := cs.checkPoint(pts); err != nil {
		return err
	}

	cs.t.Logf("stop input...")
	cs.ipt.Terminate()

	require.Equal(cs.t, mExpect, cs.mCount)

	cs.t.Logf("exit...")
	wg.Wait()

	return nil
}

func (cs *caseSpec) getPool(endpoint string) (*dockertest.Pool, error) {
	p, err := dockertest.NewPool(endpoint)
	if err != nil {
		return nil, err
	}
	err = p.Client.Ping()
	if err != nil {
		cs.t.Logf("Could not connect to Docker: %v", err)
		return nil, err
	}
	return p, nil
}

func (cs *caseSpec) getDockerFilePath() (dirName string, fileName string, err error) {
	if len(cs.dockerFileText) == 0 {
		return
	}

	tmpDir, err := ioutil.TempDir("", "dockerfiles_")
	if err != nil {
		cs.t.Logf("ioutil.TempDir failed: %s", err.Error())
		return "", "", err
	}

	tmpFile, err := ioutil.TempFile(tmpDir, "dockerfile_")
	if err != nil {
		cs.t.Logf("ioutil.TempFile failed: %s", err.Error())
		return "", "", err
	}

	_, err = tmpFile.WriteString(cs.dockerFileText)
	if err != nil {
		cs.t.Logf("TempFile.WriteString failed: %s", err.Error())
		return "", "", err
	}

	if err := os.Chmod(tmpFile.Name(), os.ModePerm); err != nil {
		cs.t.Logf("os.Chmod failed: %s", err.Error())
		return "", "", err
	}

	if err := tmpFile.Close(); err != nil {
		cs.t.Logf("Close failed: %s", err.Error())
		return "", "", err
	}

	return tmpDir, tmpFile.Name(), nil
}

func (cs *caseSpec) getMappingPorts() error {
	cs.serverPorts = make([]string, len(cs.exposedPorts))
	for k, v := range cs.exposedPorts {
		mapStr := cs.resource.GetHostPort(v)
		_, port, err := net.SplitHostPort(mapStr)
		if err != nil {
			return err
		}
		cs.serverPorts[k] = port
	}
	return nil
}

func (cs *caseSpec) portsOK(r *testutils.RemoteInfo) error {
	for _, v := range cs.serverPorts {
		if !r.PortOK(docker.Port(v).Port(), time.Minute) {
			return fmt.Errorf("service checking failed")
		}
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////

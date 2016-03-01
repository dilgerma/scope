package render_test

import (
	"net"
	"reflect"
	"testing"

	"github.com/dilgerma/scope/probe/host"
	"github.com/dilgerma/scope/render"
	"github.com/dilgerma/scope/report"
	"github.com/dilgerma/scope/test"
)

func TestReportLocalNetworks(t *testing.T) {
	r := report.MakeReport().Merge(report.Report{
		Host: report.Topology{
			Nodes: report.Nodes{
				"nonets": report.MakeNode(),
				"foo": report.MakeNode().WithSets(report.EmptySets.
					Add(host.LocalNetworks, report.MakeStringSet(
					"10.0.0.1/8", "192.168.1.1/24", "10.0.0.1/8", "badnet/33")),
				),
			},
		},
	})
	want := report.Networks([]*net.IPNet{
		mustParseCIDR("10.0.0.1/8"),
		mustParseCIDR("192.168.1.1/24"),
	})
	have := render.LocalNetworks(r)
	if !reflect.DeepEqual(want, have) {
		t.Errorf("%s", test.Diff(want, have))
	}
}

func mustParseCIDR(s string) *net.IPNet {
	_, ipNet, err := net.ParseCIDR(s)
	if err != nil {
		panic(err)
	}
	return ipNet
}

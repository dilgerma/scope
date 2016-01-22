package host_test

import (
	"reflect"
	"testing"

	"github.com/dilgerma/scope/probe/host"
	"github.com/dilgerma/scope/report"
	"github.com/dilgerma/scope/test"
)

func TestTagger(t *testing.T) {
	var (
		hostID         = "foo"
		probeID        = "a1b2c3d4"
		endpointNodeID = report.MakeEndpointNodeID(hostID, "1.2.3.4", "56789") // hostID ignored
		nodeMetadata   = report.MakeNodeWith(map[string]string{"foo": "bar"})
	)

	r := report.MakeReport()
	r.Process.AddNode(endpointNodeID, nodeMetadata)
	want := nodeMetadata.Merge(report.MakeNodeWith(map[string]string{
		report.HostNodeID: report.MakeHostNodeID(hostID),
		report.ProbeID:    probeID,
	}))
	rpt, _ := host.NewTagger(hostID, probeID).Tag(r)
	have := rpt.Process.Nodes[endpointNodeID].Copy()
	if !reflect.DeepEqual(want, have) {
		t.Error(test.Diff(want, have))
	}
}

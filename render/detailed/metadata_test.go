package detailed_test

import (
	"reflect"
	"testing"

	"github.com/dilgerma/scope/probe/docker"
	"github.com/dilgerma/scope/render/detailed"
	"github.com/dilgerma/scope/report"
	"github.com/dilgerma/scope/test"
	"github.com/dilgerma/scope/test/fixture"
)

func TestNodeMetadata(t *testing.T) {
	inputs := []struct {
		name string
		node report.Node
		want []detailed.MetadataRow
	}{
		{
			name: "container",
			node: report.MakeNodeWith(map[string]string{
				docker.ContainerID:            fixture.ClientContainerID,
				docker.LabelPrefix + "label1": "label1value",
				docker.ContainerState:         docker.StateRunning,
			}).WithTopology(report.Container).WithSets(report.EmptySets.
				Add(docker.ContainerIPs, report.MakeStringSet("10.10.10.0/24", "10.10.10.1/24")),
			),
			want: []detailed.MetadataRow{
				{ID: docker.ContainerID, Value: fixture.ClientContainerID, Prime: true},
				{ID: docker.ContainerState, Value: "running", Prime: true},
				{ID: docker.ContainerIPs, Value: "10.10.10.0/24, 10.10.10.1/24"},
			},
		},
		{
			name: "unknown topology",
			node: report.MakeNodeWith(map[string]string{
				docker.ContainerID: fixture.ClientContainerID,
			}).WithTopology("foobar").WithID(fixture.ClientContainerNodeID),
			want: nil,
		},
	}
	for _, input := range inputs {
		have := detailed.NodeMetadata(input.node)
		if !reflect.DeepEqual(input.want, have) {
			t.Errorf("%s: %s", input.name, test.Diff(input.want, have))
		}
	}
}

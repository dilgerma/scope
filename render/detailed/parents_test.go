package detailed_test

import (
	"fmt"
	"testing"

	"github.com/dilgerma/scope/render"
	"github.com/dilgerma/scope/render/detailed"
	"github.com/dilgerma/scope/report"
	"github.com/dilgerma/scope/test"
	"github.com/dilgerma/scope/test/fixture"
	"github.com/dilgerma/scope/test/reflect"
)

func TestParents(t *testing.T) {
	for _, c := range []struct {
		name string
		node report.Node
		want []detailed.Parent
	}{
		{
			name: "Node accidentally tagged with itself",
			node: render.HostRenderer.Render(fixture.Report, render.FilterNoop)[fixture.ClientHostNodeID].WithParents(
				report.EmptySets.Add(report.Host, report.MakeStringSet(fixture.ClientHostNodeID)),
			),
			want: nil,
		},
		{
			node: render.HostRenderer.Render(fixture.Report, render.FilterNoop)[fixture.ClientHostNodeID],
			want: nil,
		},
		{
			node: render.ContainerImageRenderer.Render(fixture.Report, render.FilterNoop)[fixture.ClientContainerImageNodeID],
			want: []detailed.Parent{
				{ID: fixture.ClientHostNodeID, Label: fixture.ClientHostName, TopologyID: "hosts"},
			},
		},
		{
			node: render.ContainerRenderer.Render(fixture.Report, render.FilterNoop)[fixture.ClientContainerNodeID],
			want: []detailed.Parent{
				{ID: fixture.ClientContainerImageNodeID, Label: fixture.ClientContainerImageName, TopologyID: "containers-by-image"},
				{ID: fixture.ClientHostNodeID, Label: fixture.ClientHostName, TopologyID: "hosts"},
				{ID: fixture.ClientPodNodeID, Label: "pong-a", TopologyID: "pods"},
			},
		},
		{
			node: render.ProcessRenderer.Render(fixture.Report, render.FilterNoop)[fixture.ClientProcess1NodeID],
			want: []detailed.Parent{
				{ID: fixture.ClientContainerNodeID, Label: fixture.ClientContainerName, TopologyID: "containers"},
				{ID: fixture.ClientContainerImageNodeID, Label: fixture.ClientContainerImageName, TopologyID: "containers-by-image"},
				{ID: fixture.ClientHostNodeID, Label: fixture.ClientHostName, TopologyID: "hosts"},
			},
		},
	} {
		name := c.name
		if name == "" {
			name = fmt.Sprintf("Node %q", c.node.ID)
		}
		if have := detailed.Parents(fixture.Report, c.node); !reflect.DeepEqual(c.want, have) {
			t.Errorf("%s: %s", name, test.Diff(c.want, have))
		}
	}
}

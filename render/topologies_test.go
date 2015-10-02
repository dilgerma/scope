package render_test

import (
	"reflect"
	"testing"

	"github.com/weaveworks/scope/probe/docker"
	"github.com/weaveworks/scope/render"
	"github.com/weaveworks/scope/render/expected"
	"github.com/weaveworks/scope/test"
	"github.com/weaveworks/scope/test/fixture"
)

func TestProcessRenderer(t *testing.T) {
	have := render.ProcessRenderer.Render(fixture.Report).Prune()
	want := expected.RenderedProcesses
	if !reflect.DeepEqual(want, have) {
		t.Error(test.Diff(want, have))
	}
}

func TestProcessNameRenderer(t *testing.T) {
	have := render.ProcessNameRenderer.Render(fixture.Report).Prune()
	want := expected.RenderedProcessNames
	if !reflect.DeepEqual(want, have) {
		t.Error(test.Diff(want, have))
	}
}

func TestContainerRenderer(t *testing.T) {
	have := (render.ContainerWithImageNameRenderer.Render(fixture.Report)).Prune()
	want := expected.RenderedContainers
	if !reflect.DeepEqual(want, have) {
		t.Error(test.Diff(want, have))
	}
}

func TestContainerFilterRenderer(t *testing.T) {
	// tag on of the containers in the topology and ensure
	// it is filtered out correctly.
	input := fixture.Report.Copy()
	input.Container.Nodes[fixture.ClientContainerNodeID].Metadata[docker.LabelPrefix+"works.weave.role"] = "system"
	have := render.FilterSystem(render.ContainerWithImageNameRenderer).Render(input).Prune()
	want := expected.RenderedContainers.Copy()
	delete(want, fixture.ClientContainerID)
	if !reflect.DeepEqual(want, have) {
		t.Error(test.Diff(want, have))
	}
}

func TestContainerImageRenderer(t *testing.T) {
	have := render.ContainerImageRenderer.Render(fixture.Report).Prune()
	want := expected.RenderedContainerImages
	if !reflect.DeepEqual(want, have) {
		t.Error(test.Diff(want, have))
	}
}

func TestHostRenderer(t *testing.T) {
	have := render.HostRenderer.Render(fixture.Report).Prune()
	want := expected.RenderedHosts
	if !reflect.DeepEqual(want, have) {
		t.Error(test.Diff(want, have))
	}
}

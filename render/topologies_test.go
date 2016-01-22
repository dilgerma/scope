package render_test

import (
	"reflect"
	"testing"

	"github.com/dilgerma/scope/probe/docker"
	"github.com/dilgerma/scope/probe/kubernetes"
	"github.com/dilgerma/scope/render"
	"github.com/dilgerma/scope/render/expected"
	"github.com/dilgerma/scope/test"
	"github.com/dilgerma/scope/test/fixture"
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

func TestPodRenderer(t *testing.T) {
	have := render.PodRenderer.Render(fixture.Report).Prune()
	want := expected.RenderedPods
	if !reflect.DeepEqual(want, have) {
		t.Error(test.Diff(want, have))
	}
}

func TestPodFilterRenderer(t *testing.T) {
	// tag on containers or pod namespace in the topology and ensure
	// it is filtered out correctly.
	input := fixture.Report.Copy()
	input.Pod.Nodes[fixture.ClientPodNodeID].Metadata[kubernetes.PodID] = "kube-system/foo"
	input.Pod.Nodes[fixture.ClientPodNodeID].Metadata[kubernetes.Namespace] = "kube-system"
	input.Pod.Nodes[fixture.ClientPodNodeID].Metadata[kubernetes.PodName] = "foo"
	input.Container.Nodes[fixture.ClientContainerNodeID].Metadata[docker.LabelPrefix+"io.kubernetes.pod.name"] = "kube-system/foo"
	have := render.FilterSystem(render.PodRenderer).Render(input).Prune()
	want := expected.RenderedPods.Copy()
	delete(want, fixture.ClientPodID)
	delete(want, fixture.ClientContainerID)
	if !reflect.DeepEqual(want, have) {
		t.Error(test.Diff(want, have))
	}
}

func TestPodServiceRenderer(t *testing.T) {
	have := render.PodServiceRenderer.Render(fixture.Report).Prune()
	want := expected.RenderedPodServices
	if !reflect.DeepEqual(want, have) {
		t.Error(test.Diff(want, have))
	}
}

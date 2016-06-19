package render_test

import (
	"testing"

	"github.com/dilgerma/scope/render"
	"github.com/dilgerma/scope/render/expected"
	"github.com/dilgerma/scope/test"
	"github.com/dilgerma/scope/test/fixture"
	"github.com/dilgerma/scope/test/reflect"
)

func TestEndpointRenderer(t *testing.T) {
	have := Prune(render.EndpointRenderer.Render(fixture.Report, FilterNoop))
	want := Prune(expected.RenderedEndpoints)
	if !reflect.DeepEqual(want, have) {
		t.Error(test.Diff(want, have))
	}
}

func TestProcessRenderer(t *testing.T) {
	have := Prune(render.ProcessRenderer.Render(fixture.Report, FilterNoop))
	want := Prune(expected.RenderedProcesses)
	if !reflect.DeepEqual(want, have) {
		t.Error(test.Diff(want, have))
	}
}

func TestProcessNameRenderer(t *testing.T) {
	have := Prune(render.ProcessNameRenderer.Render(fixture.Report, FilterNoop))
	want := Prune(expected.RenderedProcessNames)
	if !reflect.DeepEqual(want, have) {
		t.Error(test.Diff(want, have))
	}
}

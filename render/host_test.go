package render_test

import (
	"testing"

	"github.com/dilgerma/scope/render"
	"github.com/dilgerma/scope/render/expected"
	"github.com/dilgerma/scope/test"
	"github.com/dilgerma/scope/test/fixture"
	"github.com/dilgerma/scope/test/reflect"
)

func TestHostRenderer(t *testing.T) {
	have := Prune(render.HostRenderer.Render(fixture.Report, render.FilterNoop))
	want := Prune(expected.RenderedHosts)
	if !reflect.DeepEqual(want, have) {
		t.Error(test.Diff(want, have))
	}
}

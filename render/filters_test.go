package render_test

import (
	"testing"

	"github.com/dilgerma/scope/render"
	"github.com/dilgerma/scope/report"
	"github.com/dilgerma/scope/test"
	"github.com/dilgerma/scope/test/reflect"
)

func TestFilterRender(t *testing.T) {
	renderer := render.FilterUnconnected(
		mockRenderer{RenderableNodes: render.RenderableNodes{
			"foo": {ID: "foo", Node: report.MakeNode().WithAdjacent("bar")},
			"bar": {ID: "bar", Node: report.MakeNode().WithAdjacent("foo")},
			"baz": {ID: "baz", Node: report.MakeNode()},
		}})
	want := render.RenderableNodes{
		"foo": {ID: "foo", Node: report.MakeNode().WithAdjacent("bar")},
		"bar": {ID: "bar", Node: report.MakeNode().WithAdjacent("foo")},
	}
	have := renderer.Render(report.MakeReport()).Prune()
	if !reflect.DeepEqual(want, have) {
		t.Error(test.Diff(want, have))
	}
}

func TestFilterRender2(t *testing.T) {
	// Test adjacencies are removed for filtered nodes.
	renderer := render.Filter{
		FilterFunc: func(node render.RenderableNode) bool {
			return node.ID != "bar"
		},
		Renderer: mockRenderer{RenderableNodes: render.RenderableNodes{
			"foo": {ID: "foo", Node: report.MakeNode().WithAdjacent("bar")},
			"bar": {ID: "bar", Node: report.MakeNode().WithAdjacent("foo")},
			"baz": {ID: "baz", Node: report.MakeNode()},
		}},
	}
	want := render.RenderableNodes{
		"foo": {ID: "foo", Node: report.MakeNode()},
		"baz": {ID: "baz", Node: report.MakeNode()},
	}
	have := renderer.Render(report.MakeReport()).Prune()
	if !reflect.DeepEqual(want, have) {
		t.Error(test.Diff(want, have))
	}
}

func TestFilterUnconnectedPseudoNodes(t *testing.T) {
	// Test pseudo nodes that are made unconnected by filtering
	// are also removed.
	{
		nodes := render.RenderableNodes{
			"foo": {ID: "foo", Node: report.MakeNode().WithAdjacent("bar")},
			"bar": {ID: "bar", Node: report.MakeNode().WithAdjacent("baz")},
			"baz": {ID: "baz", Node: report.MakeNode(), Pseudo: true},
		}
		renderer := render.Filter{
			FilterFunc: func(node render.RenderableNode) bool {
				return true
			},
			Renderer: mockRenderer{RenderableNodes: nodes},
		}
		want := nodes.Prune()
		have := renderer.Render(report.MakeReport()).Prune()
		if !reflect.DeepEqual(want, have) {
			t.Error(test.Diff(want, have))
		}
	}
	{
		renderer := render.Filter{
			FilterFunc: func(node render.RenderableNode) bool {
				return node.ID != "bar"
			},
			Renderer: mockRenderer{RenderableNodes: render.RenderableNodes{
				"foo": {ID: "foo", Node: report.MakeNode().WithAdjacent("bar")},
				"bar": {ID: "bar", Node: report.MakeNode().WithAdjacent("baz")},
				"baz": {ID: "baz", Node: report.MakeNode(), Pseudo: true},
			}},
		}
		want := render.RenderableNodes{
			"foo": {ID: "foo", Node: report.MakeNode()},
		}
		have := renderer.Render(report.MakeReport()).Prune()
		if !reflect.DeepEqual(want, have) {
			t.Error(test.Diff(want, have))
		}
	}
	{
		renderer := render.Filter{
			FilterFunc: func(node render.RenderableNode) bool {
				return node.ID != "bar"
			},
			Renderer: mockRenderer{RenderableNodes: render.RenderableNodes{
				"foo": {ID: "foo", Node: report.MakeNode()},
				"bar": {ID: "bar", Node: report.MakeNode().WithAdjacent("foo")},
				"baz": {ID: "baz", Node: report.MakeNode().WithAdjacent("bar"), Pseudo: true},
			}},
		}
		want := render.RenderableNodes{
			"foo": {ID: "foo", Node: report.MakeNode()},
		}
		have := renderer.Render(report.MakeReport()).Prune()
		if !reflect.DeepEqual(want, have) {
			t.Error(test.Diff(want, have))
		}
	}
}

func TestFilterUnconnectedSelf(t *testing.T) {
	// Test nodes that are only connected to themselves are filtered.
	{
		nodes := render.RenderableNodes{
			"foo": {ID: "foo", Node: report.MakeNode().WithAdjacent("foo")},
		}
		renderer := render.FilterUnconnected(mockRenderer{RenderableNodes: nodes})
		want := render.RenderableNodes{}
		have := renderer.Render(report.MakeReport()).Prune()
		if !reflect.DeepEqual(want, have) {
			t.Error(test.Diff(want, have))
		}
	}
}

func TestFilterPseudo(t *testing.T) {
	// Test pseudonodes are removed
	{
		nodes := render.RenderableNodes{
			"foo": {ID: "foo", Node: report.MakeNode()},
			"bar": {ID: "bar", Pseudo: true, Node: report.MakeNode()},
		}
		renderer := render.FilterPseudo(mockRenderer{RenderableNodes: nodes})
		want := render.RenderableNodes{
			"foo": {ID: "foo", Node: report.MakeNode()},
		}
		have := renderer.Render(report.MakeReport()).Prune()
		if !reflect.DeepEqual(want, have) {
			t.Error(test.Diff(want, have))
		}
	}
}

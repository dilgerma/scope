package docker_test

import (
	"reflect"
	"testing"

	"github.com/dilgerma/scope/probe/docker"
	"github.com/dilgerma/scope/report"
	"github.com/dilgerma/scope/test"
)

func TestLabels(t *testing.T) {
	want := map[string]string{
		"foo1": "bar1",
		"foo2": "bar2",
	}
	nmd := report.MakeNode()

	nmd = docker.AddLabels(nmd, want)
	have := docker.ExtractLabels(nmd)

	if !reflect.DeepEqual(want, have) {
		t.Error(test.Diff(want, have))
	}
}

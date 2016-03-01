package docker_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/dilgerma/scope/common/mtime"
	"github.com/dilgerma/scope/probe/docker"
	"github.com/dilgerma/scope/probe/process"
	"github.com/dilgerma/scope/report"
)

type mockProcessTree struct {
	parents map[int]int
}

func (m *mockProcessTree) GetParent(pid int) (int, error) {
	parent, ok := m.parents[pid]
	if !ok {
		return -1, fmt.Errorf("Not found %d", pid)
	}
	return parent, nil
}

func (m *mockProcessTree) GetChildren(int) ([]int, error) {
	panic("Not implemented")
}

func TestTagger(t *testing.T) {
	mtime.NowForce(time.Now())
	defer mtime.NowReset()

	oldProcessTree := docker.NewProcessTreeStub
	defer func() { docker.NewProcessTreeStub = oldProcessTree }()

	docker.NewProcessTreeStub = func(_ process.Walker) (process.Tree, error) {
		return &mockProcessTree{map[int]int{3: 2}}, nil
	}

	var (
		pid1NodeID = report.MakeProcessNodeID("somehost.com", "2")
		pid2NodeID = report.MakeProcessNodeID("somehost.com", "3")
	)

	input := report.MakeReport()
	input.Process.AddNode(pid1NodeID, report.MakeNodeWith(map[string]string{process.PID: "2"}))
	input.Process.AddNode(pid2NodeID, report.MakeNodeWith(map[string]string{process.PID: "3"}))

	have, err := docker.NewTagger(mockRegistryInstance, nil).Tag(input)
	if err != nil {
		t.Errorf("%v", err)
	}

	// Processes should be tagged with their container ID, and parents
	for _, nodeID := range []string{pid1NodeID, pid2NodeID} {
		node, ok := have.Process.Nodes[nodeID]
		if !ok {
			t.Errorf("Expected process node %s, but not found", nodeID)
		}

		// node should have the container id added
		if have, ok := node.Latest.Lookup(docker.ContainerID); !ok || have != "ping" {
			t.Errorf("Expected process node %s to have container id %q, got %q", nodeID, "ping", have)
		}

		// node should have the container as a parent
		if have, ok := node.Parents.Lookup(report.Container); !ok || !have.Contains(report.MakeContainerNodeID("ping")) {
			t.Errorf("Expected process node %s to have container %q as a parent, got %q", nodeID, "ping", have)
		}

		// node should have the container image as a parent
		if have, ok := node.Parents.Lookup(report.ContainerImage); !ok || !have.Contains(report.MakeContainerImageNodeID("baz")) {
			t.Errorf("Expected process node %s to have container image %q as a parent, got %q", nodeID, "baz", have)
		}
	}
}

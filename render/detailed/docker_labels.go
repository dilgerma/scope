package detailed

import (
	"sort"

	"github.com/dilgerma/scope/probe/docker"
	"github.com/dilgerma/scope/report"
)

// NodeDockerLabels produces a table (to be consumed directly by the UI) based
// on an origin ID, which is (optimistically) a node ID in one of our
// topologies.
func NodeDockerLabels(nmd report.Node) []MetadataRow {
	if nmd.Topology != report.Container && nmd.Topology != report.ContainerImage {
		return nil
	}

	var rows []MetadataRow
	// Add labels in alphabetical order
	labels := docker.ExtractLabels(nmd)
	labelKeys := make([]string, 0, len(labels))
	for k := range labels {
		labelKeys = append(labelKeys, k)
	}
	sort.Strings(labelKeys)
	for _, labelKey := range labelKeys {
		rows = append(rows, MetadataRow{ID: "label_" + labelKey, Value: labels[labelKey]})
	}
	return rows
}

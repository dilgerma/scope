package render

import (
	"github.com/dilgerma/scope/report"
)

// TopologySelector selects a single topology from a report.
// NB it is also a Renderer!
type TopologySelector func(r report.Report) RenderableNodes

// Render implements Renderer
func (t TopologySelector) Render(r report.Report) RenderableNodes {
	return t(r)
}

// Stats implements Renderer
func (t TopologySelector) Stats(r report.Report) Stats {
	return Stats{}
}

// MakeRenderableNodes converts a topology to a set of RenderableNodes
func MakeRenderableNodes(t report.Topology) RenderableNodes {
	result := RenderableNodes{}
	for id, nmd := range t.Nodes {
		result[id] = NewRenderableNode(id).WithNode(nmd)
	}

	// Push EdgeMetadata to both ends of the edges
	for srcID, srcNode := range result {
		srcNode.Edges.ForEach(func(dstID string, emd report.EdgeMetadata) {
			srcNode.EdgeMetadata = srcNode.EdgeMetadata.Flatten(emd)

			dstNode := result[dstID]
			dstNode.EdgeMetadata = dstNode.EdgeMetadata.Flatten(emd.Reversed())
			result[dstID] = dstNode
		})
		result[srcID] = srcNode
	}
	return result
}

var (
	// SelectEndpoint selects the endpoint topology.
	SelectEndpoint = TopologySelector(func(r report.Report) RenderableNodes {
		return MakeRenderableNodes(r.Endpoint)
	})

	// SelectProcess selects the process topology.
	SelectProcess = TopologySelector(func(r report.Report) RenderableNodes {
		return MakeRenderableNodes(r.Process)
	})

	// SelectContainer selects the container topology.
	SelectContainer = TopologySelector(func(r report.Report) RenderableNodes {
		return MakeRenderableNodes(r.Container)
	})

	// SelectContainerImage selects the container image topology.
	SelectContainerImage = TopologySelector(func(r report.Report) RenderableNodes {
		return MakeRenderableNodes(r.ContainerImage)
	})

	// SelectAddress selects the address topology.
	SelectAddress = TopologySelector(func(r report.Report) RenderableNodes {
		return MakeRenderableNodes(r.Address)
	})

	// SelectHost selects the address topology.
	SelectHost = TopologySelector(func(r report.Report) RenderableNodes {
		return MakeRenderableNodes(r.Host)
	})

	// SelectPod selects the pod topology.
	SelectPod = TopologySelector(func(r report.Report) RenderableNodes {
		return MakeRenderableNodes(r.Pod)
	})

	// SelectService selects the service topology.
	SelectService = TopologySelector(func(r report.Report) RenderableNodes {
		return MakeRenderableNodes(r.Service)
	})
)

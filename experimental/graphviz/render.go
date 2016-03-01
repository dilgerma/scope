package main

import (
	"fmt"

	"github.com/dilgerma/scope/render"
	"github.com/dilgerma/scope/report"
)

func renderTo(rpt report.Report, topology string) (render.RenderableNodes, error) {
	renderer, ok := map[string]render.Renderer{
		"processes":           render.FilterUnconnected(render.ProcessWithContainerNameRenderer),
		"processes-by-name":   render.FilterUnconnected(render.ProcessNameRenderer),
		"containers":          render.ContainerWithImageNameRenderer,
		"containers-by-image": render.ContainerImageRenderer,
		"hosts":               render.HostRenderer,
	}[topology]
	if !ok {
		return render.RenderableNodes{}, fmt.Errorf("unknown topology %v", topology)
	}
	return renderer.Render(rpt), nil
}

package process

import (
	"strconv"

	"github.com/dilgerma/scope/common/mtime"
	"github.com/dilgerma/scope/report"
)

// We use these keys in node metadata
const (
	PID            = "pid"
	Name           = "name"
	PPID           = "ppid"
	Cmdline        = "cmdline"
	Threads        = "threads"
	CPUUsage       = "process_cpu_usage_percent"
	MemoryUsage    = "process_memory_usage_bytes"
	OpenFilesCount = "open_files_count"
)

// Reporter generates Reports containing the Process topology.
type Reporter struct {
	scope   string
	walker  Walker
	jiffies Jiffies
}

// Jiffies is the type for the function used to fetch the elapsed jiffies.
type Jiffies func() (uint64, float64, error)

// NewReporter makes a new Reporter.
func NewReporter(walker Walker, scope string, jiffies Jiffies) *Reporter {
	return &Reporter{
		scope:   scope,
		walker:  walker,
		jiffies: jiffies,
	}
}

// Name of this reporter, for metrics gathering
func (Reporter) Name() string { return "Process" }

// Report implements Reporter.
func (r *Reporter) Report() (report.Report, error) {
	result := report.MakeReport()
	processes, err := r.processTopology()
	if err != nil {
		return result, err
	}
	result.Process = result.Process.Merge(processes)
	return result, nil
}

func (r *Reporter) processTopology() (report.Topology, error) {
	t := report.MakeTopology()
	now := mtime.Now()
	deltaTotal, maxCPU, err := r.jiffies()
	if err != nil {
		return t, err
	}

	err = r.walker.Walk(func(p, prev Process) {
		pidstr := strconv.Itoa(p.PID)
		nodeID := report.MakeProcessNodeID(r.scope, pidstr)
		node := report.MakeNode()
		for _, tuple := range []struct{ key, value string }{
			{PID, pidstr},
			{Name, p.Name},
			{Cmdline, p.Cmdline},
			{Threads, strconv.Itoa(p.Threads)},
		} {
			if tuple.value != "" {
				node = node.WithLatests(map[string]string{tuple.key: tuple.value})
			}
		}

		if p.PPID > 0 {
			node = node.WithLatests(map[string]string{PPID: strconv.Itoa(p.PPID)})
		}

		if deltaTotal > 0 {
			cpuUsage := float64(p.Jiffies-prev.Jiffies) / float64(deltaTotal) * 100.
			node = node.WithMetric(CPUUsage, report.MakeMetric().Add(now, cpuUsage).WithMax(maxCPU))
		}

		node = node.WithMetric(MemoryUsage, report.MakeMetric().Add(now, float64(p.RSSBytes)))
		node = node.WithMetric(OpenFilesCount, report.MakeMetric().Add(now, float64(p.OpenFilesCount)))

		t.AddNode(nodeID, node)
	})

	return t, err
}

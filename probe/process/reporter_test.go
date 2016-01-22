package process_test

import (
	"reflect"
	"testing"

	"github.com/dilgerma/scope/common/mtime"
	"github.com/dilgerma/scope/probe/process"
	"github.com/dilgerma/scope/report"
	"github.com/dilgerma/scope/test"
)

type mockWalker struct {
	processes []process.Process
}

func (m *mockWalker) Walk(f func(process.Process)) error {
	for _, p := range m.processes {
		f(p)
	}
	return nil
}

func TestReporter(t *testing.T) {
	walker := &mockWalker{
		processes: []process.Process{
			{PID: 1, PPID: 0, Comm: "init"},
			{PID: 2, PPID: 1, Comm: "bash"},
			{PID: 3, PPID: 1, Comm: "apache", Threads: 2},
			{PID: 4, PPID: 2, Comm: "ping", Cmdline: "ping foo.bar.local"},
			{PID: 5, PPID: 1, Cmdline: "tail -f /var/log/syslog"},
		},
	}

	reporter := process.NewReporter(walker, "")
	want := report.MakeReport()
	want.Process = report.MakeTopology().AddNode(
		report.MakeProcessNodeID("", "1"), report.MakeNodeWith(map[string]string{
			process.PID:     "1",
			process.Comm:    "init",
			process.Threads: "0",
		}),
	).AddNode(
		report.MakeProcessNodeID("", "2"), report.MakeNodeWith(map[string]string{
			process.PID:     "2",
			process.Comm:    "bash",
			process.PPID:    "1",
			process.Threads: "0",
		}),
	).AddNode(
		report.MakeProcessNodeID("", "3"), report.MakeNodeWith(map[string]string{
			process.PID:     "3",
			process.Comm:    "apache",
			process.PPID:    "1",
			process.Threads: "2",
		}),
	).AddNode(
		report.MakeProcessNodeID("", "4"), report.MakeNodeWith(map[string]string{
			process.PID:     "4",
			process.Comm:    "ping",
			process.PPID:    "2",
			process.Cmdline: "ping foo.bar.local",
			process.Threads: "0",
		}),
	).AddNode(
		report.MakeProcessNodeID("", "5"), report.MakeNodeWith(map[string]string{
			process.PID:     "5",
			process.PPID:    "1",
			process.Cmdline: "tail -f /var/log/syslog",
			process.Threads: "0",
		}),
	)

	have, err := reporter.Report()
	if err != nil || !reflect.DeepEqual(want, have) {
		t.Errorf("%s (%v)", test.Diff(want, have), err)
	}
}

package overlay_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/dilgerma/scope/common/exec"
	"github.com/dilgerma/scope/probe/docker"
	"github.com/dilgerma/scope/probe/overlay"
	"github.com/dilgerma/scope/report"
	"github.com/dilgerma/scope/test"
	testExec "github.com/dilgerma/scope/test/exec"
)

func TestWeaveTaggerOverlayTopology(t *testing.T) {
	oldExecCmd := exec.Command
	defer func() { exec.Command = oldExecCmd }()
	exec.Command = func(name string, args ...string) exec.Cmd {
		return testExec.NewMockCmdString(fmt.Sprintf("%s %s %s/24\n", mockContainerID, mockContainerMAC, mockContainerIP))
	}

	s := httptest.NewServer(http.HandlerFunc(mockWeaveRouter))
	defer s.Close()

	w := overlay.NewWeave(mockHostID, s.URL)
	defer w.Stop()
	w.Tick()

	{
		have, err := w.Report()
		if err != nil {
			t.Fatal(err)
		}
		if want, have := report.MakeTopology().AddNode(
			report.MakeOverlayNodeID(mockWeavePeerName),
			report.MakeNodeWith(map[string]string{
				overlay.WeavePeerName:     mockWeavePeerName,
				overlay.WeavePeerNickName: mockWeavePeerNickName,
			}),
		), have.Overlay; !reflect.DeepEqual(want, have) {
			t.Error(test.Diff(want, have))
		}
	}

	{
		nodeID := report.MakeContainerNodeID(mockHostID, mockContainerID)
		want := report.Report{
			Container: report.MakeTopology().AddNode(nodeID, report.MakeNodeWith(map[string]string{
				docker.ContainerID:       mockContainerID,
				overlay.WeaveDNSHostname: mockHostname,
				overlay.WeaveMACAddress:  mockContainerMAC,
			}).WithSets(report.Sets{
				docker.ContainerIPs:           report.MakeStringSet(mockContainerIP),
				docker.ContainerIPsWithScopes: report.MakeStringSet(mockContainerIPWithScope),
			})),
		}
		have, err := w.Tag(report.Report{
			Container: report.MakeTopology().AddNode(nodeID, report.MakeNodeWith(map[string]string{
				docker.ContainerID: mockContainerID,
			})),
		})
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(want, have) {
			t.Error(test.Diff(want, have))
		}
	}
}

const (
	mockHostID               = "host1"
	mockWeavePeerName        = "winnebago"
	mockWeavePeerNickName    = "winny"
	mockContainerID          = "83183a667c01"
	mockContainerMAC         = "d6:f2:5a:12:36:a8"
	mockContainerIP          = "10.0.0.123"
	mockContainerIPWithScope = ";10.0.0.123"
	mockHostname             = "hostname.weave.local"
)

var (
	mockResponse = fmt.Sprintf(`{
		"Router": {
			"Peers": [{
				"Name": "%s",
				"Nickname": "%s"
			}]
		},
		"DNS": {
			"Entries": [{
				"ContainerID": "%s",
				"Hostname": "%s.",
				"Tombstone": 0
			}]
		}
	}`, mockWeavePeerName, mockWeavePeerNickName, mockContainerID, mockHostname)
)

func mockWeaveRouter(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte(mockResponse)); err != nil {
		panic(err)
	}
}

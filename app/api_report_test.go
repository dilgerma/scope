package app_test

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"

	"github.com/dilgerma/scope/app"
	"github.com/dilgerma/scope/report"
)

func topologyServer() *httptest.Server {
	router := mux.NewRouter()
	app.RegisterTopologyRoutes(StaticReport{}, router)
	return httptest.NewServer(router)
}

func TestAPIReport(t *testing.T) {
	ts := topologyServer()
	defer ts.Close()

	is404(t, ts, "/api/report/foobar")

	var body = getRawJSON(t, ts, "/api/report")
	// fmt.Printf("Body: %v\n", string(body))
	var r report.Report
	err := json.Unmarshal(body, &r)
	if err != nil {
		t.Fatalf("JSON parse error: %s", err)
	}
}

package app_test

import (
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/ugorji/go/codec"

	"github.com/dilgerma/scope/app"
	"github.com/dilgerma/scope/report"
)

func topologyServer() *httptest.Server {
	handler := app.TopologyHandler(StaticReport{}, mux.NewRouter(), nil)
	return httptest.NewServer(handler)
}

func TestAPIReport(t *testing.T) {
	ts := topologyServer()
	defer ts.Close()

	is404(t, ts, "/api/report/foobar")

	var body = getRawJSON(t, ts, "/api/report")
	// fmt.Printf("Body: %v\n", string(body))
	var r report.Report

	decoder := codec.NewDecoderBytes(body, &codec.JsonHandle{})
	if err := decoder.Decode(&r); err != nil {
		t.Fatalf("JSON parse error: %s", err)
	}
}

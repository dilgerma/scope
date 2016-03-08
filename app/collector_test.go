package app_test

import (
	"testing"
	"time"

	"golang.org/x/net/context"

	"github.com/dilgerma/scope/app"
	"github.com/dilgerma/scope/report"
	"github.com/dilgerma/scope/test"
	"github.com/dilgerma/scope/test/reflect"
)

func TestCollector(t *testing.T) {
	ctx := context.Background()
	window := time.Millisecond
	c := app.NewCollector(window)

	r1 := report.MakeReport()
	r1.Endpoint.AddNode("foo", report.MakeNode())

	r2 := report.MakeReport()
	r2.Endpoint.AddNode("bar", report.MakeNode())

	if want, have := report.MakeReport(), c.Report(ctx); !reflect.DeepEqual(want, have) {
		t.Error(test.Diff(want, have))
	}

	c.Add(ctx, r1)
	if want, have := r1, c.Report(ctx); !reflect.DeepEqual(want, have) {
		t.Error(test.Diff(want, have))
	}

	c.Add(ctx, r2)

	merged := report.MakeReport()
	merged = merged.Merge(r1)
	merged = merged.Merge(r2)
	if want, have := merged, c.Report(ctx); !reflect.DeepEqual(want, have) {
		t.Error(test.Diff(want, have))
	}
}

func TestCollectorWait(t *testing.T) {
	ctx := context.Background()
	window := time.Millisecond
	c := app.NewCollector(window)

	waiter := make(chan struct{}, 1)
	c.WaitOn(ctx, waiter)
	defer c.UnWait(ctx, waiter)
	c.(interface {
		Broadcast()
	}).Broadcast()

	select {
	case <-waiter:
	default:
		t.Fatal("Didn't unblock")
	}
}

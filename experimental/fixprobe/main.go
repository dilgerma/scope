// Publish a fixed report.
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/ugorji/go/codec"

	"github.com/dilgerma/scope/common/xfer"
	"github.com/dilgerma/scope/probe/appclient"
	"github.com/dilgerma/scope/report"
)

func main() {
	var (
		publish         = flag.String("publish", fmt.Sprintf("localhost:%d", xfer.AppPort), "publish target")
		publishInterval = flag.Duration("publish.interval", 1*time.Second, "publish (output) interval")
	)
	flag.Parse()

	if len(flag.Args()) != 1 {
		log.Fatal("usage: fixprobe [--args] report.json")
	}

	b, err := ioutil.ReadFile(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}

	var fixedReport report.Report
	decoder := codec.NewDecoderBytes(b, &codec.JsonHandle{})
	if err := decoder.Decode(&fixedReport); err != nil {
		log.Fatal(err)
	}

	client, err := appclient.NewAppClient(appclient.ProbeConfig{
		Token:    "fixprobe",
		ProbeID:  "fixprobe",
		Insecure: false,
	}, *publish, *publish, nil)
	if err != nil {
		log.Fatal(err)
	}

	rp := appclient.NewReportPublisher(client)
	for range time.Tick(*publishInterval) {
		rp.Publish(fixedReport)
	}
}

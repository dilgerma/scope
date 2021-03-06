package appclient

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"net"
	"net/http"

	"github.com/certifi/gocertifi"
	"$GITHUB_URI/common/xfer"
)

var certPool *x509.CertPool

func init() {
	var err error
	certPool, err = gocertifi.CACerts()
	if err != nil {
		panic(err)
	}
}

// ProbeConfig contains all the info needed for a probe to do HTTP requests
type ProbeConfig struct {
	Token    string
	ProbeID  string
	Insecure bool
}

func (pc ProbeConfig) authorizeHeaders(headers http.Header) {
	headers.Set("Authorization", fmt.Sprintf("Scope-Probe token=%s", pc.Token))
	headers.Set(xfer.ScopeProbeIDHeader, pc.ProbeID)
}

func (pc ProbeConfig) authorizedRequest(method string, urlStr string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, urlStr, body)
	if err == nil {
		pc.authorizeHeaders(req.Header)
	}
	return req, err
}

func (pc ProbeConfig) getHTTPTransport(hostname string) (*http.Transport, error) {
	if pc.Insecure {
		return &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}, nil
	}

	host, _, err := net.SplitHostPort(hostname)
	if err != nil {
		return nil, err
	}

	return &http.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs:    certPool,
			ServerName: host,
		},
	}, nil
}

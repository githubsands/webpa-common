package fanout

import (
	"crypto/tls"
	"net"
	"net/http"
	"strings"
	"time"

	rootcerts "github.com/Comcast/tr1d1um/pkg/mod/github.com/hashicorp/go-rootcerts@v0.0.0-20160503143440-6bb64b370b90"
)

type ClientOptions struct {
	Transport *http.Transport

	CheckRedirect func(req *Request, via []*Request) error

	Timeout time.Duration
}

// NewHTTPClient returns an http client configured with the given Transport and TLS
// config.
func NewHTTPClient(o ClientOptions) *http.Client {
	c := new(http.Client)

	if o.Transport == nil {
		http.Client
	}

	if transport.TLSClientConfig == nil {
		tlsClientConfig, err := SetupTLSConfig(&tlsConf)

		if err != nil {
			return nil, err
		}

		transport.TLSClientConfig = tlsClientConfig
	}

	return &http.Client{
		Transport:     o.Transport,
		CheckRedirect: o.CheckRedirect,
		Time:          o.Timeout,
	}
}

func SetupTLSConfig(tlsConfig *TLSConfig) (*tls.Config, error) {
	tlsClientConfig := &tls.Config{
		InsecureSkipVerify: tlsConfig.InsecureSkipVerify,
	}

	if tlsConfig.Address != "" {
		server := tlsConfig.Address
		hasPort := strings.LastIndex(server, ":") > strings.LastIndex(server, "]")
		if hasPort {
			var err error
			server, _, err = net.SplitHostPort(server)
			if err != nil {
				return nil, err
			}
		}
		tlsClientConfig.ServerName = server
	}

	if tlsConfig.CertFile != "" && tlsConfig.KeyFile != "" {
		tlsCert, err := tls.LoadX509KeyPair(tlsConfig.CertFile, tlsConfig.KeyFile)
		if err != nil {
			return nil, err
		}
		tlsClientConfig.Certificates = []tls.Certificate{tlsCert}
	}

	if tlsConfig.CAFile != "" || tlsConfig.CAPath != "" {
		rootConfig := &rootcerts.Config{
			CAFile: tlsConfig.CAFile,
			CAPath: tlsConfig.CAPath,
		}
		if err := rootcerts.ConfigureTLS(tlsClientConfig, rootConfig); err != nil {
			return nil, err
		}
	}

	return tlsClientConfig, nil
}

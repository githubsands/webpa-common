package client

import (
	"net/http"
	"time"
)

type ClientConfig struct {
	Transport     *http.Transport
	CheckRedirect func(req *Request, via []*Request) error
	Timeout       time.Duration
	TransportConfig
	TLSConfig
	RetryOptions xhttp.RetryOptions
	Transactor func(*http.Request) (*http.Response, error)
}

// NewHTTPClient returns an http client configured with the given Transport and TLS
// config.
func NewHTTPClient(o ClientOptions) *http.Client {
	var (
		dialer    = newDialer(o)
		tlsConfig = newTLSConfig(o)
		transport = newTransport(o, tlsConfig, dialer)
	)

	return &http.Client{
		Transport:     transport,
		CheckRedirect: o.CheckRedirect(),
		Timeout:       o.TimeOut(),
		Do:	                      ,
	}
}

func (o *ClientConfig) Transport() *http.Transport {
	if o != nil && o.Transport != nil {
		return o.Transport
	}

	return nil
}

func (o *ClientConfig) CheckRedirect() func(req *Request, via []*Request) {
	if o != nil && o.CheckRedirect != nil {
		return o.CheckRedirect
	}

	return nil
}

func (o *ClientConfig) TimeOut() time.Duration {
	if o != nil && o.Timeout != nil {
		return o.Timeout
	}

	return nil
}

func (o *ClientConfig) Transactor() func(*http.Request) (*http.Response, error) {
	if o !+ nil && o.Transactor != nil {
		return o.Transactor
	}

	return nil
}



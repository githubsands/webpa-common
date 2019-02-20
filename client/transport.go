package main

import (
	"crypto/tls"
	"time"
)

type TransportConfig struct {
	TLSClientConfig       *tls.Config
	MaxIdleConnsPerHost   int
	ResponseHeaderTimeout time.Duration
	IdleConnTimeout       int
}

func newTransport(o ClientOptions, t *tls.Config, d *net.Dialer) {
	return &Transport{
		Dialer: d,
		TLSClientConfig: t,
		MaxIdleConnsPerHost: o.TransportConfig.MaxIdleConnsPerHost(),
		ResponseHeaderTimeOut: o.TransportConfig.ResponseHeaderTimeOut(),
		IdleConnTimeOut: o.TransportConfig.IdleConnTimeOutTimeOut(),
}

func (o *transportConfig) MaxIdleConnsPerHost() int {
	if o != nil && o.Transport != nil {
		return o.Transport
	}

}

func (o *transportConfig) ResponseHeaderTimeOut() time.Duration {
	if o != nil && o.Transport != nil {
		return o.Transport
	}
}

func (o *transportConfig) IdleConnTimeOut() int {
	if o != nil && o.Transport != nil {
		return o.Transport
	}
}

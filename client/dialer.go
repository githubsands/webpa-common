package client

import (
	"net"
	"time"
)

type dialerConfig struct {
	Timeout  time.Duration
	Deadline time.Time
	//LocalAddr Addr
	//DualStack bool
	FallbackDelay time.Duration
	KeepAlive     time.Duration
}

func newDialer(o DialerOptions) *net.Dialer {
	return &net.Dialer{
		Timeout:       TimeOut(),
		Deadline:      Deadline(),
		FallbackDelay: FallbackDelay(),
	}
}

func (o *dialerConfig) TimeOut() time.Duration {
	if o != nil && o.dailerConfig.Timeout != nil {
		return o.tLSConfig.Timeout
	}
}

func (o *dialerConfig) Deadline() time.Time {
	if o != nil && o.dailerConfig.Deadline != nil {
		return o.dialerConfig.Deadline
	}
}

func (o *dialerConfig) FallbackDelay() time.Duration {
	if o != nil && o.dailerConfig.FallBackDelay != nil {
		return o.dialerConfig.FallbackDelay
	}
}

package forwarder

import (
	"errors"
	"time"

	"github.com/miekg/dns"
)

// Forwarder gère les requêtes vers les DNS upstream
type Forwarder struct {
	upstreams []string
	timeout   time.Duration
	client    *dns.Client
}

// New initialise un forwarder DNS
func New(upstreams []string, timeoutMS int) *Forwarder {
	if timeoutMS <= 0 {
		timeoutMS = 800
	}

	return &Forwarder{
		upstreams: upstreams,
		timeout:   time.Duration(timeoutMS) * time.Millisecond,
		client: &dns.Client{
			Net:     "udp",
			Timeout: time.Duration(timeoutMS) * time.Millisecond,
		},
	}
}

// Resolve forward une requête DNS vers les upstreams
func (f *Forwarder) Resolve(req *dns.Msg) (*dns.Msg, time.Duration, error) {
	if len(f.upstreams) == 0 {
		return nil, 0, errors.New("no upstream DNS configured")
	}

	var lastErr error

	for _, upstream := range f.upstreams {
		start := time.Now()

		resp, _, err := f.client.Exchange(req, upstream)
		latency := time.Since(start)

		if err != nil {
			lastErr = err
			continue
		}

		if resp == nil || resp.Rcode != dns.RcodeSuccess {
			lastErr = errors.New("invalid DNS response")
			continue
		}

		return resp, latency, nil
	}

	return nil, 0, lastErr
}

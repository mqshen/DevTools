package traceroute

import (
	"encoding/json"
	"fmt"
	"net"
	"net/netip"
	"time"
)

const (
	DesMinPort        = 33434
	DesMaxPort        = 33534
	DefaultFirstTTL   = 1
	DefaultMaxTTL     = 64
	DefaultMinWaitSec = 1
	DefaultMaxWaitSec = 10
)

type TraceResult struct {
	TTL         int           `json:"ttl" yaml:"ttl"`
	NextHot     string        `json:"host" yaml:"host"`
	ElapsedTime time.Duration `json:"elapsedTime" yaml:"elapsedTime"`
	Replied     bool          `json:"replied" yaml:"replied"`
}

func (t TraceResult) MarshalJSON() (b []byte, err error) {
	isPrivate := false
	if t.NextHot != "*" {
		addr, err := netip.ParseAddr(t.NextHot)
		if err == nil {
			isPrivate = addr.IsPrivate()
		}
	}
	result := map[string]interface{}{
		"ttl":         t.TTL,
		"host":        t.NextHot,
		"isPrivate":   isPrivate,
		"elapsedTime": t.ElapsedTime / time.Microsecond,
		"replied":     t.Replied,
	}
	return json.Marshal(result)
}

type TraceConfig struct {
	FirstTTL int
	Retry    int
	MaxTTL   int
	Debug    bool
	WaitSec  int64
}

func (conf *TraceConfig) check() error {
	if conf.MaxTTL <= 0 {
		return fmt.Errorf("invalid max ttl: %d", conf.MaxTTL)
	}

	if conf.FirstTTL <= 0 {
		conf.FirstTTL = DefaultFirstTTL
	}

	if conf.MaxTTL > DefaultMaxTTL {
		conf.MaxTTL = DefaultMaxTTL
	}

	if conf.WaitSec <= 0 {
		conf.WaitSec = DefaultMinWaitSec
	} else if conf.WaitSec >= DefaultMaxWaitSec {
		conf.WaitSec = DefaultMaxWaitSec
	}

	return nil
}

func Traceroute(dest string) (<-chan *TraceResult, <-chan error, error) {
	conf := &TraceConfig{
		FirstTTL: DefaultFirstTTL,
		MaxTTL:   DefaultMaxTTL,
		Retry:    1,
		WaitSec:  DefaultMinWaitSec,
	}
	addr, err := netip.ParseAddr(dest)
	var ips []net.IP
	if err != nil {
		ips, err = net.LookupIP(dest)
		if len(ips) > 0 {
			addr, err = netip.ParseAddr(ips[0].String())
		}
	}
	if err != nil {
		return nil, nil, err
	}

	if addr.Is4() {
		routeChan, errChan := trace4(conf, addr)
		return routeChan, errChan, nil
	} else {
		routeChan, errChan := trace6(conf, addr)
		return routeChan, errChan, nil
	}
}

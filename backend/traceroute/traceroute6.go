package traceroute

import (
	"fmt"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv6"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sys/unix"
	"net/netip"
	"time"
)

func trace6(conf *TraceConfig, addr netip.Addr) (<-chan *TraceResult, <-chan error) {
	ttl := conf.FirstTTL
	try := conf.Retry
	destPort := DesMinPort
	destAddr := addr.As16()

	var (
		eg        errgroup.Group
		routeChan = make(chan *TraceResult)
		errChan   = make(chan error)
	)

	eg.Go(func() error {
		recvSocket, err := unix.Socket(unix.AF_INET6, unix.SOCK_RAW, unix.IPPROTO_ICMPV6)
		if err != nil {
			errChan <- err
			return nil
		}
		sendSocket, err := unix.Socket(unix.AF_INET6, unix.SOCK_DGRAM, unix.IPPROTO_UDP)
		if err != nil {
			errChan <- err
			return nil
		}
		defer unix.Close(recvSocket)
		defer unix.Close(sendSocket)

		if err := unix.SetsockoptTimeval(recvSocket, unix.SOL_SOCKET, unix.SO_RCVTIMEO,
			&unix.Timeval{Sec: conf.WaitSec, Usec: 0}); err != nil {
			errChan <- err
			return nil
		}

		//var results []TraceResult
		begin := time.Now()

		for {
			begin = time.Now()
			if err := unix.SetsockoptInt(sendSocket, unix.IPPROTO_IPV6, unix.IPV6_UNICAST_HOPS, ttl); err != nil {
				errChan <- err
				return nil
			}
			if err := unix.Sendto(sendSocket, []byte("hello"), 0, &unix.SockaddrInet6{Port: destPort, Addr: destAddr}); err != nil {
				errChan <- err
				return nil
			}

			var p = make([]byte, 4096)
			result := &TraceResult{TTL: ttl, ElapsedTime: time.Since(begin), Replied: false}
			n, from, err := unix.Recvfrom(recvSocket, p, 0)
			if err == nil {
				try = 0
				icmpReply, err := icmp.ParseMessage(58, p[:n])
				if err != nil {
					errChan <- err
					return nil
				}

				if icmpReply.Type == ipv6.ICMPTypeTimeExceeded || icmpReply.Type == ipv6.ICMPTypeDestinationUnreachable {
					fromAddr := from.(*unix.SockaddrInet6).Addr
					result.Replied = true
					result.NextHot = netip.AddrFrom16(fromAddr).String()
					//results = append(results, result)
					routeChan <- result
					if conf.Debug {
						fmt.Printf("ttl %d receive from:%v time:%v icmpReply:%+v\n", ttl, result.NextHot, time.Since(begin), icmpReply)
					}

					if icmpReply.Type == ipv6.ICMPTypeTimeExceeded {
						ttl++
					}
					if icmpReply.Type == ipv6.ICMPTypeDestinationUnreachable || ttl > conf.MaxTTL || fromAddr == destAddr {
						return nil
					}
				}
			} else {
				if conf.Debug {
					fmt.Printf("ttl %d * err: %s \n", ttl, err.Error())
				}
				result.NextHot = "*"
				routeChan <- result
				try++
				if try > conf.Retry {
					try = 0
					ttl++
				}
			}

			destPort++
			if destPort > DesMaxPort {
				destPort = DesMinPort
			}
		}
	})
	go func() {
		// Wait for the aggregation routines to shut down to avoid writing to closed channels
		_ = eg.Wait() // will never error
		close(routeChan)
		close(errChan)
	}()
	return routeChan, errChan
}

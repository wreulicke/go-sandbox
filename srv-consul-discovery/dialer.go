// This file is forked from https://github.com/koron-go/dialsrv/blob/d374719f1db5f61a200c1f884dd001ae1baccdbf/dial.go

package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
)

// Dialer wraps net.Dialer with SRV lookup.
type Dialer struct {
	nd *net.Dialer
}

// New creates a new Dialer with base *net.Dialer.
func New(d *net.Dialer) *Dialer {
	if d == nil {
		d = &net.Dialer{}
	}
	return &Dialer{nd: d}
}

// Dial connects to the address on the named network.
func (d *Dialer) Dial(network, address string) (net.Conn, error) {
	return d.DialContext(context.Background(), network, address)
}

// DialContext connects to the address on the named network using
// the provided context.
func (d *Dialer) DialContext(ctx context.Context, network, address string) (net.Conn, error) {
	if fa := parseAddr(network, address); fa != nil {
		return d.dialSRV(ctx, fa)
	}
	return d.nd.DialContext(ctx, network, address)
}

func (d Dialer) dialSRV(ctx context.Context, fa *FlavoredAddr) (net.Conn, error) {
	r := d.nd.Resolver
	if r == nil {
		r = net.DefaultResolver
	}
	host, err := splitHost(fa.Name)
	if err != nil {
		return nil, err
	}
	_, addrs, err := r.LookupSRV(ctx, fa.Service, fa.Proto, host)
	if err != nil {
		log.Println("lookupSRV is failed", err)
		return nil, err
	}
	if len(addrs) == 0 {
		return nil, fmt.Errorf("no SRV records for %s", fa.String())
	}
	return d.nd.DialContext(ctx, fa.Network, address(addrs[0]))
}

func splitHost(s string) (string, error) {
	if strings.IndexByte(s, ':') < 0 {
		return s, nil
	}
	h, _, err := net.SplitHostPort(s)
	return h, err
}

// FlavoredAddr represents SRV flavored address.
type FlavoredAddr struct {
	Network string
	Service string
	Proto   string
	Name    string
}

func parseAddr(network, address string) *FlavoredAddr {
	const prefix = "srv+"
	if !strings.HasPrefix(address, prefix) {
		return nil
	}
	address = address[len(prefix):]
	n := strings.Index(address, "+")
	if n < 0 {
		return &FlavoredAddr{Network: network, Name: address}
	}
	return &FlavoredAddr{
		Network: network,
		Service: address[:n],
		Proto:   network,
		Name:    address[n+1:],
	}
}

// String returns FlavoredAddr's string representation.
func (fa *FlavoredAddr) String() string {
	if fa.Service == "" && fa.Proto == "" {
		return fa.Name
	}
	return "_" + fa.Service + "._" + fa.Proto + "." + fa.Name
}

func address(srv *net.SRV) string {
	return srv.Target + ":" + strconv.FormatUint(uint64(srv.Port), 10)
}

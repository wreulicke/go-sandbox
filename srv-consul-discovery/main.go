package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httptrace"
	"time"
)

var d = net.Dialer{}
var resolver = &net.Resolver{
	PreferGo: true,
	Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
		log.Println("dialing", network, address)
		return d.DialContext(ctx, network, ":8600")
	},
}

// HTTPTransport is replacement for http.DefaultTransport
var HTTPTransport = &http.Transport{
	Proxy: http.ProxyFromEnvironment,
	DialContext: New(&net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
		DualStack: true,
		Resolver:  resolver,
	}).DialContext,
	MaxIdleConns:          100,
	IdleConnTimeout:       90 * time.Second,
	TLSHandshakeTimeout:   10 * time.Second,
	ExpectContinueTimeout: 1 * time.Second,
}

// HTTPClient is replacement for http.DefaultClient
var HTTPClient = &http.Client{
	Transport: HTTPTransport,
}

func main() {
	trace := &httptrace.ClientTrace{
		DNSStart: func(dnsInfo httptrace.DNSStartInfo) {
			fmt.Printf("DNS Start Info: %+v\n", dnsInfo)
		},
		// GetConn: func(hostPort string) {
		// 	fmt.Printf("Get Conn: %s\n", hostPort)
		// },
		// GotConn: func(connInfo httptrace.GotConnInfo) {
		// 	fmt.Printf("Got Conn: %s\n", connInfo.Conn.RemoteAddr())
		// },
		DNSDone: func(dnsInfo httptrace.DNSDoneInfo) {
			fmt.Printf("DNS Done Info: %+v\n", dnsInfo)
		},
		// ConnectStart: func(network, addr string) {
		// 	fmt.Printf("Connect Start: %s:%s\n", network, addr)
		// },
		// PutIdleConn: func(err error) {
		// 	fmt.Printf("PutIdleConn: %+v\n", err)
		// },
		// ConnectDone: func(network, addr string, err error) {
		// 	fmt.Printf("Connect Done: %s:%s %+v\n", network, addr, err)
		// },
		// Wait100Continue: func() {
		// 	fmt.Printf("Wait100Continue")
		// },
	}
	fmt.Println(trace)
	http.Handle("/", http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		httpReq, _ := http.NewRequest("GET", "http://srv+service.service.consul/backend", nil)
		httpRes, err := HTTPClient.Do(httpReq.WithContext(httptrace.WithClientTrace(httpReq.Context(), trace)))
		// httpRes, err := HTTPClient.Do(httpReq)
		if err != nil {
			log.Println(err)
			res.WriteHeader(500)
			return
		}
		defer httpRes.Body.Close()
		io.Copy(res, httpRes.Body)
	}))
	http.Handle("/backend", http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(200)
		res.Write([]byte("ok"))
	}))
	http.ListenAndServe(":8080", nil)
}

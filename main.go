package main

import (
	"context"
	"net"
	"time"
)

func extip() string {

	r := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{Timeout: time.Millisecond * time.Duration(10000)}
			return d.DialContext(ctx, network, "resolver1.opendns.com:53")
		},
	}

	ip, _ := r.LookupHost(context.Background(), "myip.opendns.com")

	return ip[0]

}

func main() {
	println(extip())
}

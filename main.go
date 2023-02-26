package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

const resolverTimeout = 10000
const defaultPort = "53"

const googleResolver = "ns1.google.com"
const openDNSResolver = "resolver1.opendns.com"

const googleHost = "o-o.myaddr.1.google.com"
const openDNSHost = "myip.opendns.com"

func UseResolver(resolver string) *net.Resolver {

	return &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{Timeout: time.Millisecond * time.Duration(resolverTimeout)}
			return d.DialContext(ctx, network, resolver+":"+defaultPort)
		},
	}

}

func GoogleExtIP() string {

	r := UseResolver(googleResolver)
	txts, _ := r.LookupTXT(context.Background(), googleHost)

	return txts[0]

}

func OpenDNSExtIP() string {

	r := UseResolver(openDNSResolver)
	ip, _ := r.LookupHost(context.Background(), openDNSHost)

	return ip[0]

}

func main() {

	opendns := OpenDNSExtIP()
	google := GoogleExtIP()

	if google == opendns {
		fmt.Println(opendns)
	} else {

		l := log.New(os.Stderr, "", 1)

		l.Println("Google and OpenDNS have different ideas regarding your external IP address.")
		l.Printf("Google thinks it is: %s\n", google)
		l.Printf("OpenDNS thinks it is: %s\n", opendns)

		os.Exit(1)

	}

}

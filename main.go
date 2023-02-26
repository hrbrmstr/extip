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
			return d.DialContext(ctx, network, resolver + ":" + defaultPort)
		},
	}

}

func GoogleExtIP() ([]string, error) {

	r := UseResolver(googleResolver)
	txts, err := r.LookupTXT(context.Background(), googleHost)

	return txts, err

}

func OpenDNSExtIP() ([]string, error) {

	r := UseResolver(openDNSResolver)
	ip, err := r.LookupHost(context.Background(), openDNSHost)

	return ip, err

}

func main() {

  l := log.New(os.Stderr, "", 1)

	opendns, oerr := OpenDNSExtIP()
	google, gerr := GoogleExtIP()

  if (oerr != nil) && (gerr != nil) {
    l.Println("Neither DNS resolution worked.")
		os.Exit(2)
	} else if (oerr != nil) {
		l.Println("OpenDNS resolver query failed.")
		fmt.Println(google[0])
	} else if (gerr != nil) {
		l.Println("Google resolver query failed")
		fmt.Println(opendns[0])
	} else if google[0] == opendns[0] {
		fmt.Println(opendns[0])
	} else {

		l.Println("Google and OpenDNS have different ideas regarding your external IP address.")
		l.Printf("Google thinks it is: %s\n", google)
		l.Printf("OpenDNS thinks it is: %s\n", opendns)

		os.Exit(1)

	}

}

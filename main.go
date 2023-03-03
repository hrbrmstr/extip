package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/alexflint/go-arg"
)

const(
	version = "0.2.0"
  resolverTimeout = 10000
  defaultPort = 53

  googleResolver = "ns1.google.com"
  openDNSResolver = "resolver1.opendns.com"

 googleHost = "o-o.myaddr.1.google.com"
 openDNSHost = "myip.opendns.com"
 akamaiHost = "whoami.ds.akahelp.net"
 akamaiDomain = "akamai.com"
)

// Test if all strings in a list are equal
func AllEqual(a []string) bool {
	for i := 1; i < len(a); i++ {
		if a[i] != a[0] {
			return false
		}
	}
	return true
}

// Remove last char in a string if it is `suffix`
func TrimSuffix(s, suffix string) string {
  hasSuf := strings.HasSuffix(s, suffix)
	if hasSuf {
    s = s[:len(s)-len(suffix)]
  }
  return s
}

// Setup a specific resovler to use for DNS lookups
func UseResolver(resolver string, port int) *net.Resolver {
	
	return &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{Timeout: time.Millisecond * time.Duration(resolverTimeout)}
			return d.DialContext(ctx, network, resolver + ":" + strconv.Itoa(port))
		},
	}
	
}

// Get one of Akamai's authoritative nameservers
func AkamaiResolver() (string, error) {
	
	ns, err := net.LookupNS(akamaiDomain)
	
	nsHost := ""
	if err == nil {
		nsHost = TrimSuffix((*ns[0]).Host, ".")
	}

	return nsHost, err
	
}

// Get our local, external IP via Akamai's hack
// 
// Akamai returns "ns" "ip.ad.dr.ess" so we have to get rid of cruft
func AkamaiExtIP() ([]string, error) {

	regEx := regexp.MustCompile(`[^0-9\.\:]`)
	akamaiResolver, err := AkamaiResolver()
	
	if (err != nil) {
		return []string{""}, err
	}

	r := UseResolver(akamaiResolver, defaultPort)
	txts, err := r.LookupTXT(context.Background(), akamaiHost)

	if (err != nil) {
		return []string{""}, err
	}
	
	txts[0] = regEx.ReplaceAllString(txts[0], "")
	
	return txts, err
	
}

// Get our local, external IP via Google's hack
func GoogleExtIP() ([]string, error) {
	
	r := UseResolver(googleResolver, defaultPort)
	txts, err := r.LookupTXT(context.Background(), googleHost)
	
	if (err != nil) {
		txts = []string{""}
	}

	return txts, err
	
}

// Get our local, external IP via OpenDNS's hack
func OpenDNSExtIP() ([]string, error) {
	
	r := UseResolver(openDNSResolver, defaultPort)
	ips, err := r.LookupHost(context.Background(), openDNSHost)

	if (err != nil) {
		ips = []string{""}
	}

	return ips, err
	
}

// Get our local, external IP via custom extip resolver
func CustomExtIP(resolver string, host string, rectype string, port int) ([]string, error) {
	
	r := UseResolver(resolver, port)

	var res [](string)
	var err error

	if (rectype == "A") {
	  res, err = r.LookupHost(context.Background(), host)
	} else {
		res, err = r.LookupTXT(context.Background(), host)
	}

	if (err != nil) {
		res = []string{""}
	}

	return res, err
	
}

type args struct {
	Server string `arg:"-s,--server,env:EXTIP_SERVER" help:"extip-svr IP/FQDN. e.g., ip.rudis.net" placeholder:"EXTIP_SERVER"`
	Domain string `arg:"-d,--domain,env:EXTIP_DOMAIN" help:"Domain to use for IP Lookup. e.g., myip.is" placeholder:"DOMAIN"`
	RecordType string `arg:"-r,--record-type,env:EXTIP_RECORD_TYPE" help:"DNS record type to lookup. One of TXT or A." placeholder:"RECORD_TYPE" default:"TXT"`
	Port int    `arg:"-p,--port,env:EXTIP_PORT" help:"Port extip resolver is listening on." placeholder:"PORT" default:"53"`
}

func (args) Description() string {
  return "Lookup external IP address via DNS.\n\nDefaults to using Akamai, OpenDNS, and Google services.\nYou can specify an extip server via the following command line options.\nNOTE: Both server and domain should be specified to override default behavior.\nMore info about running an extip server can be found at <https://github.com/hrbrmstr/extip-svr>.\n"
}

func (args) Version() string {
	return "extip " + version
}

// TODO: Make this a proper CLI with cmdline options since we have 3 services
func main() {

	l := log.New(os.Stderr, "", 1)

	var args args

  arg.MustParse(&args)

  if args.Server != "" && args.Domain != "" {
	  custom, cerr := CustomExtIP(args.Server, args.Domain, args.RecordType, args.Port)
  	if (cerr != nil) {
	  	l.Println("No DNS resolutions worked.")
		  os.Exit(3)
	  } else {
			fmt.Println(custom[0])
			os.Exit(0)
		}

	}
	
	opendns, oerr := OpenDNSExtIP()
	google, gerr := GoogleExtIP()
	akamai, aerr := AkamaiExtIP()
	
	if (oerr != nil) && (gerr != nil) && (aerr != nil) {
		l.Println("No DNS resolutions worked.")
		os.Exit(2)
	}
	
	if (oerr != nil) {
		l.Println("OpenDNS resolver query failed.")
		opendns[0] = "FAILED"
	}
	
	if (gerr != nil) {
		l.Println("Google resolver query failed")
		google[0] = "FAILED"
	}
	
	if (aerr != nil) {
		l.Println("Akamai resovler query failed")
		akamai[0] = "FAILED"
	}
	
	// If at least one worked, compare the three; if not all equal then error out
	// otherwise return one of them.

	if AllEqual([]string{akamai[0], google[0], opendns[0]}) {
		
		fmt.Println(opendns[0])
		
	} else {
			
	  l.Println("Resolvers have different ideas regarding your external IP address.")
		l.Printf("Akamai thinks it is: %s\n", akamai[0])
		l.Printf("Google thinks it is: %s\n", google[0])
		l.Printf("OpenDNS thinks it is: %s\n", opendns[0])
			
		os.Exit(1)
			
	}
		
}
	
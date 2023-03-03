# extip - retrieve your external IP address via DNS

Small Golang package/cli that uses special DNS resolvers to return your external IP address.

By default it uses Google, OpenDNS, and Akamai. If there is a conflict between the resolver answers a message will be delivered on stderr with the conflicting values.

Alternatively, you can specify an [extip server](https://github.com/hrbrmstr/extip-svr) to use. See below for how to do that.

## References

- [Bizarre and Unusual Uses of DNS](https://fosdem.org/2023/schedule/event/dns_bizarre_and_unusual_uses_of_dns/)
- [Akamai blog](https://www.akamai.com/blog/developers/introducing-new-whoami-tool-dns-resolver-information)

## Build

```
just build # requires https://github.com/casey/just
```

## Install

```
go install -ldflags "-s -w" github.com/hrbrmstr/extip@latest
```

## Options

```
Lookup external IP address via DNS.

Defaults to using Akamai, OpenDNS, and Google services.
You can specify an extip server via the following command line options.
NOTE: Both server and domain should be specified to override default behavior.
More info about running an extip server can be found at <https://github.com/hrbrmstr/extip-svr>.

extip 0.2.0
Usage: extip [--server EXTIP_SERVER] [--domain DOMAIN] [--record-type RECORD_TYPE] [--port PORT]

Options:
  --server EXTIP_SERVER, -s EXTIP_SERVER
                         extip-svr IP/FQDN. e.g., ip.rudis.net [env: EXTIP_SERVER]
  --domain DOMAIN, -d DOMAIN
                         Domain to use for IP Lookup. e.g., myip.is [env: EXTIP_DOMAIN]
  --record-type RECORD_TYPE, -r RECORD_TYPE
                         DNS record type to lookup. One of TXT or A. [default: TXT, env: EXTIP_RECORD_TYPE]
  --port PORT, -p PORT   Port extip resolver is listening on. [default: 53, env: EXTIP_PORT]
  --help, -h             display this help and exit
```

## Usage

Default:

```
extip
```

Use an extip server:

```
extip -s ip.rudis.net -d ip.is
```
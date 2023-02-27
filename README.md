# extip - retrieve your external IP address via DNS

Small Golang package/cli that uses special DNS resolvers to return your external IP address.

Presently it uses Google, OpenDNS, and Akamai. If there is a conflict between the resolver answers a message will be delivered on stderr with the conflicting values.

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

## TODO

- [ ] Spruce up the CLI to allow folks to choose from the existing list vs use all of them
- [ ] Allow folks to specify an endpoint and record type so it can be used with the companion server <https://github.com/hrbrmstr/extip>

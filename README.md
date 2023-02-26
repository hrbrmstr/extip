# extip - retrieve your external IP address via DNS

Small Golang package/cli that uses special DNS resolvers to return your external IP address.

Presently it uses Google and OpenDNS. If there is a conflict between the resolver answers a message will be delivered on stderr with the conflicting values.

## Build

```
just build # requires https://github.com/casey/just
```

## Install

```
go install -ldflags "-s -w" github.com/hrbrmstr/extip
```

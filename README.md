# extip - retrieve your external IP address via DNS

Small Golang package/cli that uses a special DNS resolver to return your external IP address.

## Build

```
just build # requires https://github.com/casey/just
```

## Install

```
go install -ldflags "-s -w" github.com/hrbrmstr/extip
```

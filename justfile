#!/usr/bin/env just --justfile

set shell := ["zsh", "-cu"] 

# Lists the justfile commands
@default:
  @just --list

# Build the package
@build:
  go build -ldflags "-s -w" -o bin/extip

@run: build
	./extip

# Be a good citizen
@fmt:
  go fmt

# Build for other platforms
@release-build:
  GOOS=windows GOARCH=386 go build -ldflags "-s -w" -o bin/extip.exe main.go && cd bin && zip extip.zip extip.exe && rm extip.exe
  GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o bin/extip-linux-amd64 main.go && cd bin && gzip extip-linux-amd64
  GOOS=linux GOARCH=arm64 go build -ldflags "-s -w" -o bin/extip-linux-arm64 main.go && cd bin && gzip extip-linux-arm64
  GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w" -o bin/extip-darwin-amd64 main.go && cd bin && gzip extip-darwin-amd64
  GOOS=darwin GOARCH=arm64 go build -ldflags "-s -w" -o bin/extip-darwin-arm64 main.go && cd bin && gzip extip-darwin-arm64

# release binaries
@release:
  gh workflow run release-binaries.yml   

# Check results against dig. Requires dig.
@test: build
  [ "$(dig myip.opendns.com @resolver1.opendns.com +short)" = "$(./bin/extip)" ] && echo "Passed OpenDNS test"
  [ "$(dig o-o.myaddr.1.google.com @ns1.google.com TXT +short | tr -d '"')" = "$(./bin/extip)" ] && echo "Passed Google test"
  [ "$(dig +short TXT whoami.ds.akahelp.net @$(dig +short +answer NS akamai.com | head -1) | grep ns | sed -e 's/[^0-9\.\:]//g')" = "$(./bin/extip)" ] && echo "Passed Akamai test"
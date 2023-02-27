#!/usr/bin/env just --justfile

set shell := ["zsh", "-cu"] 

# Lists the justfile commands
@default:
  @just --list

# Build the package
@build:
  go build -ldflags "-s -w"

@run: build
	./extip

# Be a good citizen
@fmt:
  go fmt

# Check results against dig. Requires dig.
@test: build
  [ "$(dig myip.opendns.com @resolver1.opendns.com +short)" = "$(./extip)" ] && echo "Passed OpenDNS test"
  [ "$(dig o-o.myaddr.1.google.com @ns1.google.com TXT +short | tr -d '"')" = "$(./extip)" ] && echo "Passed Google test"
  [ "$(dig +short TXT whoami.ds.akahelp.net @$(dig +short +answer NS akamai.com | head -1) | grep ns | sed -e 's/[^0-9\.\:]//g')" = "$(./extip)" ] && echo "Passed Akamai test"
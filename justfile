#!/usr/bin/env just --justfile

set shell := ["zsh", "-cu"] 

# Lists the justfile commands
@default:
  @just --list

# Build the package
@build:
  go build -ldflags "-s -w"

# Be a good citizen
@fmt:
  go fmt

# Check results against dig. Requires dig.
@test: build
  [ "$(dig myip.opendns.com @resolver1.opendns.com +short)" = "$(./extip)" ] && echo "Passed"
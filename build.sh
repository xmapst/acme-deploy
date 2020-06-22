#!/bin/sh

# disable go modules
export GOPATH=""

# disable cgo
export CGO_ENABLED=0

set -e
set -x

# linux
GOOS=linux GOARCH=amd64 go build -o release/linux/amd64/acme
tar -cvzf release/acme_linux_amd64.tar.gz -C release/linux/amd64 acme
GOOS=linux GOARCH=arm64 go build -o release/linux/arm64/acme
tar -cvzf release/acme_linux_arm64.tar.gz -C release/linux/arm64 acme
GOOS=linux GOARCH=arm   go build -o release/linux/arm/acme
tar -cvzf release/acme_linux_arm.tar.gz   -C release/linux/arm   acme
GOOS=linux GOARCH=386   go build -o release/linux/386/acme
tar -cvzf release/acme_linux_386.tar.gz   -C release/linux/386   acme

# windows
GOOS=windows GOARCH=amd64 go build -o release/windows/amd64/acme.exe
tar -cvzf release/acme_windows_amd64.tar.gz -C release/windows/amd64 acme.exe
GOOS=windows GOARCH=386 go build -o release/windows/386/acme.exe
tar -cvzf release/acme_windows_386.tar.gz   -C release/windows/386   acme.exe

# darwin
GOOS=darwin GOARCH=amd64 go build -o release/darwin/amd64/acme
tar -cvzf release/acme_darwin_amd64.tar.gz -C release/darwin/amd64  acme

# freebsd
GOOS=freebsd GOARCH=amd64 go build -o release/freebsd/amd64/acme
tar -cvzf release/acme_freebsd_amd64.tar.gz -C release/freebsd/amd64 acme
GOOS=freebsd GOARCH=arm   go build -o release/freebsd/arm/acme
tar -cvzf release/acme_freebsd_arm.tar.gz   -C release/freebsd/arm   acme
GOOS=freebsd GOARCH=386   go build -o release/freebsd/386/acme
tar -cvzf release/acme_freebsd_386.tar.gz   -C release/freebsd/386   acme

# netbsd
GOOS=netbsd GOARCH=amd64 go build -o release/netbsd/amd64/acme
tar -cvzf release/acme_netbsd_amd64.tar.gz -C release/netbsd/amd64 acme
GOOS=netbsd GOARCH=arm   go build -o release/netbsd/arm/acme
tar -cvzf release/acme_netbsd_arm.tar.gz   -C release/netbsd/arm   acme

# openbsd
GOOS=openbsd GOARCH=amd64 go build -o release/openbsd/amd64/acme
tar -cvzf release/acme_openbsd_amd64.tar.gz -C release/openbsd/amd64 acme
GOOS=openbsd GOARCH=arm   go build -o release/openbsd/arm/acme
tar -cvzf release/acme_openbsd_arm.tar.gz   -C release/openbsd/arm   acme
GOOS=openbsd GOARCH=386   go build -o release/openbsd/386/acme
tar -cvzf release/acme_openbsd_386.tar.gz   -C release/openbsd/386   acme

# dragonfly
GOOS=dragonfly GOARCH=amd64 go build -o release/dragonfly/amd64/acme
tar -cvzf release/acme_dragonfly_amd64.tar.gz -C release/dragonfly/amd64  acme

# solaris
GOOS=solaris GOARCH=amd64 go build -o release/solaris/amd64/acme
tar -cvzf release/acme_solaris_amd64.tar.gz -C release/solaris/amd64  acme

# generate shas for tar files
shasum release/*.tar.gz > release/acme_checksums.txt

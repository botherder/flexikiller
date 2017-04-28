#!/bin/bash

# go get github.com/Sirupsen/logrus
# go get github.com/mattn/go-colorable
# go get github.com/mattn/go-sqlite3
# go get golang.org/x/sys/windows/registry

GO_VERSION="$(go version)"
GO_VERSION="$(echo $GO_VERSION | awk '{print $3}')"
if [[ $GO_VERSION = "go1.8" ]]; then
	$GOPATH/bin/rsrc -manifest flexikiller.manifest -ico icon.ico -o rsrc.syso
	GOOS=windows GOARCH=386 CC=i686-w64-mingw32-gcc CGO_ENABLED=1 go build --ldflags '-s -w -extldflags "-static"' -o FlexiKiller.exe
else
	echo "Error: Build currently only works with go1.8"
fi

#!/bin/bash

if [ "$(which git)" != "" ]; then
    commit=$(git rev-parse HEAD | head -c 7);
    branch=$(git rev-parse --abbrev-ref HEAD);
    version="$commit ($branch) - UNOFFICIAL BUILD"
else
   version="UNOFFICIAL BUILD"
fi

dstArch="$GOARCH"
if [ "$dstArch" = "386" ]; then
    dstArch="i386"
fi

dstName="engarde"
if [ "$GOOS" = "windows" ]; then
    dstName="$dstName.exe"
fi

echo "Building for $GOOS $dstArch - ver. $version"
go build -ldflags "-s -w -X 'main.Version=$version'" -o "dist/$GOOS/$dstArch/$dstName" ./cmd/engarde

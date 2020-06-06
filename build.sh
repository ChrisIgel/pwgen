#!/bin/sh

APP=pwgen
VERSION=v1.0.0

mkdir -p bin dist

for OS in darwin linux windows
do
  for ARCH in amd64
  do
    export GOOS=$OS
    export GOARCH=$ARCH
    NAME=$APP-$VERSION-$GOOS-$GOARCH
    echo "Building $NAME"

    if [ $GOOS = "windows" ]; then
      go build -ldflags="-s -w" -o bin/$NAME.exe main.go
      zip dist/$NAME.zip bin/$NAME.exe
    else
      go build -ldflags="-s -w" -o bin/$NAME main.go
      tar -cvzf dist/$NAME.tar.gz bin/$NAME
    fi

    echo "Building $NAME finished!\n"
  done
done

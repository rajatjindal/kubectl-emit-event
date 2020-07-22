#!/bin/bash

version=$1

if [ -z "$version" ]; then
    echo "version is needed"
    exit -1
fi

rm -rf _dist || true
mkdir -p _dist
# env GOOS=windows GOARCH=amd64 go build -ldflags "-X github.com/rajatjindal/kubectl-whoami/cmd.Version=$version" -o _dist/kubectl-whoami-windows-amd64-$version.exe main.go
env GOOS=darwin  GOARCH=amd64 go build -ldflags "-w -s -X github.com/rajatjindal/kubectl-whoami/pkg/cmd.Version=$version" -o _dist/darwin-amd64/kubectl-whoami      main.go
cp LICENSE _dist/darwin-amd64/
cd _dist/darwin-amd64/ && tar -cvzf darwin-amd64-$version.tar.gz kubectl-whoami LICENSE
shasum -a 256 darwin-amd64-$version.tar.gz | awk '{print $1}' > sha256-darwin-amd64-$version
cd -

env GOOS=linux   GOARCH=amd64 go build -ldflags "-w -s -X github.com/rajatjindal/kubectl-whoami/pkg/cmd.Version=$version" -o _dist/linux-amd64/kubectl-whoami       main.go
cp LICENSE _dist/linux-amd64/
cd _dist/linux-amd64/  && tar -cvzf linux-amd64-$version.tar.gz kubectl-whoami LICENSE
shasum -a 256 linux-amd64-$version.tar.gz | awk '{print $1}' > sha256-linux-amd64-$version

cd -

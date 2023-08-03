#!/usr/bin/env bash

set +x
set -e

egg=dist/tfmock.egg-info
rm -fr dist
mkdir -p $egg

GOOS=windows GOARCH=amd64 go build -o $egg/bin/win32/az.exe ./modules/az/
GOOS=linux GOARCH=amd64 go build -o $egg/bin/linux/az ./modules/az/
GOOS=darwin GOARCH=amd64 go build -o $egg/bin/darwin/az ./modules/az/
chmod +x $egg/bin/linux/az
chmod +x $egg/bin/darwin/az

GOOS=windows GOARCH=amd64 go build -o $egg/bin/win32/CloudMock.exe
GOOS=linux GOARCH=amd64 go build -o $egg/bin/linux/CloudMock
GOOS=darwin GOARCH=amd64 go build -o $egg/bin/darwin/CloudMock

chmod +x $egg/bin/linux/CloudMock
chmod +x $egg/bin/darwin/CloudMock

mkdir -p $egg/bin/requests/
cp -fr requests $egg/bin/requests/
cp -fr *.md dist/
cp -fr *.pem $egg/bin/
cp -fr ./modules/az/provider_override.tf $egg/bin
cp -fr ./tfmock/* dist/

pushd dist/
python -m build
popd
#!/usr/bin/env bash

GOOS=windows GOARCH=amd64 go build -o ./build/windows/pcd2wc4.exe
GOOS=darwin GOARCH=amd64 go build -o ./build/darwin/pcd2wc4
GOOS=linux GOARCH=amd64 go build -o ./build/linux/pcd2wc4

cp README.md ./build/darwin
cp README.md ./build/linux
cp README.md ./build/windows
#!/usr/bin/env bash

rm -rf ./build

GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o ./build/windows/pcd2wc4.exe
GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o ./build/darwin/pcd2wc4
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ./build/linux/pcd2wc4

upx --brute ./build/darwin/pcd2wc4
upx --brute ./build/linux/pcd2wc4
upx --brute ./build/windows/pcd2wc4.exe

cp README.md ./build/darwin
cp README.md ./build/linux
cp README.md ./build/windows

zip -j -9 ./build/pcd2wc4_v1.1_darwin.zip ./build/darwin/*
zip -j -9 ./build/pcd2wc4_v1.1_linux.zip ./build/linux/*
zip -j -9 ./build/pcd2wc4_v1.1_windows.zip ./build/windows/*
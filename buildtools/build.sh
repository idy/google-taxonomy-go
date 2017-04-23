#!/bin/sh

go-bindata -o data/data.go -pkg data -ignore "\.go$" -nometadata -prefix data data

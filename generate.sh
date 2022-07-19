#!/bin/bash
export PATH=$(go env GOPATH)/bin:$PATH
protoc -I . ./application.proto --go_out=:.
protoc -I . ./application.proto --go-grpc_out=:.


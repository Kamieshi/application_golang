#!/bin/bash
export PATH=$(go env GOPATH)/bin:$PATH
swag init -g main.go --output docs/app


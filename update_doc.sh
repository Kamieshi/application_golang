#!/bin/bash
export PATH=$(go env GOPATH)/bin:$PATH
swag init -g Server.go --output docs/app


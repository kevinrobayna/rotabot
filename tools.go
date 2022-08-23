//go:build tools
// +build tools

package tools

import (
	_ "github.com/cespare/reflex"
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "gotest.tools/gotestsum"
)
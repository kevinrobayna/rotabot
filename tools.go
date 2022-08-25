//go:build tools
// +build tools

package tools

import (
	_ "github.com/cespare/reflex"
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "goa.design/goa/v3/cmd/goa"
	_ "gotest.tools/gotestsum"
)

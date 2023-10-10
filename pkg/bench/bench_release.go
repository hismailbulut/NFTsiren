//go:build !debug

package bench

import "io"

func Begin() func(name ...string) { return func(name ...string) {} }

func PrintResults(out io.Writer) {}

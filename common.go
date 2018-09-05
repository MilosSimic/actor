package main

import (
	"strings"
)

func join(parts ...string) string {
	return strings.Join(parts, "/")
}

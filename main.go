package main

import (
	"fmt"
	"pastezero-cli/cmd"
)

var (
	version   = "dev"
	buildTime = "unknown"
)

func main() {
	cmd.VersionInfo = fmt.Sprintf("PasteZero CLI %s (%s)", version, buildTime)
	cmd.Execute()
}
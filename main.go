package main

import (
	"os"

	"github.com/crazyuploader/zfs-dash/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

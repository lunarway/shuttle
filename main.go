package main

import (
	"os"

	"github.com/lunarway/shuttle/cmd"
)

func main() {
	cmd.Execute(os.Stdout, os.Stderr)
}

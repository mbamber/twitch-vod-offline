package main

import (
	"github.com/mbamber/tvo/cmd"
	"github.com/xfxdev/xlog"
)

func main() {
	setupLogging()
	cmd.Execute()
}

func setupLogging() {
	xlog.SetLevel(xlog.DebugLevel)
}

package main

import (
	"github.com/mbamber/tvo/event"
	"github.com/mbamber/tvo/twitch"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/xfxdev/xlog"
)

func main() {
	setupLogging()
	lambda.Start(pollingHandler)
}

func setupLogging() {
	xlog.SetLevel(xlog.DebugLevel)
}

func pollingHandler(evt event.Event) (err error) {
	xlog.Debugf("Polling lambda handler invoked")
	vods, err := twitch.ListTwitchVODs(evt.ChannelNames)

	if evt.NotifyURL == "" {
		xlog.Infof("Notify URL was not supplied so not notifying anywhere of the result")
		return
	}

	xlog.Infof("Would post requests here for %d vods", len(vods))

	return
}

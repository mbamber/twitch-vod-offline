package twitch

import (
	"github.com/xfxdev/xlog"
)

// ListTwitchVODs lists details about each VOD uploaded from each of the given channel names
func ListTwitchVODs(channelNames []string) (vods []VOD, err error) {
	xlog.Debugf("Getting VODs for channels %+v", channelNames)

	conn := NewConnection("sa2qabz4cm25y5ep8dnib2lusi0xuc", "ukj89zfbpne62biwt599g0au6cieju")
	channels, err := conn.GetChannels(channelNames)
	if err != nil {
		return
	}

	xlog.Debugf("Got %d channels from %d channel names", len(channels), len(channelNames))

	// Get the channel IDs
	ids := make([]string, 0)
	for _, c := range channels {
		ids = append(ids, c.ID)
	}
	xlog.Infof("Got channel IDs as: %+v", ids)

	vods, err = conn.GetVODs(ids)
	if err != nil {
		return
	}

	xlog.Infof("Got %d vods", len(vods))
	return
}

package twitch

import "fmt"

func ListTwitchVODs(channelNames []string) (err error) {
	conn := Connection{ClientID: "sa2qabz4cm25y5ep8dnib2lusi0xuc"}

	channels, err := conn.GetChannels(channelNames)
	if err != nil {
		return
	}

	fmt.Println(channels)

	return nil
}

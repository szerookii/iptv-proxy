package main

import (
	"fmt"
	"github.com/szerookii/iptv-proxy/iptv/stb"
)

func main() {
	c, err := stb.NewClient("", "")
	if err != nil {
		panic(err)
	}

	mainInfo, err := c.MainInfo()
	if err != nil {
		panic(err)
	}

	fmt.Printf("MAC: %s\nSubscription ends on: %s\n", mainInfo.Mac, mainInfo.Phone)

	link, err := c.CreateLink("ffmpeg http:\\/\\/localhost\\/ch\\/403102_")
	if err != nil {
		panic(err)
	}

	fmt.Printf("Link: %s\n", link)
}

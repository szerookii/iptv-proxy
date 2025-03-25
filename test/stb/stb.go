package main

import (
	"fmt"
	"github.com/szerookii/iptv-proxy/iptv/stb"
)

func main() {
	c, err := stb.NewClient("http://ns31617126.ip-162-19-37.eu:8080/c", "00:1A:79:A8:D6:73")
	if err != nil {
		panic(err)
	}

	username, pass, err := c.ConvertToXtream()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Username: %s\nPassword %s\n", username, pass)
}

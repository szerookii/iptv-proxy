package xtream

import (
	"fmt"
	"github.com/gofiber/fiber/v3"
	"github.com/szerookii/iptv-proxy/config"
	"github.com/szerookii/iptv-proxy/iptv/stb"
	"github.com/szerookii/iptv-proxy/iptv/xtream"
	"github.com/szerookii/iptv-proxy/utils"
	"strconv"
	"strings"
	"time"
)

func PlayerAPI(c fiber.Ctx) error {
	username := c.Query("username")
	password := c.Query("password")
	action := c.Query("action")

	if username != config.Get().Xtream.Username || password != config.Get().Xtream.Password {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	switch action {
	case "":
		var accountInfo = new(xtream.AccountInfo)

		accountInfo.UserInfo.Username = username
		accountInfo.UserInfo.Password = password
		accountInfo.UserInfo.Message = ""
		accountInfo.UserInfo.Auth = 1

		switch remote := config.Get().Remote.Data.(type) {
		case *config.StbRemote:
			stb, err := stb.NewClient(remote.URL, remote.MacAddress)
			if err != nil {
				fmt.Println(err)
				return c.SendStatus(fiber.StatusInternalServerError)
			}

			mainInfo, err := stb.MainInfo()
			if err != nil {
				fmt.Println(err)
				return c.SendStatus(fiber.StatusInternalServerError)
			}

			parsedDate, err := time.Parse("January 2, 2006 15:04", mainInfo.Phone)
			if err != nil {
				parsedDate = time.Now().Add(24 * time.Hour)
			}

			if parsedDate.Before(time.Now()) {
				accountInfo.UserInfo.Status = "Expired"
			} else {
				accountInfo.UserInfo.Status = "Active"
			}

			accountInfo.UserInfo.ExpDate = fmt.Sprintf("%d", parsedDate.Unix())
			accountInfo.UserInfo.IsTrial = "0"
			accountInfo.UserInfo.ActiveCons = "0"
			accountInfo.UserInfo.CreatedAt = fmt.Sprintf("%d", time.Now().Add(-72*time.Hour).Unix())
			accountInfo.UserInfo.MaxConnections = "10" // TODO: idk do we need this
			accountInfo.UserInfo.AllowedOutputFormats = []string{"ts"}

			accountInfo.ServerInfo.Url = c.BaseURL()
			accountInfo.ServerInfo.Port = fmt.Sprintf("%d", config.Get().Port)
			accountInfo.ServerInfo.HttpsPort = fmt.Sprintf("%d", config.Get().Port)
			accountInfo.ServerInfo.RtmpPort = ""
			accountInfo.ServerInfo.Timezone = "Europe/Kiev"
			accountInfo.ServerInfo.TimestampNow = int(time.Now().Unix())
			accountInfo.ServerInfo.TimeNow = time.Now().Format("2006-01-02 15:04:05")
			accountInfo.ServerInfo.Process = true
			break

			// TODO: Implement XtreamRemote
		}

		return c.JSON(accountInfo)

	case "get_live_categories":
		var liveCategories []*xtream.LiveCategory

		switch remote := config.Get().Remote.Data.(type) {
		case *config.StbRemote:
			stb, err := stb.NewClient(remote.URL, remote.MacAddress)
			if err != nil {
				return c.SendStatus(fiber.StatusInternalServerError)
			}

			genres, err := stb.Genres()
			if err != nil {
				return c.SendStatus(fiber.StatusInternalServerError)
			}

			for _, genre := range genres {
				categoryId, err := strconv.Atoi(genre.Id)
				if err != nil {
					continue
				}

				liveCategories = append(liveCategories, &xtream.LiveCategory{
					CategoryId:   categoryId,
					CategoryName: genre.Title,
					ParentId:     0, // TODO: try to find out what is this
				})
			}

			// TODO: Implement XtreamRemote
		}

		return c.JSON(liveCategories)
	case "get_live_streams":
		var liveStreams []*xtream.LiveStream

		switch remote := config.Get().Remote.Data.(type) {
		case *config.StbRemote:
			stbC, err := stb.NewClient(remote.URL, remote.MacAddress)
			if err != nil {
				return c.SendStatus(fiber.StatusInternalServerError)
			}

			var channels []*stb.ItvChannel
			if c.Query("category_id") == "" || strings.ToLower(c.Query("category_id")) == "all" {
				channels, err = stbC.Channels()
			} else {
				channels, err = stbC.GenreChannels(c.Query("category_id"))
			}

			for _, channel := range channels {
				streamId, err := strconv.Atoi(channel.Id)
				if err != nil {
					continue
				}

				liveStreams = append(liveStreams, &xtream.LiveStream{
					Num:               streamId,
					Name:              channel.Name,
					StreamType:        "live",
					StreamId:          streamId,
					StreamIcon:        channel.Logo,
					EpgChannelId:      "",
					Added:             "0",
					IsAdult:           "0",
					CategoryId:        channel.TvGenreId,
					CustomSid:         "0",
					TvArchive:         0,
					DirectSource:      "",
					TvArchiveDuration: "0",
				})
			}

			return c.JSON(liveStreams)

			// TODO: Implement XtreamRemote
		}

	case "get_short_epg":
		// TODO: Implement
		return utils.SendJSON(c, 200, fiber.Map{})
	}

	return c.SendStatus(fiber.StatusNotFound)
}

func Live(c fiber.Ctx) error {
	username := c.Params("username")
	password := c.Params("password")
	streamId := c.Params("stream")

	if username != config.Get().Xtream.Username || password != config.Get().Xtream.Password {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	switch remote := config.Get().Remote.Data.(type) {
	case *config.StbRemote:
		channelId := strings.Split(streamId, ".")[0]
		client, err := stb.NewClient(remote.URL, remote.MacAddress)
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		// TODO: more optimized way to get channel by id
		channels, err := client.Channels()
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		for _, channel := range channels {
			if channel.Id == channelId {
				link, err := client.CreateLink(channel.Cmd)
				if err != nil {
					return c.SendStatus(fiber.StatusInternalServerError)
				}
				return c.Redirect().To(link)
			}
		}

		return c.SendStatus(fiber.StatusNotFound)

		// TODO: Implement XtreamRemote
	}

	return c.SendStatus(fiber.StatusNotFound)
}

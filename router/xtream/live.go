package xtream

import (
	"fmt"
	"github.com/gofiber/fiber/v3"
	"github.com/szerookii/iptv-proxy/config"
	"github.com/szerookii/iptv-proxy/iptv/stb"
	"github.com/szerookii/iptv-proxy/utils"
	"strconv"
	"strings"
)

func Media(c fiber.Ctx) error {
	t := c.Params("type")
	username := c.Params("username")
	password := c.Params("password")

	if username != config.Get().Xtream.Username || password != config.Get().Xtream.Password {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	switch t {
	case "live":
		return Live(c)

	case "movie":
		return Movie(c)
	}

	return c.SendStatus(fiber.StatusNotFound)
}

func Live(c fiber.Ctx) error {
	streamId := c.Params("stream")

	switch remote := config.Get().Remote.Data.(type) {
	case *config.StbRemote:
		channelIdStr := strings.Split(streamId, ".")[0]
		stbC, err := stb.NewClient(remote.URL, remote.MacAddress)
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		channelIdMerged, err := strconv.ParseInt(channelIdStr, 10, 64)
		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		genreId, channelId := utils.SplitNumbers(channelIdMerged)
		channels, err := stbC.GenreChannels(fmt.Sprintf("%d", genreId))
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		for _, channel := range channels {
			if channel.Id == fmt.Sprintf("%d", channelId) {
				link, err := stbC.CreateLink("itv", channel.Cmd)
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

func Movie(c fiber.Ctx) error {
	streamId := c.Params("stream")

	switch remote := config.Get().Remote.Data.(type) {
	case *config.StbRemote:
		streamIdStr := strings.Split(streamId, ".")[0]
		stbC, err := stb.NewClient(remote.URL, remote.MacAddress)
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		streamIdMerged, err := strconv.ParseInt(streamIdStr, 10, 64)
		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		categoryId, vodId := utils.SplitNumbers(streamIdMerged)
		vods, err := stbC.CategoryVods(fmt.Sprintf("%d", categoryId))
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		vodIdStr := fmt.Sprintf("%d", vodId)

		for _, vod := range vods {
			if vod.Id == vodIdStr {
				link, err := stbC.CreateLink("vod", vod.Cmd)
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

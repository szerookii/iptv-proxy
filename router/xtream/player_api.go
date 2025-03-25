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
			stbC, err := stb.NewClient(remote.URL, remote.MacAddress)
			if err != nil {
				return c.SendStatus(fiber.StatusInternalServerError)
			}

			mainInfo, err := stbC.MainInfo()
			if err != nil {
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
		liveCategories := make([]*xtream.LiveCategory, 0)

		switch remote := config.Get().Remote.Data.(type) {
		case *config.StbRemote:
			stbC, err := stb.NewClient(remote.URL, remote.MacAddress)
			if err != nil {
				return c.SendStatus(fiber.StatusInternalServerError)
			}

			genres, err := stbC.ItvGenres()
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

		case *config.XtreamRemote:

		}

		return c.JSON(liveCategories)
	case "get_live_streams":
		liveStreams := make([]*xtream.LiveStream, 0)

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
				channelId, err := strconv.Atoi(channel.Id)
				if err != nil {
					continue
				}

				genreId, err := strconv.Atoi(channel.TvGenreId)
				if err != nil {
					continue
				}

				liveStreams = append(liveStreams, &xtream.LiveStream{
					Num:               channelId,
					Name:              channel.Name,
					StreamType:        "live",
					StreamId:          int(utils.MergeNumbers(int32(genreId), int32(channelId))),
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

	case "get_vod_categories":
		vodCategories := make([]*xtream.VodCategory, 0)

		switch remote := config.Get().Remote.Data.(type) {
		case *config.StbRemote:
			stbC, err := stb.NewClient(remote.URL, remote.MacAddress)
			if err != nil {
				return c.SendStatus(fiber.StatusInternalServerError)
			}

			categories, err := stbC.VodCategories()
			if err != nil {
				return c.SendStatus(fiber.StatusInternalServerError)
			}

			for _, genre := range categories {
				categoryId, err := strconv.Atoi(genre.Id)
				if err != nil {
					continue
				}

				vodCategories = append(vodCategories, &xtream.VodCategory{
					CategoryId:   categoryId,
					CategoryName: genre.Title,
					ParentId:     0, // TODO: try to find out what is this
				})
			}

			// TODO: Implement XtreamRemote
		}

		return c.JSON(vodCategories)

	case "get_vod_streams":
		vodStreams := make([]*xtream.VodStream, 0)

		switch remote := config.Get().Remote.Data.(type) {
		case *config.StbRemote:
			stbC, err := stb.NewClient(remote.URL, remote.MacAddress)
			if err != nil {
				return c.SendStatus(fiber.StatusInternalServerError)
			}

			var vods []*stb.Vod
			if c.Query("category_id") == "" || strings.ToLower(c.Query("category_id")) == "all" {
				// TODO: implement all vods
			} else {
				vods, err = stbC.CategoryVods(c.Query("category_id"))
			}

			fmt.Println(len(vods))

			for _, vod := range vods {
				vodId, err := strconv.Atoi(vod.Id)
				if err != nil {
					continue
				}

				vodCategoryId, err := strconv.Atoi(vod.CategoryId)
				if err != nil {
					continue
				}

				rating, err := strconv.ParseFloat(vod.RatingImdb, 64)
				if err != nil {
					rating = 0.0
				}

				vodStreams = append(vodStreams, &xtream.VodStream{
					Num:                vodId,
					Name:               vod.Name,
					StreamType:         "movie",
					StreamId:           int(utils.MergeNumbers(int32(vodCategoryId), int32(vodId))),
					StreamIcon:         vod.ScreenshotUri,
					Rating:             rating,
					Rating5Based:       rating,
					Added:              vod.Added,
					IsAdult:            "0", // TODO: parse age and check if it's adult
					CategoryId:         vod.CategoryId,
					ContainerExtension: "mp4",
					CustomSid:          "",
					DirectSource:       "",
				})
			}

			return c.JSON(vodStreams)

			// TODO: Implement XtreamRemote
		}

	case "get_vod_info":
		vodInfo := new(xtream.VodInfo)

		switch remote := config.Get().Remote.Data.(type) {
		case *config.StbRemote:
			stbC, err := stb.NewClient(remote.URL, remote.MacAddress)
			if err != nil {
				return c.SendStatus(fiber.StatusInternalServerError)
			}

			mergedVodId, err := strconv.ParseInt(c.Query("vod_id"), 10, 64)
			if err != nil {
				return c.SendStatus(fiber.StatusInternalServerError)
			}

			vodCategory, vodId := utils.SplitNumbers(mergedVodId)
			vods, err := stbC.CategoryVods(fmt.Sprintf("%d", vodCategory))
			vodIdStr := fmt.Sprintf("%d", vodId)

			for _, vod := range vods {
				if vod.Id == vodIdStr {
					vodInfo.Info.Name = vod.Name
					vodInfo.Info.OName = vod.Name
					vodInfo.Info.MovieImage = vod.ScreenshotUri
					vodInfo.Info.Releasedate = vod.Year
					vodInfo.Info.EpisodeRunTime = "0"
					vodInfo.Info.YoutubeTrailer = ""
					vodInfo.Info.Director = vod.Director
					vodInfo.Info.Actors = vod.Actors
					vodInfo.Info.Cast = vod.Actors
					vodInfo.Info.Description = vod.Description
					vodInfo.Info.Plot = vod.Description
					vodInfo.Info.Age = vod.Age
					vodInfo.Info.MpaaRating = "0"
					vodInfo.Info.RatingCountKinopoisk = 0
					vodInfo.Info.Country = ""
					vodInfo.Info.Genre = vod.GenresStr
					vodInfo.Info.BackdropPath = []string{vod.ScreenshotUri}
					// TODO: fix duration
					vodInfo.Info.DurationSecs = 0
					vodInfo.Info.Duration = "0"
					vodInfo.Info.Video = []interface{}{}
					vodInfo.Info.Audio = []interface{}{}
					vodInfo.Info.Bitrate = 0
					vodInfo.Info.Rating = vod.RatingImdb
					vodInfo.MovieData.StreamId = int(mergedVodId)
					vodInfo.MovieData.Name = vod.Name
					vodInfo.MovieData.Added = vod.Added
					vodInfo.MovieData.CategoryId = vod.CategoryId
					vodInfo.MovieData.ContainerExtension = "mp4"
					vodInfo.MovieData.CustomSid = ""
					break
				}
			}

			return c.JSON(vodInfo)

			// TODO: Implement XtreamRemote
		}
	}

	return c.SendStatus(fiber.StatusNotFound)
}

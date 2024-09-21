package middleware

import (
	"fmt"
	"github.com/gofiber/fiber/v3"
	"github.com/phuslu/log"
	"github.com/szerookii/iptv-proxy/config"
)

func LogsMiddleware(c fiber.Ctx) error {
	if config.Get().EnableLogs {
		queryString := ""

		for k, v := range c.Queries() {
			queryString += fmt.Sprintf("%s=%s&", k, v)
		}

		if len(queryString) > 0 {
			queryString = "?" + queryString[:len(queryString)-1]
		}

		log.Info().Msgf("=> %s%s from %s", c.Path(), queryString, c.IP())
	}

	return c.Next()
}

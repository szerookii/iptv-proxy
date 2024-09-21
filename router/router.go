package router

import (
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v3"
	"github.com/phuslu/log"
	"github.com/szerookii/iptv-proxy/router/middleware"
	"github.com/szerookii/iptv-proxy/router/xtream"
)

func Init() *fiber.App {
	log.Info().Msgf("Initializing router...")

	app := fiber.New(fiber.Config{
		JSONEncoder: sonic.Marshal,
		JSONDecoder: sonic.Unmarshal,
	})

	app.Use(middleware.LogsMiddleware)

	// Xtream Codes API
	app.Get("/player_api.php", xtream.PlayerAPI)
	app.Get("/live/:username/:password/:stream", xtream.Live)

	return app
}

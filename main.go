package main

import (
	"fmt"
	"github.com/gofiber/fiber/v3"
	"github.com/phuslu/log"
	"github.com/szerookii/iptv-proxy/config"
	"github.com/szerookii/iptv-proxy/router"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	if log.IsTerminal(os.Stderr.Fd()) {
		log.DefaultLogger = log.Logger{
			TimeFormat: "15:04:05",
			Caller:     1,
			Writer: &log.ConsoleWriter{
				ColorOutput:    true,
				QuoteString:    true,
				EndWithMessage: true,
			},
		}
	}

	log.Info().Msgf("Starting IPTV proxy...")

	log.Info().Msg("Loading config...")
	if config.Get().Remote.Data == nil {
		log.Warn().Msg("No remote data found in config.json")
		return
	}

	r := router.Init()

	go r.Listen(fmt.Sprintf(":%d", config.Get().Port), fiber.ListenConfig{
		DisableStartupMessage: true,
	})

	log.Info().Msgf("Started router on port %d", config.Get().Port)

	sChan := make(chan os.Signal, 1)
	defer close(sChan)
	signal.Notify(sChan, syscall.SIGINT, syscall.SIGTERM)

	log.Info().Msg("IPTV proxy started successfully!")

	<-sChan

	if err := r.Shutdown(); err != nil {
		log.Error().Err(err).Msg("Failed to shutdown.")
	}
}

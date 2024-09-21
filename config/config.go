package config

import (
	"encoding/json"
	"github.com/phuslu/log"
	"os"
)

type Config struct {
	Port       int     `json:"port"`
	EnableLogs bool    `json:"enable_logs"`
	Remote     *Remote `json:"remote"`
	Xtream     *Xtream `json:"local_xtream"`
}

type Xtream struct {
	Enabled  bool   `json:"enabled"`
	Username string `json:"username"`
	Password string `json:"password"`
}

var config *Config

func init() {
	if _, err := os.Stat("config.json"); os.IsNotExist(err) {
		c := &Config{
			Port:   3000,
			Remote: new(Remote),
			Xtream: new(Xtream),
		}

		bytes, err := json.MarshalIndent(&c, "", "    ")
		if err != nil {
			log.Panic().Err(err).Msg("Failed to marshal config")
		}

		err = os.WriteFile("config.json", bytes, 0644)
		if err != nil {
			log.Panic().Err(err).Msg("Failed to write config.json")
		}
	}

	data, err := os.ReadFile("config.json")

	if err != nil {
		log.Panic().Err(err).Msg("Failed to read config.json")
	}

	c := new(Config)
	_ = json.Unmarshal(data, &c)

	config = c
}

func Get() *Config {
	return config
}

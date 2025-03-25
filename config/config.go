package config

import (
	"github.com/joho/godotenv"
	"os"
	"strconv"
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
	_ = godotenv.Load()
	if _, err := os.Stat("config.json"); os.IsNotExist(err) {
		c := &Config{
			Port:   3000,
			Remote: new(Remote),
			Xtream: new(Xtream),
		}

		if port, err := strconv.Atoi(os.Getenv("PORT")); err == nil {
			c.Port = port
		}

		if enableLogs, err := strconv.ParseBool(os.Getenv("ENABLE_LOGS")); err == nil {
			c.EnableLogs = enableLogs
		}

		remoteType := os.Getenv("REMOTE_TYPE")
		switch remoteType {
		case "stb":
			c.Remote.Data = new(StbRemote)
			c.Remote.Type = RemoteTypeStb
			_ = c.Remote.Data.UnmarshalJSON([]byte(`{"url":"` + os.Getenv("REMOTE_URL") + `","mac_address":"` + os.Getenv("REMOTE_MAC_ADDRESS") + `"}`))

		case "xtream":
			c.Remote.Data = new(XtreamRemote)
			c.Remote.Type = RemoteTypeXtream
			_ = c.Remote.Data.UnmarshalJSON([]byte(`{"url":"` + os.Getenv("REMOTE_URL") + `","username":"` + os.Getenv("REMOTE_USERNAME") + `","password":"` + os.Getenv("REMOTE_PASSWORD") + `"}`))
		}

		if enabled, err := strconv.ParseBool(os.Getenv("LOCAL_XTREAM_ENABLED")); err == nil {
			c.Xtream.Enabled = enabled
			c.Xtream.Username = os.Getenv("LOCAL_XTREAM_USERNAME")
			c.Xtream.Password = os.Getenv("LOCAL_XTREAM_PASSWORD")
		}

		config = c
	}
}

func Get() *Config {
	return config
}

package xtream

type VodCategory struct {
	CategoryId   int    `json:"category_id"`
	CategoryName string `json:"category_name"`
	ParentId     int    `json:"parent_id"`
}

type VodStream struct {
	Num                int     `json:"num"`
	Name               string  `json:"name"`
	StreamType         string  `json:"stream_type"`
	StreamId           int     `json:"stream_id"`
	StreamIcon         string  `json:"stream_icon"`
	Rating             float64 `json:"rating"`
	Rating5Based       float64 `json:"rating_5based"`
	Added              string  `json:"added"`
	IsAdult            string  `json:"is_adult"`
	CategoryId         string  `json:"category_id"`
	ContainerExtension string  `json:"container_extension"`
	CustomSid          string  `json:"custom_sid"`
	DirectSource       string  `json:"direct_source"`
}

type VodInfo struct {
	Info struct {
		TmdbUrl              string        `json:"tmdb_url"`
		TmdbId               string        `json:"tmdb_id"`
		Name                 string        `json:"name"`
		OName                string        `json:"o_name"`
		CoverBig             string        `json:"cover_big"`
		MovieImage           string        `json:"movie_image"`
		Releasedate          string        `json:"releasedate"`
		EpisodeRunTime       string        `json:"episode_run_time"`
		YoutubeTrailer       string        `json:"youtube_trailer"`
		Director             string        `json:"director"`
		Actors               string        `json:"actors"`
		Cast                 string        `json:"cast"`
		Description          string        `json:"description"`
		Plot                 string        `json:"plot"`
		Age                  string        `json:"age"`
		MpaaRating           string        `json:"mpaa_rating"`
		RatingCountKinopoisk int           `json:"rating_count_kinopoisk"`
		Country              string        `json:"country"`
		Genre                string        `json:"genre"`
		BackdropPath         []string      `json:"backdrop_path"`
		DurationSecs         int           `json:"duration_secs"`
		Duration             string        `json:"duration"`
		Video                []interface{} `json:"video"`
		Audio                []interface{} `json:"audio"`
		Bitrate              int           `json:"bitrate"`
		Rating               string        `json:"rating"`
	} `json:"info"`
	MovieData struct {
		StreamId           int    `json:"stream_id"`
		Name               string `json:"name"`
		Added              string `json:"added"`
		CategoryId         string `json:"category_id"`
		ContainerExtension string `json:"container_extension"`
		CustomSid          string `json:"custom_sid"`
		DirectSource       string `json:"direct_source"`
	} `json:"movie_data"`
}

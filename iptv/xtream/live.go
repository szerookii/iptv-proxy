package xtream

type LiveCategory struct {
	CategoryId   int    `json:"category_id"`
	CategoryName string `json:"category_name"`
	ParentId     int    `json:"parent_id"`
}

type LiveStream struct {
	Num               int    `json:"num"`
	Name              string `json:"name"`
	StreamType        string `json:"stream_type"`
	StreamId          int    `json:"stream_id"`
	StreamIcon        string `json:"stream_icon"`
	EpgChannelId      string `json:"epg_channel_id"`
	Added             string `json:"added"`
	IsAdult           string `json:"is_adult"`
	CategoryId        string `json:"category_id"`
	CustomSid         string `json:"custom_sid"`
	TvArchive         int    `json:"tv_archive"`
	DirectSource      string `json:"direct_source"`
	TvArchiveDuration string `json:"tv_archive_duration"`
}

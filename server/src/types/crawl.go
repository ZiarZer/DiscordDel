package types

type CrawlingInfo struct {
	ChannelId    string `json:"channel_id"`
	OldestReadId string `json:"oldest_read_id"`
	NewestReadId string `json:"newest_read_id"`
	ReachedTop   bool   `json:"reached_top"`
}

package types

type CrawlingInfo struct {
	ChannelId           string `json:"channel_id"`
	OldestReadMessageId string `json:"oldest_read_message_id"`
	NewestReadMessageId string `json:"newest_read_message_id"`
	ReachedTop          bool   `json:"reached_top"`
}

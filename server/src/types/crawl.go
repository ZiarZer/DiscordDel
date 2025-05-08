package types

type CrawlingInfo struct {
	ChannelId    Snowflake `json:"channel_id"`
	OldestReadId Snowflake `json:"oldest_read_id"`
	NewestReadId Snowflake `json:"newest_read_id"`
	ReachedTop   bool      `json:"reached_top"`
}

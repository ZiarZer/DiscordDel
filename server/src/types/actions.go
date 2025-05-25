package types

import (
	"time"
)

type CrawlingInfo struct {
	ChannelId    Snowflake `json:"channel_id"`
	OldestReadId Snowflake `json:"oldest_read_id"`
	NewestReadId Snowflake `json:"newest_read_id"`
	ReachedTop   bool      `json:"reached_top"`
}

type ActionType string
type Scope string
type Action struct {
	Id          *int64
	Type        *ActionType
	Scope       *Scope
	TargetId    *Snowflake
	Description string
	StartTime   time.Time
	LogFunc     func(message string, logLevel *LogLevel)
	LogEndTime  bool
}

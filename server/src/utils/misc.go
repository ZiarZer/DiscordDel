package utils

import (
	"fmt"
	"strconv"

	"github.com/ZiarZer/DiscordDel/types"
)

func MakePointer[T any](value T) *T {
	return &value
}

func GetTimestampFromSnowflake(snowflakeId types.Snowflake) int {
	intSnowflake, err := strconv.Atoi(string(snowflakeId))
	if err != nil {
		return 0
	}
	return intSnowflake >> 22
}

func FormatDuration(seconds int64) string {
	h := seconds / 3600
	m := seconds / 60 % 60
	s := seconds % 60
	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}

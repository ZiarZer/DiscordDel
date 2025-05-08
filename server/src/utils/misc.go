package utils

import (
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

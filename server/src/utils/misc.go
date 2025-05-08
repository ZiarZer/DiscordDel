package utils

import "strconv"

func MakePointer[T any](value T) *T {
	return &value
}

func GetTimestampFromSnowflake(snowflakeId string) int {
	intSnowflake, err := strconv.Atoi(snowflakeId)
	if err != nil {
		return 0
	}
	return intSnowflake >> 22
}

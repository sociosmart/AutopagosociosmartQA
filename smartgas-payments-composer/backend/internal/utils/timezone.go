package utils

import (
	"time"
)

func InitTimezone(timezone string) {
	loc, err := time.LoadLocation(timezone)

	if err != nil {
		panic(err)
	}

	time.Local = loc
}

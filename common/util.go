package common

import (
	"time"
	"unicode"

	"github.com/google/uuid"
)

func GenerateUUID() string {
	id := uuid.New()

	return id.String()
}

var countryTz = map[string]string{
	"Jakarta":  "Asia/Jakarta",
	"Jayapura": "Asia/Jayapura",
	"Makassar": "Asia/Makassar",
}

func TimeIn(name string) time.Time {
	loc, err := time.LoadLocation(countryTz[name])
	if err != nil {
		panic(err)
	}
	return time.Now().In(loc)
}

func TimeTo(t time.Time, name string) time.Time {
	loc, err := time.LoadLocation(countryTz[name])
	if err != nil {
		panic(err)
	}
	return t.In(loc)
}

func TruncateToLength(s string, max int) string {
	lastSpaceIx := -1
	len := 0
	for i, r := range s {
		if unicode.IsSpace(r) {
			lastSpaceIx = i
		}
		len++
		if len >= max {
			if lastSpaceIx != -1 {
				return s[:lastSpaceIx] + "..."
			}

			break
		}
	}

	return s
}

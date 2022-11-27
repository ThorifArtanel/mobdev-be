package common

import (
	"time"

	guuid "github.com/google/uuid"
)

func GenerateUUID() string {
	id := guuid.New()

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

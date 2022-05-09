package util

import (
	"strings"
	"time"
)

func ParseTime(dateTime string) string {

	dateTimeN := strings.SplitN(dateTime, "+", -1)[0]

	// set location America/Sao_Paulo

	loc, _ := time.LoadLocation("America/Sao_Paulo")

	time, _ := time.Parse("2006-01-02T15:04:05", dateTimeN)

	return time.In(loc).Format("2006-01-02 15:04:05")
}

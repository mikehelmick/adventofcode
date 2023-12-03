package straid

import (
	"strconv"
	"strings"
)

func Field(s string, sep string, pos int) string {
	parts := strings.Split(s, sep)
	return parts[pos]
}

func AsInt(s string) int64 {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		panic(err)
	}
	return i
}

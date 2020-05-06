package controllers

import (
	"strconv"
	"strings"
)

const defaultMaxResultCount = 30

func convertStrToArrInt64(str string) []int64 {
	idsStr := strings.Split(strings.TrimSpace(str), ",")
	var ids []int64
	for _, idstr := range idsStr {
		id, _ := strconv.ParseInt(idstr, 10, 64)
		if id != 0 {
			ids = append(ids, id)
		}
	}
	return ids
}

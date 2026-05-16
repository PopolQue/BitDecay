package game

import (
	"fmt"
	"strconv"
	"strings"
)

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func ParseTile(key string) (int, int, error) {
	parts := strings.Split(key, ",")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid tile key")
	}
	q, err1 := strconv.Atoi(parts[0])
	r, err2 := strconv.Atoi(parts[1])
	if err1 != nil || err2 != nil {
		return 0, 0, fmt.Errorf("invalid tile coordinates")
	}
	return q, r, nil
}

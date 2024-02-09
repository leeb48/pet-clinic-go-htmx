package handlers

import (
	"strconv"
	"strings"
)

func atoiWithDefault(numStr string, num int) int {
	intVal := num

	if strings.TrimSpace(numStr) != "" {
		intVal, _ = strconv.Atoi(numStr)
	}

	return intVal
}

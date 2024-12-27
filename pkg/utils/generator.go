package utils

import (
	"fmt"
	"math/rand/v2"
	"strconv"
	"time"
)

func GenerateFolderStructFromDate() string {
	date := time.Now()

	// folder struct
	year := date.Year()
	month := date.Month()
	day := date.Day()

	return fmt.Sprintf("%d/%d/%d/", year, month, day)
}

func GenerateRandomUint(limit int) string {
	name := ""
	for i := 0; i < limit; i++ {
		name += strconv.Itoa(rand.IntN(9))
	}
	return name
}

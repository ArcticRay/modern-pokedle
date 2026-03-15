package game

import (
	"crypto/sha256"
	"time"
)

func DailyPokemonID(maxID int) int {
	date := time.Now().Format("2006-01-02")

	hash := sha256.Sum256([]byte(date))

	num := int(hash[0])<<24 | int(hash[1])<<16 | int(hash[2])<<8 | int(hash[3])

	return (num % maxID) + 1
}

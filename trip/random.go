package trip

import (
	"math/rand"
)

func RandInt(min int, max int) int {
	return min + rand.Intn(max-min)
}
func RandInt16(min int, max int) uint16 {
	return uint16(min + rand.Intn(max-min))
}
func RandFloat(min float64, max float64) float64 {
	return min + rand.Float64()*(max-min)
}
func Chance(n int) bool {
	return RandInt(0, n) == 0
}

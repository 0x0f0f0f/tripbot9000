package trip

import (
	"math/rand"
)

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}
func randInt16(min int, max int) uint16 {
	return uint16(min + rand.Intn(max-min))
}
func randFloat(min float64, max float64) float64 {
	return min + rand.Float64()*(max-min)
}
func chance(n int) bool {
	return randInt(0, n) == 0
}

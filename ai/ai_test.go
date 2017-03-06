package ai

import (
	"testing"
	"time"

	"github.com/lukechampine/tenten/game"
)

func BenchmarkMove(b *testing.B) {
	g := game.New()
	start := time.Now()
	n := 0
	for Move(&g, g.NextBag()) {
		n++
	}
	b.Log("Score:", g.Score())
	b.Log("Bags:", n)
	b.Log("Time/move:", time.Since(start)/time.Duration(n))
}

package ai

import (
	"testing"
	"time"

	"github.com/lukechampine/tenten/game"
)

func BenchmarkGame(b *testing.B) {
	g := game.New()
	n := 0
	start := time.Now()
lost:
	for {
		moves := BestMoves(g.Board(), g.NextBag())
		for _, m := range moves {
			if !g.Place(m.Piece, m.X, m.Y) {
				break lost
			}
		}
		n++
	}
	b.Log("Score:", g.Score())
	b.Log("Bags:", n)
	b.Log("Time/move:", time.Since(start)/time.Duration(n))
}

package ai

import (
	"testing"
	"time"

	"github.com/fatih/color"
	"github.com/lukechampine/tenten/game"
)

func printBoard(b *game.Board) {
	for _, row := range b {
		for _, c := range row {
			switch c {
			case game.Empty:
				print("□")
			case game.Red:
				print(color.RedString("■"))
			case game.Pink:
				print(color.New(color.FgMagenta, color.Bold).Sprint("■"))
			case game.Orange:
				print(color.New(color.FgRed, color.Bold).Sprint("■"))
			case game.Yellow:
				print(color.YellowString("■"))
			case game.Green:
				print(color.GreenString("■"))
			case game.Teal:
				print(color.New(color.FgCyan, color.Faint).Sprint("■"))
			case game.Cyan:
				print(color.CyanString("■"))
			case game.Blue:
				print(color.BlueString("■"))
			case game.Purple:
				print(color.New(color.FgHiMagenta, color.Faint).Sprint("■"))
			}
			print(" ")
		}
		println()
	}
}

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

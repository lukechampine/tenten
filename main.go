package main

import (
	"fmt"
	"time"

	"github.com/fatih/color"

	"github.com/lukechampine/tenten/ai"
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

func main() {
	g := game.New()
	n := 0
	start := time.Now()
	for ai.Move(&g, g.NextBag()) {
		n++
	}
	elapsed := time.Since(start)
	printBoard(g.Board())
	fmt.Printf("\nFinal Score: %v\nPlayed %v bags in %v\nAverage move time: %v/bag\n", g.Score(), n, elapsed, elapsed/time.Duration(n))
}

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
				color.New(color.FgRed).Print("■")
			case game.Pink:
				color.New(color.FgMagenta, color.Bold).Print("■")
			case game.Orange:
				color.New(color.FgRed, color.Bold).Print("■")
			case game.Yellow:
				color.New(color.FgYellow).Print("■")
			case game.Green:
				color.New(color.FgGreen).Print("■")
			case game.Teal:
				color.New(color.FgCyan, color.Faint).Print("■")
			case game.Cyan:
				color.New(color.FgCyan).Print("■")
			case game.Blue:
				color.New(color.FgBlue).Print("■")
			case game.Purple:
				color.New(color.FgHiMagenta, color.Faint).Print("■")
			}
			print(" ")
		}
		println()
	}
}

func printMoves(moves []ai.Move) {
	var grid [5][18]game.Color
	for i, m := range moves {
		for _, d := range m.Piece.Dots() {
			x := d.X + (len(moves)-i-1)*6
			grid[d.Y][x] = m.Piece.Color()
		}
	}

	for _, row := range grid {
		for _, c := range row {
			switch c {
			case game.Empty:
				print(" ")
			case game.Red:
				color.New(color.FgRed).Print("■")
			case game.Pink:
				color.New(color.FgMagenta, color.Bold).Print("■")
			case game.Orange:
				color.New(color.FgRed, color.Bold).Print("■")
			case game.Yellow:
				color.New(color.FgYellow).Print("■")
			case game.Green:
				color.New(color.FgGreen).Print("■")
			case game.Teal:
				color.New(color.FgCyan, color.Faint).Print("■")
			case game.Cyan:
				color.New(color.FgCyan).Print("■")
			case game.Blue:
				color.New(color.FgBlue).Print("■")
			case game.Purple:
				color.New(color.FgHiMagenta, color.Faint).Print("■")
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
lost:
	for {
		moves := ai.BestMoves(g.Board(), g.NextBag())
		for i, m := range moves {
			print("\033[H\033[2J")
			printBoard(g.Board())
			println()
			printMoves(moves[i:])
			time.Sleep(1000 * time.Millisecond)

			if !g.Place(m.Piece, m.X, m.Y) {
				break lost
			}
		}
		n++
	}
	elapsed := time.Since(start)
	print("\033[H\033[2J")
	printBoard(g.Board())
	fmt.Printf("\nFinal Score: %v\nPlayed %v bags in %v\nAverage move time: %v/bag\n", g.Score(), n, elapsed, elapsed/time.Duration(n))
}

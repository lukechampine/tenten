package main

import (
	"fmt"
	"time"

	"github.com/fatih/color"

	"github.com/lukechampine/tenten/ai"
	"github.com/lukechampine/tenten/game"
)

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
	seed := time.Now().Unix()
	g := game.New(seed)
	n := 0
	start := time.Now()
lost:
	for {
		moves := ai.BestMoves(g.Board(), g.NextBag())
		print("\033[H\033[2J")
		println("Score:", g.Score())
		println("Heuristic:", ai.Heuristic(g.Board()))
		println()
		println(g.Board().String())
		println()
		printMoves(moves[:])

		for i, m := range moves {
			time.Sleep(0 * time.Millisecond)
			if !g.Place(m.Piece, m.X, m.Y) {
				break lost
			}

			print("\033[H\033[2J")
			println("Score:", g.Score())
			println("Heuristic:", ai.Heuristic(g.Board()))
			println()
			println(g.Board().String())
			println()
			if i+1 < len(moves) {
				printMoves(moves[i+1:])
			} else {
				println("\nThinking...")
			}
		}
		n++
	}
	elapsed := time.Since(start)
	print("\033[H\033[2J")
	println(g.Board().String())
	fmt.Printf("\nFinal Score: %v\nPlayed %v bags in %v\nAverage move time: %v/bag\nSeed: %v\n", g.Score(), n, elapsed, elapsed/time.Duration(n), seed)
}

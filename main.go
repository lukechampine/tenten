package main

import (
	"fmt"
	"time"

	"github.com/fatih/color"

	"github.com/lukechampine/tenten/ai"
	"github.com/lukechampine/tenten/game"
)

const moveDelay = 500 * time.Millisecond
const thinkingTime = 1000 * time.Millisecond

func printBag(pieces [3]game.Piece, placed [3]bool) {
	var grid [5][18]game.Color
	for i, p := range pieces {
		if placed[i] {
			continue
		}
		for _, d := range p.Dots() {
			x := d.X + (len(pieces)-i-1)*6
			grid[d.Y][x] = p.Color()
		}
	}

	for _, row := range grid {
		for _, c := range row {
			switch c {
			case game.Empty:
				fmt.Print(" ")
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
			fmt.Print(" ")
		}
		fmt.Println()
	}
}

func formatEvaled(n int) string {
	if n > 100000 {
		return color.GreenString("%v", n)
	} else if n > 10000 {
		return color.YellowString("%v", n)
	} else {
		return color.RedString("%v", n)
	}
}

func main() {
	seed := time.Now().Unix()
	g := game.New(seed)
	n := 0
	start := time.Now()
lost:
	for {
		bag := g.NextBag()
		fmt.Printf("\033[H\033[2JScore: %v\nEvaluated ... board states\n\n%v\n", g.Score(), g.Board())
		placed := [3]bool{false, false, false}
		printBag(bag, placed)

		moves, evaled := ai.BestMoves(g.Board(), bag, thinkingTime)
		for _, m := range moves {
			if !g.Place(m.Piece, m.X, m.Y) {
				break lost
			}
			// remove piece
			for i := range bag {
				if bag[i] == m.Piece && !placed[i] {
					placed[i] = true
					break
				}
			}
			fmt.Printf("\033[H\033[2JScore: %v\nEvaluated %v board states\n\n%v\n",
				g.Score(), formatEvaled(evaled), g.Board())
			printBag(bag, placed)
			time.Sleep(moveDelay)
		}
		n++
	}
	elapsed := time.Since(start)
	fmt.Print("\033[H\033[2J", g.Board().String())
	fmt.Printf("\n\nFinal Score: %v\nPlayed %v bags in %v\nAverage move time: %v/bag\nSeed: %v\n", g.Score(), n, elapsed, elapsed/time.Duration(n), seed)
}

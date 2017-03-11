package ai

import (
	"math"

	"github.com/lukechampine/tenten/game"
)

func Heuristic(b *game.Board) float64 {
	return makeHBoard(b, 0).heuristic()
}

type Move struct {
	Piece game.Piece
	X, Y  int
}

func BestMoves(b *game.Board, bag [3]game.Piece) [3]Move {
	perms := [][3]game.Piece{
		{bag[0], bag[1], bag[2]},
		{bag[0], bag[2], bag[1]},
		{bag[1], bag[0], bag[2]},
		{bag[1], bag[2], bag[0]},
		{bag[2], bag[0], bag[1]},
		{bag[2], bag[1], bag[0]},
	}
	bestPerm := perms[0]
	var bestX, bestY [3]int
	maxH := math.Inf(-1)
	for _, perm := range perms {
		hb1 := makeHBoard(b, perm[0])
		hb2 := makeHBoard(b, perm[1])
		hb3 := makeHBoard(b, perm[2])
		for x1 := 0; x1 <= 10-perm[0].Width(); x1++ {
			for y1 := 0; y1 <= 10-perm[0].Height(); y1++ {
				if !hb1.place(x1, y1) {
					continue
				}
				copyHBoard(hb2, hb1)
				for x2 := 0; x2 <= 10-perm[1].Width(); x2++ {
					for y2 := 0; y2 <= 10-perm[1].Height(); y2++ {
						if !hb2.place(x2, y2) {
							continue
						}
						copyHBoard(hb3, hb2)
						for x3 := 0; x3 <= 10-perm[2].Width(); x3++ {
							for y3 := 0; y3 < 10-perm[2].Height(); y3++ {
								if !hb3.place(x3, y3) {
									continue
								}
								if h := hb3.heuristic(); h > maxH {
									maxH = h
									bestPerm = perm
									bestX[0], bestX[1], bestX[2], bestY[0], bestY[1], bestY[2] = x1, x2, x3, y1, y2, y3
								}
							}
						}
					}
				}
			}
		}
	}
	return [3]Move{
		{bestPerm[0], bestX[0], bestY[0]},
		{bestPerm[1], bestX[1], bestY[1]},
		{bestPerm[2], bestX[2], bestY[2]},
	}
}

package ai

import "github.com/lukechampine/tenten/game"

func heuristic(b *game.Board) int {
	// maximize empty spaces
	empty := 0
	for i := range b {
		for _, c := range b[i] {
			if c == game.Empty {
				empty++
			}
		}
	}
	return empty
}

func Move(g *game.Game, bag [3]game.Piece) bool {
	perms := [][3]game.Piece{
		{bag[0], bag[1], bag[2]},
		{bag[0], bag[2], bag[1]},
		{bag[1], bag[0], bag[2]},
		{bag[1], bag[2], bag[0]},
		{bag[2], bag[0], bag[1]},
		{bag[2], bag[1], bag[0]},
	}
	var scratch game.Board
	var bestPerm [3]game.Piece
	var bestX, bestY [3]int
	maxH := 0
	for _, perm := range perms {
		for x1 := range scratch {
			for y1 := range scratch {
				for x2 := range scratch {
					for y2 := range scratch {
						for x3 := range scratch {
							for y3 := range scratch {
								g.Board().Copy(&scratch)
								if scratch.Place(perm[0], x1, y1) > 0 && scratch.Place(perm[1], x2, y2) > 0 && scratch.Place(perm[2], x3, y3) > 0 {
									if h := heuristic(&scratch); h > maxH {
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
	}
	for i, p := range bestPerm {
		if !g.Place(p, bestX[i], bestY[i]) {
			return false
		}
	}
	return true
}

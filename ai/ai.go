package ai

import "github.com/lukechampine/tenten/game"

func heuristic(b *game.Board) int {
	empty := 0
	for i := range b {
		for _, c := range b[i] {
			// maximize empty spaces
			if c == game.Empty {
				empty++
			}
		}
	}
	// apply weights
	h := empty * 1
	if !holdsLine1x5(b) {
		h -= 50
	}
	if !holdsLine5x1(b) {
		h -= 50
	}
	if !holdsSq3x3(b) {
		h -= 100
	}
	return h
}

func holdsLine1x5(b *game.Board) bool {
	for y := range b {
		for x := 0; x < 10; x++ {
			if b.IsEmpty(x, y) {
				n := 1
				for ; x < 10 && b.IsEmpty(x, y); x++ {
					n++
				}
				if n >= 5 {
					return true
				}
			}
		}
	}
	return false
}

func holdsLine5x1(b *game.Board) bool {
	for x := range b {
		for y := 0; y < 10; y++ {
			if b.IsEmpty(x, y) {
				n := 1
				for ; y < 10 && b.IsEmpty(x, y); y++ {
					n++
				}
				if n >= 5 {
					return true
				}
			}
		}
	}
	return false
}

func holdsSq3x3(b *game.Board) bool {
	for y := 0; y < 7; y++ {
		for x := 0; x < 7; x++ {
			if b.IsEmpty(x, y) && b.IsEmpty(x+1, y) && b.IsEmpty(x+2, y) &&
				b.IsEmpty(x, y+1) && b.IsEmpty(x+1, y+1) && b.IsEmpty(x+2, y+1) &&
				b.IsEmpty(x, y+2) && b.IsEmpty(x+1, y+2) && b.IsEmpty(x+2, y+2) {
				return true
			}
		}
	}
	return false
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
	maxH := -1000000
	for _, perm := range perms {
		for x1 := range scratch {
		loop1:
			for y1 := range scratch {
				for x2 := range scratch {
				loop2:
					for y2 := range scratch {
						for x3 := range scratch {
						loop3:
							for y3 := range scratch {
								g.Board().Copy(&scratch)
								if scratch.Place(perm[0], x1, y1) <= 0 {
									continue loop1
								} else if scratch.Place(perm[1], x2, y2) <= 0 {
									continue loop2
								} else if scratch.Place(perm[2], x3, y3) <= 0 {
									continue loop3
								}
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
	for i, p := range bestPerm {
		if !g.Place(p, bestX[i], bestY[i]) {
			return false
		}
	}
	return true
}

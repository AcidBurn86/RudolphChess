package main

import (
	"fmt"

	"github.com/notnil/chess"
)

var boardForCuts [8][8]int8

func gameToBoard8x8(game *chess.Game) [8][8]int8 {

	fen := game.FEN()

	fmt.Println(fen)

	return boardForCuts
}

func hasToCut(game *chess.Game) {
	gameToBoard8x8(game)

}

func printPrettyBoard2(boardForCuts [8][8]int8) {
	var b string
	for i := 0; i < 8; i++ {
		b = ""
		for j := 0; j < 8; j++ {
			b = b + fmt.Sprintf("%d ", boardForCuts[i][j])
		}
		fmt.Println(b)
	}

}

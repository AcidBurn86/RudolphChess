package main

import (
	"fmt"
)

var board [8][8]int8

type Move struct {
	FromX, FromY, ToX, ToY int8
}

var moveValidsFromTo []Move

func setupInitialBoard() {
	board[0][0] = 4
	board[0][1] = 2
	board[0][2] = 3
	board[0][3] = 5
	board[0][4] = 6
	board[0][5] = 3
	board[0][6] = 2
	board[0][7] = 4

	board[1][0] = 1
	board[1][1] = 1
	board[1][2] = 1
	board[1][3] = 1
	board[1][4] = 1
	board[1][5] = 1
	board[1][6] = 1
	board[1][7] = 1

	for i := 2; i < 6; i++ {
		for j := 0; j < 8; j++ {
			board[i][j] = 0
		}
	}

	board[6][0] = -1
	board[6][1] = -1
	board[6][2] = -1
	board[6][3] = -1
	board[6][4] = -1
	board[6][5] = -1
	board[6][6] = -1
	board[6][7] = -1

	board[7][0] = -4
	board[7][1] = -2
	board[7][2] = -3
	board[7][3] = -5
	board[7][4] = -6
	board[7][5] = -3
	board[7][6] = -2
	board[7][7] = -4

}

func printPrettyBoard() {
	var b string
	for i := 0; i < 8; i++ {
		b = ""
		for j := 0; j < 8; j++ {
			b = b + fmt.Sprintf("%d ", board[i][j])
		}
		fmt.Println(b)
	}

}

func generateMoves() {
	var x int8
	var y int8
	for y = 0; y < 8; y++ {
		for x = 0; x < 8; x++ {
			generateMovesFrom(x, y)
			//fmt.Println(board[x][y])
		}
	}

	for _, move := range moveValidsFromTo {
		fmt.Println(move)
	}

}

func generateMovesFrom(x int8, y int8) {
	if board[x][y] == 2 {
		generateMovesWitheKnights(x, y)
	}
	if board[x][y] == 1 {
		generateMovesWhitePawn(x, y)
	}
}

func generateMovesWitheKnights(x, y int8) {
	addMoveIfValid := func(toX, toY int8) {
		if toX >= 0 && toX < 8 && toY >= 0 && toY < 8 && board[toX][toY] <= 0 && board[toX][toY] != 5 {
			moveValidsFromTo = append(moveValidsFromTo, Move{FromX: x, FromY: y, ToX: toX, ToY: toY})
		}
	}

	addMoveIfValid(x+2, y+1)
	addMoveIfValid(x+2, y-1)
	addMoveIfValid(x-2, y+1)
	addMoveIfValid(x-2, y-1)
	addMoveIfValid(x+1, y+2)
	addMoveIfValid(x+1, y-2)
	addMoveIfValid(x-1, y+2)
	addMoveIfValid(x-1, y-2)
}

func generateMovesWhitePawn(x, y int8) {
	fmt.Println("Pawn: ", x, y)
	if x == 1 {
		if board[x][y+1] == 0 {
			moveValidsFromTo = append(moveValidsFromTo, Move{FromX: x, FromY: y, ToX: x, ToY: y + 1})
		}
		if board[x][y+2] == 0 && board[x][y+1] == 0 {
			moveValidsFromTo = append(moveValidsFromTo, Move{FromX: x, FromY: y, ToX: x, ToY: y + 2})
		}
	} else {
		if y+1 < 8 && board[x][y+1] == 0 {
			moveValidsFromTo = append(moveValidsFromTo, Move{FromX: x, FromY: y, ToX: x, ToY: y + 1})
		}
	}
}

package main

import (
	"fmt"
	"os"

	"github.com/notnil/chess"
)

/*func nullMovePruning(gameP *chess.Game, bestSequence []string, isWhiteTurn bool) []*chess.Move {

	// get fen from gameP
	gameP.Position().ChangeTurn()
	originalFen := gameP.FEN()
	originalFen2, _ := chess.FEN(originalFen)
	gameP2 := chess.NewGame(originalFen2, chess.UseNotation(chess.UCINotation{}))

	beta := lastBestOpponentSideScore
	if isWhiteTurn {
		beta = lastBestOwnSideScore
	}

	moves := gameP2.ValidMoves()

	//nonOptimalMoves := make([]*chess.Move, 0, len(moves))
	OptimalMoves := make([]*chess.Move, 0, len(moves))

	var score int16 = 0

	for _, move := range moves {
		gameP2 := chess.NewGame(originalFen2, chess.UseNotation(chess.UCINotation{}))
		gameP2.Move(move)

		score, _ = SearchBestMove(gameP2, 1, -beta, -beta+10, !isWhiteTurn, 2, bestSequence, 2, false, true)
		if score*-1 >= beta {
			//nonOptimalMoves = append(nonOptimalMoves, move)
		} else {
			OptimalMoves = append(OptimalMoves, move)
		}
	}
	//fmt.Println("beta: ", beta)
	if len(OptimalMoves) == 0 {
		return moves
	}

	//fmt.Println("Moves before null move prunning: ", len(moves))
	//	fmt.Println("Moves after null move prunning: ", len(OptimalMoves))
	return OptimalMoves

}*/

func cutEarlyQueenTowerMoves(game *chess.Game, moves []*chess.Move) []*chess.Move {

	if moveCountInternal > 6 {
		return moves
	}

	validMoves := make([]*chess.Move, 0, len(moves))
	board := game.Position().Board()

	for _, move := range moves {
		piece := board.Piece(move.S1())
		isCapture := board.Piece(move.S2()) != chess.NoPiece
		isCheck := move.IsCheck()

		// Si es un movimiento de torre, debe ser una captura o un jaque para ser válido
		if piece.Type() == chess.Rook { //|| piece.Type() == chess.Queen {
			if isCapture || isCheck {
				validMoves = append(validMoves, move)
			}
		} else {
			validMoves = append(validMoves, move)
		}
	}

	return validMoves
}

func cutBadKingSafety(game *chess.Game, moves []*chess.Move, isWhiteTurn bool) []*chess.Move {
	evalOnlyKingSafety = true
	fmt.Println("init cutBadKingSafety")

	eval11 := EvaluateBoard(game.Position().Board(), game, isWhiteTurn, false, 1, 1)
	fmt.Println("Eval1: ", eval11)
	validMoves := make([]*chess.Move, 0, len(moves))

	for _, move := range moves {
		clonedGame := game.Clone()

		clonedGame.Move(move)
		eval2 := EvaluateBoard(clonedGame.Position().Board(), clonedGame, !isWhiteTurn, false, 1, 1) * -1

		fmt.Println("Eval2: ", eval2, ">=", eval11, "Move: ", move)
		if eval2 >= eval11+30 {
			validMoves = append(validMoves, move)
		}
		os.Exit(0)
	}

	if len(validMoves) < 2 {
		return moves
	}
	evalOnlyKingSafety = false

	fmt.Println("Moves before king safety prunning: ", len(moves))
	fmt.Println("Moves after king safety prunning: ", len(validMoves))
	return validMoves
}

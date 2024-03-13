package main

import (
	//import vendor notnil

	"github.com/notnil/chess"
)

func SearchQsc(gameQsc *chess.Game, depth int8, alpha int16, beta int16, isWhiteTurn bool, originalDepth int8, bestSequence []string) (int16, []string) {

	if depth == originalDepth+1 {

		eval := EvaluateBoard(gameQsc.Position().Board(), gameQsc, isWhiteTurn, true, originalDepth, 2)

		return eval, bestSequence
	}

	var bestScore int16
	var sequenceForThisDepth []string
	var currentSequence []string
	var score int16
	var sequence []string

	if game.Position().Turn() == chess.White {
		bestScore = -9999
	} else {
		bestScore = 9999
	}

	validMoves := gameQsc.ValidMoves()

	validMoves = leaveOnlyCaptureOrCheckMoves(gameQsc, validMoves)

	for _, move := range validMoves {
		gameClone := gameQsc.Clone()
		gameClone.Move(move)
		currentSequence = append(bestSequence, move.String())

		score, sequence = SearchQsc(gameClone, depth+1, alpha, beta, !isWhiteTurn, originalDepth, currentSequence)

		if (!isWhiteTurn && score < bestScore) || (isWhiteTurn && score > bestScore) {
			bestScore = score
			sequenceForThisDepth = sequence
		}

		// Luego, dependiendo de isWhiteTurn y depth, actualizamos alpha y beta
		if (!isWhiteTurn && depth%2 == 0) || (isWhiteTurn && depth%2 != 0) {
			alpha = max(alpha, bestScore)
		} else {
			beta = min(beta, bestScore)
		}

		if beta <= alpha {
			break
		}
	}

	return bestScore, sequenceForThisDepth
}

func leaveOnlyCaptureOrCheckMoves(gameQscClean *chess.Game, moves []*chess.Move) []*chess.Move {
	if len(moves) == 0 {
		return moves
	}

	var validMovesN []*chess.Move

	for _, move := range moves {
		gameClone := gameQscClean.Clone()

		isCapture := gameClone.Position().Board().Piece(move.S2()) != chess.NoPiece
		gameClone.Move(move)
		isCheck := move.IsCheck()
		if isCheck || isCapture {
			validMovesN = append(validMovesN, move)
		}
	}
	if len(validMovesN) == 0 {
		return moves
	}
	return validMovesN
}

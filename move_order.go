package main

import (
	"sort"

	"github.com/notnil/chess"
)

type ScoredMove struct {
	move  *chess.Move
	score int16
}

func reorderMovesForCaptures1(gameR *chess.Game, moves []*chess.Move, prunn bool, depth int8) []*chess.Move {
	if len(moves) == 0 {
		return nil
	}

	var scoredMoves []ScoredMove
	var score int16

	for _, move := range moves {
		movePiece := gameR.Position().Board().Piece(move.S1())
		targetPiece := gameR.Position().Board().Piece(move.S2())
		//board := gameR.Position().Board()
		score = 0 // Base score for moves not fitting any special criteria
		if move.String() == bestMovePrevDepth {
			score += 150
		} else if move.String() == bestMovePrevDepth2 {
			score += 100
		}
		// Check and Best Move Boost
		if move.IsCheck() {
			score += 600
		}

		// Capture Moves Enhancement
		if targetPiece != chess.NoPiece {
			score += 300 - abs((pieceValue(targetPiece.Type()))-(pieceValue(movePiece.Type())))
		}

		if (movePiece.Type() == chess.Pawn) && ((move.S2().Rank() == 7 && movePiece.Color() == chess.White) || (move.S2().Rank() == 0 && movePiece.Color() == chess.Black)) {
			score += 150 // Boost for moves that advance pawns closer to promotion
		}
		if (movePiece.Type() == chess.Pawn) && ((move.S2().Rank() == 6 && movePiece.Color() == chess.White) || (move.S2().Rank() == 1 && movePiece.Color() == chess.Black)) {
			score += 125 // Boost for moves that advance pawns closer to promotion
		}

		if movePiece.Type() == chess.Rook && move.S2().Rank() == 6 && movePiece.Color() == chess.White {
			score += 20
		}

		if movePiece.Type() == chess.Rook && move.S2().Rank() == 1 && movePiece.Color() == chess.Black {
			score += 20
		}

		scoredMoves = append(scoredMoves, ScoredMove{move, score})

	}

	if len(scoredMoves) == 0 {
		return moves
	}

	// Sort the moves based on the calculated scores
	sort.Slice(scoredMoves, func(i, j int) bool {
		return scoredMoves[i].score > scoredMoves[j].score
	})

	// Reconstruct the ordered list of moves based on their scores
	orderedMoves := make([]*chess.Move, len(scoredMoves))
	for i, scoredMove := range scoredMoves {
		orderedMoves[i] = scoredMove.move
	}

	return orderedMoves
}

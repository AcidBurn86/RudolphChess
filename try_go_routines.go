package main

import "github.com/notnil/chess"

type MoveResult struct {
	Score    int16
	Sequence []string
}

func SplitValidMoves(game *chess.Game, depth int8, alpha int16, beta int16, isWhiteTurn bool, originalDepth int8, bestSequence []string, maxDepth int8, canCountNodes bool, notNull bool, lwc bool, lbc bool) (int16, []string) {
	resultChan1 := make(chan MoveResult)
	resultChan2 := make(chan MoveResult)

	go func() {
		score, sequence := SearchBestMove(game, depth, alpha, beta, isWhiteTurn, originalDepth, bestSequence, maxDepth, canCountNodes, notNull, lwc, lbc)
		resultChan1 <- MoveResult{Score: score, Sequence: sequence}
	}()

	go func() {
		score, sequence := SearchBestMove(game, depth, alpha, beta, isWhiteTurn, originalDepth, bestSequence, maxDepth, canCountNodes, notNull, lwc, lbc)
		resultChan2 <- MoveResult{Score: score, Sequence: sequence}
	}()

	result1 := <-resultChan1
	result2 := <-resultChan2

	var bestResult MoveResult
	if isWhiteTurn {
		if result1.Score > result2.Score {
			bestResult = result1
		} else {
			bestResult = result2
		}
	} else {
		if result1.Score < result2.Score {
			bestResult = result1
		} else {
			bestResult = result2
		}
	}

	// Retorna el mejor resultado
	return bestResult.Score, bestResult.Sequence

}

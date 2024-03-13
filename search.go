package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/notnil/chess"
)

type MostCutMoves struct {
	moves   string
	counter int
}

var bestMovePrevDepth string
var bestMovePrevDepth2 string

var bestCutMovesMap map[string]*MostCutMoves

var cutLines []string

func FindAndIncreaseCutMove(moveCut string) {
	// Crea una clave única para los movimientos
	key := (moveCut)

	// Busca la clave en el mapa
	if move, found := bestCutMovesMap[key]; found {
		move.counter++
	} else {
		// Si la clave no existe, agrega la nueva entrada al mapa
		bestCutMovesMap[key] = &MostCutMoves{key, 1} // Almacena la clave como string
	}
}

func GetCutMoveCount(move string) int {
	// Crea una clave única para los movimientos
	key := (move)

	// Busca la clave en el mapa
	if moveData, found := bestCutMovesMap[key]; found {
		return moveData.counter
	}
	return 0 // Si el movimiento no se encuentra, devuelve 0
}

var cutedForDepth1 []*chess.Move

func cutFrom1To5(game *chess.Game, validMoves []*chess.Move) []*chess.Move {
	var movesNotInCutLines []*chess.Move

	for _, move := range validMoves {
		isInCutLines := false
		for _, cutLine := range cutLines {
			if move.String() == cutLine {
				isInCutLines = true
				break
			}
		}

		//piece1 := game.Position().Board().Piece(move.S1())
		piece2 := game.Position().Board().Piece(move.S2())

		if isInCutLines && !move.IsCheck() && piece2 == chess.NoPiece &&
			//piece1.Type() != chess.King &&
			move.String() != bestMovePrevDepth {
		} else {
			movesNotInCutLines = append(movesNotInCutLines, move)
		}

	}

	if len(movesNotInCutLines) <= 3 {
		cutedForDepth1 = validMoves

		return validMoves
	} else {

		if len(movesNotInCutLines) > 0 {
			cutedForDepth1 = movesNotInCutLines
			return movesNotInCutLines
		} else {
			cutedForDepth1 = validMoves
			return validMoves
		}
	}

}

var lastIsWhiteCapturing bool
var lastIsBlackCapturing bool

var threshold int16 = 120

var depthToCatchMoves int8 = 5

//var isTest2 bool = false

var qs bool = false

func qsc(pos []string, gameCloneReprocess *chess.Game) (int16, []string) {
	fmt.Println("inicia qsc")
	moves := pos
	qs = true

	isWhiteTurn := gameCloneReprocess.Position().Turn() == chess.White

	alpha := int16(-9999)
	beta := int16(9999)

	score, sequence := SearchBestMove(gameCloneReprocess, 1, alpha, beta, !isWhiteTurn, 2, []string{}, 2, true, false, false, false)

	sequence = append(moves, sequence...)

	qs = false

	fmt.Println("termina qsc", score, sequence)
	//os.Exit(0)
	return score, sequence

}

var foundedHash int

func SearchBestMove(game *chess.Game, depth int8, alpha int16, beta int16, isWhiteTurn bool, originalDepth int8, bestSequence []string, maxDepth int8, canCountNodes bool, notNull bool, lwc bool, lbc bool) (int16, []string) {
	if depth == originalDepth+1 {

		var edgeScore int16 = 0

		if lwc {
			edgeScore = edgeScore_value
		}
		if lbc {
			edgeScore = -edgeScore_value
		}

		/*if score, found := getPositionScore(game.FEN()); found {
			foundedHash++
			return int16(score) + edgeScore, bestSequence
		} else {*/
		return EvaluateBoard(game.Position().Board(), game, isWhiteTurn, canCountNodes, originalDepth, maxDepth) + edgeScore, bestSequence
		//}
	}

	if isOpening {
		threshold = threshold_isOpening
	}
	if isEndgame {
		threshold = threshold_isEndgame
	}
	if isMiddleGame {
		threshold = threshold_isMiddleGame
	}

	var prunn bool = false

	var bestScore int16
	if isWhiteTurn {
		bestScore = -9999
	} else {
		bestScore = 9999
	}
	sequenceForThisDepth := make([]string, 0)

	var validMoves []*chess.Move

	if len(cutedForDepth1) > 0 && depth == 100 && !isReProcessBestLines {
		fmt.Println("Lines Now= ", len(validMoves))
	} else {
		validMoves = game.ValidMoves(false)
		if len(validMoves) > 35 && depth >= 15 {
			validMoves = validMoves[:25]
		}
	}

	lenValidMoves := len(validMoves)
	validMoves = reorderMovesForCaptures(game, validMoves, prunn, depth)
	validMoves = cutEarlyQueenTowerMoves(game, validMoves)

	if originalDepth == depthToCatchMoves+1 && depth == 100 && !isReProcessBestLines {
		//validMoves = cutFrom1To5(game, validMoves)
		fmt.Println("Lines depth 1= ", lenValidMoves, "Now= ", len(validMoves))
	}

	for _, move := range validMoves {

		gameClone := game.Clone()
		gameClone.Move(move)

		currentSequence := append(bestSequence, move.String())

		if isWhiteTurn {
			lastIsWhiteCapturing = game.Position().Board().Piece(move.S2()) != chess.NoPiece
			lastIsBlackCapturing = false
		} else {
			lastIsBlackCapturing = game.Position().Board().Piece(move.S2()) != chess.NoPiece
			lastIsWhiteCapturing = false
		}

		score, sequence := SearchBestMove(gameClone, depth+1, alpha, beta, !isWhiteTurn, originalDepth, currentSequence, maxDepth, canCountNodes, false, lastIsWhiteCapturing, lastIsBlackCapturing)
		lastIsWhiteCapturing = false
		lastIsBlackCapturing = false
		if isWhiteTurn {
			if score > bestScore {
				/*sequenceStr := strings.Join(sequence, " ")
				fmt.Printf("info depth %d score cp %d pv %s nodes %d\n", originalDepth, score, sequenceStr, nodeCount)*/

				if depth == 1 && sequenceForThisDepth != nil && !isReProcessBestLines && canCountNodes && !qs {
					sequenceStr := strings.Join(sequence, " ")
					if !hide_output {
						fmt.Printf("info depth %d score cp %d pv %s nodes %d\n", originalDepth, score, sequenceStr, nodeCount)
					}
					//joinedString := strings.Join(sequence, " ")
					//AgregarElemento(joinedString, (score))
				}
				bestScore = score
				sequenceForThisDepth = sequence

				if initEval[1]-threshold > bestScore && (depth == depthToCatchMoves) {
					existe := false
					if maxDepth > depthToCatchMoves {
						for _, valor := range cutLines {
							if valor == sequence[0] {
								existe = true
								break
							}
						}

						if !existe {
							cutLines = append(cutLines, sequence[0])
						}
					}
					break
				}
				if score >= beta {
					if canCountNodes {
						FindAndIncreaseCutMove(move.String())
					}
					break // Corte Beta
				}
				alpha = max(alpha, score)
			} else {
				if score+180 > bestScore {
					if depth == 1 && sequenceForThisDepth != nil && !isReProcessBestLines && canCountNodes && !qs {
						sequenceStr := strings.Join(sequence, " ")
						if !hide_output {
							fmt.Printf("info depth %d score cp %d pv %s nodes %d\n", originalDepth, score, sequenceStr, nodeCount)
						}
						joinedString := strings.Join(sequence, " ")
						AgregarElemento(joinedString, (score))
					}
				}
			}
		} else {
			/*sequenceStr := strings.Join(sequence, " ")
			if depth == 6 && strings.HasPrefix(sequenceStr, "c6d6") {
				fmt.Printf("info depth %d score cp %d pv %v nodes %d\n", originalDepth, score, sequenceStr, nodeCount)
			}*/
			if score < bestScore {
				if depth == 1 && sequenceForThisDepth != nil && !isReProcessBestLines && canCountNodes && !qs {
					sequenceStr := strings.Join(sequence, " ")
					if !hide_output {
						fmt.Printf("info depth %d score cp %d pv %v nodes %d\n", originalDepth, score, sequenceStr, nodeCount)
					}
					joinedString := strings.Join(sequence, " ")
					AgregarElemento(joinedString, (score))
				}
				bestScore = score
				sequenceForThisDepth = sequence

				if initEval[1]+threshold > bestScore && (depth == depthToCatchMoves) {
					existe := false
					if maxDepth > depthToCatchMoves {
						for _, valor := range cutLines {
							if valor == sequence[0] {
								existe = true
								break
							}
						}

						if !existe {
							cutLines = append(cutLines, sequence[0])
						}
					}
					break
				}
				if score <= alpha {
					if canCountNodes {
						FindAndIncreaseCutMove(move.String())
					}

					break // Corte Alpha
				}
				beta = min(beta, score)
			} else {
				if score-180 < bestScore {
					if depth == 1 && sequenceForThisDepth != nil && !isReProcessBestLines && canCountNodes && !qs {
						sequenceStr := strings.Join(sequence, " ")
						if !hide_output {
							fmt.Printf("info depth %d score cp %d pv %s nodes %d\n", originalDepth, score, sequenceStr, nodeCount)
						}
						joinedString := strings.Join(sequence, " ")
						AgregarElemento(joinedString, (score))
					}
				}
			}
		}

	}

	if len(sequenceForThisDepth) > 0 {
		return bestScore, sequenceForThisDepth
	} else {
		return bestScore, bestSequence
	}
}

func reorderMovesForCaptures(gameR *chess.Game, moves []*chess.Move, prunn bool, depth int8) []*chess.Move {
	if len(moves) == 0 {
		return nil
	}

	var scoredMoves []ScoredMove
	var score int16

	for _, move := range moves {
		movePiece := gameR.Position().Board().Piece(move.S1())
		targetPiece := gameR.Position().Board().Piece(move.S2())
		score = 100 // Base score for moves not fitting any special criteria

		// Check and Best Move Boost
		if move.IsCheck() {
			score += 1970
		} else if move.String() == bestMovePrevDepth {
			score += 899
		} else if move.String() == bestMovePrevDepth2 {
			score += 898
		}

		// Capture Moves Enhancement
		if targetPiece != chess.NoPiece {
			captureScore := 1000 - abs((pieceValue(targetPiece.Type()))-(pieceValue(movePiece.Type())))
			if captureScore > score {
				score = captureScore
			}
			if pieceValue(targetPiece.Type()) < pieceValue(movePiece.Type()) {
				score += 205
			}
			if movePiece.Type() == chess.Pawn && targetPiece.Type() != chess.Pawn {
				score += 190
			}

		}

		if (movePiece.Type() == chess.Pawn) && ((move.S2().Rank() == 7 && movePiece.Color() == chess.White) || (move.S2().Rank() == 0 && movePiece.Color() == chess.Black)) {
			score += 150 // Boost for moves that advance pawns closer to promotion
		}
		if (movePiece.Type() == chess.Pawn) && ((move.S2().Rank() == 6 && movePiece.Color() == chess.White) || (move.S2().Rank() == 1 && movePiece.Color() == chess.Black)) {
			score += 120 // Boost for moves that advance pawns closer to promotion
		}

		if depth <= 6 {
			scoredMoves = append(scoredMoves, ScoredMove{move, score})
		}
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

func abs(x int16) int16 {
	if x < 0 {
		return -x
	}
	return x
}

func max(a, b int16) int16 {
	if a > b {
		return a
	}
	return b
}

func min(a, b int16) int16 {
	if a < b {
		return a
	}
	return b
}

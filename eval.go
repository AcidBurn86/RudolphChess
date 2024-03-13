package main

import (
	"fmt"
	"sync"

	"github.com/notnil/chess"
)

var (
	hashScores sync.Map // Utilizando sync.Map para manejar la concurrencia
)

var evalThis bool = false

func squareToIndex(square string) int16 {

	col := square[0] // a-h
	row := square[1] // 1-8

	colIndex := int(col - 'a') // Convertir a-h a 0-7
	rowIndex := int(row - '1') // Convertir 1-8 a 0-7

	return int16(rowIndex*8 + colIndex)
}

var evalOnlyKingSafety bool = false

var isOpening bool = false
var isMiddleGame bool = false
var isEndgame bool = false

func MaterialEval(board *chess.Board, gameMovility *chess.Game, isWhiteTurn bool, canCountNodes bool) int16 {
	var wmaterial int16
	var bmaterial int16

	for sq := chess.A1; sq <= chess.H8; sq++ { // Iterar sobre todas las casillas
		piece := board.Piece(sq)               // Obtener la pieza en la casilla actual
		if piece.Type() != chess.NoPieceType { // Verificar si hay una pieza en la casilla
			value := pieceValue(piece.Type()) // Obtener el valor de la pieza

			if piece.Color() == chess.White {
				wmaterial += int16(value)
			} else {
				bmaterial -= int16(value)
			}

		}
	}
	return wmaterial + bmaterial
}

func EvaluateBoard(board *chess.Board, gameMovility *chess.Game, isWhiteTurn bool, canCountNodes bool, originalDepth int8, maxDepth int8) int16 {
	if canCountNodes {
		nodeCount++
	}

	/*hashValue := gameMovility.Position().Hash()

	if score, ok := hashScores.Load(hashValue); ok && !evalThis && !evalOnlyKingSafety {
		//foundedHash++
		//return score.(int16)
	}*/

	var score int16 = 0
	var positionalScore int16 = 0
	var blackPawnScore int16 = 0
	var whitePawnScore int16 = 0

	var countWhiteBishops int8 = 0
	var countBlackBishops int8 = 0
	var sqInt int16
	var balckMovility, whiteMovility int16

	var wROR bool = false
	var bROR bool = false
	//var difMovility int16
	if whiteRookInOpenRank() {
		score += RookInOpenRank
		wROR = true
	}
	if blackRookInOpenRank() {
		score -= RookInOpenRank
		bROR = true
	}

	if !areQueensInTheBoard() && (wROR || bROR) {
		isEndgame = true
		isMiddleGame = false
		isOpening = false
	} else {
		if moveCountInternal < 8 {
			isOpening = true
			isMiddleGame = false
		} else if moveCountInternal > 20 && moveCountInternal < 40 {
			isMiddleGame = true
			isOpening = false
		} else {
			isEndgame = false
		}
	}

	balckMovility, whiteMovility = countMovesByBoard(gameMovility, isWhiteTurn)

	/*	if whiteMovility > 70 {
			whiteMovility += (whiteMovility - 60)
		}
		if balckMovility > 70 {
			balckMovility += (balckMovility - 60)
		}*/

	//difference between movility

	score += checkDoubledPawns(board)

	createBitBoard(gameMovility.Position().Board())

	/*score += evalPiecesCompany(0, 2)
	score += evalPiecesCompany(1, 2)

	score -= checkIsolatedWhitePawns(whitePawns) * 15
	score += checkIsolatedBlackPawns(blackPawns) * 15*/

	kingSafetyWhite := int16(0)
	kingSafetyBlack := int16(0)

	if moveCountInternal > 1 {
		//fmt.Println("in kingsafety")
		kingSafetyWhite = kingSafety(0)
		kingSafetyBlack = kingSafety(1)
	}

	/*difKingSafety := int16(abs((kingSafetyWhite) - (kingSafetyBlack)))
	difMovility = (abs((balckMovility) - (whiteMovility)))
	if difKingSafety > 45 && difMovility > 40 {
		if isWhiteTurn {
			score = score - (difMovility / 2)
		} else {
			score = score + (difMovility / 2)
		}
	}*/

	score += (kingSafetyWhite * kingSafetyMulti) + (kingSafetyBlack * kingSafetyMulti)

	centerScore := (evalCenterControl(isWhiteTurn))

	score += centerScore * center_multi_1

	centerScore2 := (validateCenterControl2(isWhiteTurn))

	score += centerScore2 * center_multi_2

	wmaterial := int16(0)
	bmaterial := int16(0)

	//gameMovility.Position().CastleRights().CanCastle(chess.White, chess.KingSide)

	//fmt.Println("Center Score: ", centerScore)

	for sq := chess.A1; sq <= chess.H8; sq++ { // Iterar sobre todas las casillas
		piece := board.Piece(sq)               // Obtener la pieza en la casilla actual
		if piece.Type() != chess.NoPieceType { // Verificar si hay una pieza en la casilla
			value := pieceValue(piece.Type()) // Obtener el valor de la pieza

			sqInt = squareToIndex(sq.String())

			if piece.Color() == chess.White {
				wmaterial += int16(value)
			} else {
				bmaterial -= int16(value)
			}

			switch piece.Type() {
			case chess.Queen:
				if piece.Color() == chess.White {
					if moveCountInternal < 10 {
						sqText := sq.String()
						if sqText == "d1" {
							score += queenInStartPosition
						}
					}
				} else {
					if moveCountInternal < 10 {
						sqText := sq.String()
						if sqText == "d8" {
							score -= queenInStartPosition
						}
					}
				}
			case chess.Pawn:
				pawn_eval := pawnBonus(int(sq.Rank()), int(sq.File()), isWhiteTurn)

				if piece.Color() == chess.White {
					//countWhitePawns++
					if isEndgame {
						whitePawnScore += pawnTableEnd[sqInt] + pawn_eval
					} else {
						whitePawnScore += pawnTable[sqInt] + pawn_eval
					}
				} else {
					//	countBlackPawns++
					if isEndgame {
						blackPawnScore += pawnTableEnd[flipSquare(sqInt)] + pawn_eval
					} else {
						blackPawnScore -= pawnTable[flipSquare(sqInt)] + pawn_eval
					}
				}

			case chess.Knight:
				if piece.Color() == chess.White {
					if sq.Rank() == 1 && moveCountInternal < 15 {
						positionalScore -= 30
					}
					if isEndgame {
						positionalScore += KnightTableEnd[sqInt]
					} else {
						positionalScore += KnightTable[sqInt]
					}

				} else {
					if sq.Rank() == 8 && moveCountInternal < 15 {
						positionalScore += 30
					}
					if isEndgame {
						positionalScore -= KnightTableEnd[flipSquare(sqInt)]
					} else {
						positionalScore -= KnightTable[flipSquare(sqInt)]
					}
				}

			case chess.Bishop:
				if piece.Color() == chess.White {
					countWhiteBishops++
					if isEndgame {
						positionalScore += BishopTableEnd[sqInt]
					} else {
						if sq.Rank() == 1 {
							positionalScore -= 60
						}
						positionalScore += BishopTable[sqInt]
					}
				} else {
					countBlackBishops++
					if isEndgame {
						positionalScore -= BishopTableEnd[flipSquare(sqInt)]
					} else {
						if sq.Rank() == 1 {
							positionalScore += 60
						}
						positionalScore -= BishopTable[flipSquare(sqInt)]
					}
				}

			case chess.King:
				if piece.Color() == chess.White {
					if isEndgame {
						positionalScore += KingTableEnd[sqInt]
					} else {
						if moveCountInternal < 15 {
							sqText := sq.String()
							if sqText == "e1" {
								score += 130
							}
						}
						positionalScore += KingTable[sqInt]
						pieceh1 := board.Piece(7)
						if (sq.String() == "g1" || sq.String() == "h1") && !gameMovility.Position().CastleRights().CanCastle(chess.White, chess.KingSide) && (pieceh1.Type() == chess.NoPieceType) {
							positionalScore += 230
						}
						if (sq.String() == "f1" || sq.String() == "e2") && game.Position().CastleRights().CanCastle(chess.White, chess.KingSide) {
							//fmt.Println("acá")
							positionalScore -= 230
						}
					}
				} else {
					if isEndgame {
						positionalScore -= KingTableEnd[flipSquare(sqInt)]
					} else {
						if moveCountInternal < 15 {
							sqText := sq.String()
							if sqText == "e8" {
								score -= 130
							}
						}
						positionalScore -= KingTable[flipSquare(sqInt)]
						pieceh8 := board.Piece(63)
						if (sq.String() == "g8" || sq.String() == "h8") && !gameMovility.Position().CastleRights().CanCastle(chess.Black, chess.KingSide) && (pieceh8.Type() == chess.NoPieceType) {
							positionalScore -= 230
						}
						if (sq.String() == "f8" || sq.String() == "e7") && game.Position().CastleRights().CanCastle(chess.Black, chess.KingSide) {
							positionalScore += 230
						}
					}
				}
			}
		}
	}

	var bishopPair int16 = 0
	var mobility int16 = 0

	if countBlackBishops == 2 {
		bishopPair = bishopPair - bishopPairBonus
	}

	if countWhiteBishops == 2 {
		bishopPair = bishopPair + bishopPairBonus
	}

	positionalScore = positionalScore + whitePawnScore + blackPawnScore + bishopPair
	mobility = (whiteMovility - balckMovility)

	whitePawnsDefending := countDefendedByWhitePawns()
	blackPawnsDefending := countDefendedByBlackPawns()

	score += int16(whitePawnsDefending * PawnsDefending_Multi)
	score -= int16(blackPawnsDefending * PawnsDefending_Multi)

	if evalThis {
		fmt.Println("Material: W", wmaterial, "B", bmaterial)
		fmt.Println("White Pawns Defending: ", whitePawnsDefending)
		fmt.Println("Black Pawns Defending: ", blackPawnsDefending)
		fmt.Println("White Mobility: ", whiteMovility, "Black Mobility: ", balckMovility)
		fmt.Println("King Safety White: ", kingSafetyWhite, "King Safety Black: ", kingSafetyBlack)
		fmt.Println("Score: ", score, "Move Count: ", moveCountInternal, "Positional Score: ", positionalScore, "WhitePawns: ", whitePawnScore, "BlackPawns: ", blackPawnScore)
	}
	//return score
	finalScore := (score + positionalScore + mobility + wmaterial + bmaterial)

	// get randon number from 0 to 10

	//randomNumber := rand.Intn(4)
	//fmt.Println("Random Number: ", randomNumber)
	//hashScores.Store(hashValue, finalScore)

	return finalScore //+ int16(randomNumber)
}

func checkDoubledPawns(board *chess.Board) int16 {
	var countBlackPawns [8]int16
	var countWhitePawns [8]int16

	for sq := chess.A2; sq <= chess.H8; sq++ {
		piece := board.Piece(sq)
		if piece.Type() == chess.Pawn {
			col := sq.File()
			if piece.Color() == chess.White {
				countWhitePawns[col]++
			} else {
				countBlackPawns[col]++
			}
		}
	}

	var penalty int16
	for i := 0; i < 8; i++ {
		if countBlackPawns[i] > 1 {
			penalty += doubledPawnPenalty // Penalización por cada columna adicional con peones negros doblados
		}
		if countWhitePawns[i] > 1 {
			penalty -= doubledPawnPenalty
		}
	}

	return penalty
}

type mobilityResult struct {
	blackMobility int16
	whiteMobility int16
}

// Define una variable global para la caché
var cache map[string]mobilityResult

func countMovesByBoard(tempGame *chess.Game, isWhiteTurn bool) (int16, int16) {

	var blackMovilty int16 = 0
	var whiteMovilty int16 = 0

	var wBishopMoves int16 = 0
	var wKnightMoves int16 = 0

	var bBishopMoves int16 = 0
	var bKnightMoves int16 = 0

	// Utiliza el hash para cachear el resultado
	//cacheKey := boardHash

	validMoves := tempGame.ValidMoves(false)
	board := tempGame.Position().Board()

	for i := 1; i < 3; i++ {

		if i == 2 {

			validMoves = tempGame.ValidMoves(true)
			isWhiteTurn = !isWhiteTurn

		}

		if isWhiteTurn {

			for _, move := range validMoves {
				p := board.Piece(move.S1()).Type()
				if p == chess.Knight {
					wKnightMoves++
					whiteMovilty = whiteMovilty + knightMoveValue
				} else if p == chess.Bishop {
					wBishopMoves++
					whiteMovilty = whiteMovilty + bishopMoveValue
				} else if p == chess.Rook {
					whiteMovilty = whiteMovilty + rookMoveValue
				} else if p == chess.Queen {
					whiteMovilty = whiteMovilty + queenMoveValue
				}
			}

		} else {

			for _, move := range validMoves {
				p := board.Piece(move.S1()).Type()
				if p == chess.Knight {
					bKnightMoves++
					blackMovilty = blackMovilty + knightMoveValue
				} else if p == chess.Bishop {
					bBishopMoves++
					blackMovilty = blackMovilty + bishopMoveValue
				} else if p == chess.Rook {
					blackMovilty = blackMovilty + rookMoveValue
				} else if p == chess.Queen {
					blackMovilty = blackMovilty + queenMoveValue
				}
			}
		}

	}

	if wKnightMoves < 2 && moveCountInternal > 10 {
		whiteMovilty = whiteMovilty - mobilityPenalty
	}

	if bKnightMoves < 2 && moveCountInternal > 10 {
		blackMovilty = blackMovilty - mobilityPenalty
	}

	if wBishopMoves < 2 && moveCountInternal > 10 {
		whiteMovilty = whiteMovilty - mobilityPenalty
	}

	if bBishopMoves < 2 && moveCountInternal > 10 {
		blackMovilty = blackMovilty - mobilityPenalty
	}

	return blackMovilty, whiteMovilty
}

func flipSquare(sq int16) int16 {
	return sq ^ 56
}

var pawnTable = [64]int16{
	0, 0, 0, 0, 0, 0, 0, 0,
	10, 10, 0, -10, -10, 0, 10, 10,
	5, 0, 0, 5, 5, 0, 0, 5,
	0, 0, 10, 20, 20, 10, 0, 0,
	5, 5, 5, 10, 10, 5, 5, 5,
	10, 10, 10, 20, 20, 10, 10, 10,
	20, 20, 20, 30, 30, 20, 20, 20,
	0, 0, 0, 0, 0, 0, 0, 0,
}

var pawnTableEnd = [64]int16{
	0, 0, 0, 0, 0, 0, 0, 0,
	-24, -22, -34, -46, -46, -34, -22, -24,
	32, 20, -4, -18, -18, -4, 20, 32,
	36, 14, 10, -4, -4, 10, 14, 26,
	36, 12, 16, 12, 12, 16, 12, 16,
	48, 46, 44, 44, 44, 44, 46, 48,
	68, 66, 68, 66, 66, 68, 66, 68,
	0, 0, 0, 0, 0, 0, 0, 0,
}

var KnightTable = [64]int16{
	-218, -114, -132, -34, -34, -132, -114, -218,
	-78, -60, 10, -30, -30, 10, -60, -78,
	-12, 48, 40, 56, 56, 40, 48, -12,
	14, 50, 46, 58, 58, 46, 50, 14,
	12, 40, 44, 42, 42, 44, 40, 12,
	6, 40, 40, 38, 38, 40, 40, 6,
	-8, 6, 18, 32, 32, 18, 6, -8,
	-30, 2, -6, 12, 12, -6, 2, -30,
}
var KnightTableEnd = [64]int16{
	-16, 6, 36, 18, 18, 36, 6, -16,
	4, 36, 12, 42, 42, 12, 36, 4,
	-6, 10, 32, 30, 30, 32, 10, -6,
	18, 18, 38, 40, 40, 38, 18, 18,
	14, 18, 32, 40, 40, 32, 18, 14,
	8, 4, 12, 30, 30, 12, 4, 8,
	-12, 8, 10, 14, 14, 10, 8, -12,
	-10, -6, 8, 16, 16, 8, -6, -10,
}

var BishopTable = [64]int16{

	-26, 4, -92, -88, -88, -92, 4, -26,
	-52, -14, -2, -38, -38, -2, -14, -52,
	36, 48, 40, 30, 30, 40, 48, 36,
	18, 34, 40, 54, 54, 40, 34, 18,
	32, 40, 44, 60, 60, 44, 40, 32,
	36, 52, 52, 42, 42, 52, 52, 36,
	34, 62, 48, 44, 44, 48, 62, 34,
	8, 36, 30, 50, 50, 30, 36, 8,
}

var BishopTableEnd = [64]int16{
	-34, -43, -14, -16, -16, -22, -34, -47,
	-24, -24, -15, -24, -22, -30, -22, -30,
	-23, -32, -28, -27, -27, -29, -28, -23,
	-31, -20, -17, -21, -14, -23, -26, -25,
	-34, -30, -18, -15, -29, -21, -34, -35,
	-39, -34, -25, -24, -21, -34, -36, -40,
	-47, -47, -41, -37, -29, -43, -49, -58,
	-52, -40, -52, -38, -41, -43, -35, -44,
}

var KingTable = [64]int16{

	0, 5, 5, -10, -10, 0, 10, 5,
	-150, -150, -150, -150, -150, -150, -150, -150,
	-150, -150, -150, -150, -150, -150, -150, -150,
	-170, -170, -170, -170, -170, -170, -170, -170,
	-170, -170, -170, -170, -170, -170, -170, -170,
	-170, -170, -170, -170, -170, -170, -170, -170,
	-170, -170, -170, -170, -170, -170, -170, -170,
	-170, -170, -170, -170, -170, -170, -170, -170,
}

var KingTableEnd = [64]int16{
	-104, -82, 14, -52, -52, 14, -82, -104,
	-28, 22, 50, 28, 28, 50, 22, -28,
	-4, 36, 36, 34, 34, 36, 36, -4,
	-4, 36, 44, 46, 46, 44, 36, -4,
	-12, 10, 28, 40, 40, 28, 10, -12,
	-18, 10, 18, 24, 24, 18, 10, -18,
	-38, -10, 20, 24, 24, 20, -10, -38,
	-72, -40, -24, -30, -30, -24, -40, -72,
}

func pieceValue(pieceType chess.PieceType) int16 {
	threshold := int16(0)
	switch pieceType {
	case chess.Pawn:
		return 100 + threshold
	case chess.Knight:
		return 305 + threshold
	case chess.Bishop:
		return 333 + threshold
	case chess.Rook:
		return 563 + threshold
	case chess.Queen:
		return 950 + threshold
	case chess.King:
		return 0
	default:
		return 0
	}
}

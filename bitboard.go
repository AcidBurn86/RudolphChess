package main

import (
	"github.com/notnil/chess"
)

type Bitboard uint64

func setBit(b Bitboard, position int) Bitboard {
	return b | (1 << position)
}

func squareToBitboardIndex(sq chess.Square) int {
	return int(sq)
}

var (
	whitePawns   Bitboard
	whiteKnights Bitboard
	whiteBishops Bitboard
	whiteRooks   Bitboard
	whiteQueens  Bitboard
	whiteKings   Bitboard

	blackPawns   Bitboard
	blackKnights Bitboard
	blackBishops Bitboard
	blackRooks   Bitboard
	blackQueens  Bitboard
	blackKings   Bitboard

	allWhites   Bitboard
	allBlacks   Bitboard
	allOccupied Bitboard
)

func createBitBoard(board *chess.Board) Bitboard {

	whitePawns = 0
	whiteKnights = 0
	whiteBishops = 0
	whiteRooks = 0
	whiteQueens = 0
	whiteKings = 0

	blackPawns = 0
	blackKnights = 0
	blackBishops = 0
	blackRooks = 0
	blackQueens = 0
	blackKings = 0

	allWhites = 0
	allBlacks = 0

	for sq := chess.A1; sq <= chess.H8; sq++ {
		piece := board.Piece(sq)
		position := squareToBitboardIndex(sq)

		switch piece {
		case chess.WhitePawn:
			whitePawns = setBit(whitePawns, position)
			allWhites = setBit(allWhites, position)
		case chess.WhiteKnight:
			whiteKnights = setBit(whiteKnights, position)
			allWhites = setBit(allWhites, position)
		case chess.WhiteBishop:
			whiteBishops = setBit(whiteBishops, position)
			allWhites = setBit(allWhites, position)
		case chess.WhiteRook:
			whiteRooks = setBit(whiteRooks, position)
			allWhites = setBit(allWhites, position)
		case chess.WhiteQueen:
			whiteQueens = setBit(whiteQueens, position)
			allWhites = setBit(allWhites, position)
		case chess.WhiteKing:
			whiteKings = setBit(whiteKings, position)
			allWhites = setBit(allWhites, position)
		case chess.BlackPawn:
			blackPawns = setBit(blackPawns, position)
			allBlacks = setBit(allBlacks, position)
		case chess.BlackKnight:
			blackKnights = setBit(blackKnights, position)
			allBlacks = setBit(allBlacks, position)
		case chess.BlackBishop:
			blackBishops = setBit(blackBishops, position)
			allBlacks = setBit(allBlacks, position)
		case chess.BlackRook:
			blackRooks = setBit(blackRooks, position)
			allBlacks = setBit(allBlacks, position)
		case chess.BlackQueen:
			blackQueens = setBit(blackQueens, position)
			allBlacks = setBit(allBlacks, position)
		case chess.BlackKing:
			blackKings = setBit(blackKings, position)
			allBlacks = setBit(allBlacks, position)
		}

	}

	allOccupied = allWhites | allBlacks

	return allOccupied

}

func howManyBigPices() int {
	return countBits2(whiteRooks) + countBits2(blackBishops) + countBits2(blackRooks) + countBits2(whiteBishops) + countBits2(blackKnights) + countBits2(whiteKnights)
}

func howManyPawns() int {
	return countBits2(whitePawns) + countBits2(blackPawns)
}

func countBits2(b Bitboard) int {
	var count int
	for b > 0 {
		b &= b - 1
		count++
	}
	return count
}

func areQueensInTheBoard() bool {
	return whiteQueens != 0 || blackQueens != 0
}

type boardState struct {
	whiteRooks  Bitboard
	blackRooks  Bitboard
	allOccupied Bitboard
}

// Mapas de caché para almacenar los resultados de las funciones.
var (
	cacheWhiteRookInOpenRank = make(map[boardState]bool)
	cacheBlackRookInOpenRank = make(map[boardState]bool)
)

const (
	ColumnAMask Bitboard = 0x0101010101010101
	ColumnBMask Bitboard = 0x0202020202020202
	ColumnCMask Bitboard = 0x0404040404040404
	ColumnDMask Bitboard = 0x0808080808080808
	ColumnEMask Bitboard = 0x1010101010101010
	ColumnFMask Bitboard = 0x2020202020202020
	ColumnGMask Bitboard = 0x4040404040404040
	ColumnHMask Bitboard = 0x8080808080808080
)

// whiteRookInOpenRank evalúa si hay una torre blanca en una columna abierta o semiabierta.
func whiteRookInOpenRank() bool {

	state := boardState{whiteRooks: whiteRooks, allOccupied: allOccupied}
	if result, found := cacheWhiteRookInOpenRank[state]; found {
		return result // Retorna el resultado desde la caché si está presente.
	}

	columnMasks := [8]Bitboard{
		ColumnAMask, ColumnBMask, ColumnCMask, ColumnDMask,
		ColumnEMask, ColumnFMask, ColumnGMask, ColumnHMask,
	}

	for _, columnMask := range columnMasks {
		if whiteRooks&columnMask != 0 {
			if (allOccupied & columnMask) == (whiteRooks & columnMask) {
				cacheWhiteRookInOpenRank[state] = true
				return true
			}
		}
	}
	cacheWhiteRookInOpenRank[state] = false
	return false
}

func blackRookInOpenRank() bool {

	state := boardState{blackRooks: blackRooks, allOccupied: allOccupied}
	if result, found := cacheBlackRookInOpenRank[state]; found {
		return result // Retorna el resultado desde la caché si está presente.
	}
	columnMasks := [8]Bitboard{
		ColumnAMask, ColumnBMask, ColumnCMask, ColumnDMask,
		ColumnEMask, ColumnFMask, ColumnGMask, ColumnHMask,
	}
	for _, columnMask := range columnMasks {
		// Verificar si hay una torre negra en la columna actual.
		if blackRooks&columnMask != 0 {
			// Verificar si la columna está abierta (sin piezas) o semiabierta (solo piezas negras).
			if (allOccupied & columnMask) == (blackRooks & columnMask) {
				cacheBlackRookInOpenRank[state] = true
				return true // Hay al menos una torre negra en una columna abierta/semiabierta.
			}
		}
	}
	cacheBlackRookInOpenRank[state] = false
	return false // No hay torres negras en columnas abiertas/semiabiertas.
}

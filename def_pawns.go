package main

import (
	"github.com/notnil/chess"
)

func whitePawnsDefendedSquares(pawns Bitboard) Bitboard {
	leftDefended := (pawns << 9) & 0xFEFEFEFEFEFEFEFE  // Evita el desbordamiento a la siguiente fila
	rightDefended := (pawns << 7) & 0x7F7F7F7F7F7F7F7F // Lo mismo aquí

	return leftDefended | rightDefended
}

// Cuenta las piezas propias defendidas por peones blancos
func countDefendedByWhitePawns() int {
	defendedSquares := whitePawnsDefendedSquares(whitePawns)
	// Piezas propias excluyendo los peones
	ownPieces := allWhites &^ whitePawns

	// Cuenta cuántas de esas casillas están ocupadas por tus piezas
	defendedPieces := defendedSquares & ownPieces

	// Contar los bits (piezas) en defendedPieces
	count := 0
	for defendedPieces != 0 {
		count++
		defendedPieces &= defendedPieces - 1 // Elimina el bit menos significativo
	}

	return count
}

func blackPawnsDefendedSquares(pawns Bitboard) Bitboard {
	leftDefended := (pawns >> 7) & 0xFEFEFEFEFEFEFEFE
	rightDefended := (pawns >> 9) & 0x7F7F7F7F7F7F7F7F

	return leftDefended | rightDefended
}

func countDefendedByBlackPawns() int {
	defendedSquares := blackPawnsDefendedSquares(blackPawns)
	ownPieces := allBlacks &^ blackPawns

	defendedPieces := defendedSquares & ownPieces

	count := 0
	for defendedPieces != 0 {
		count++
		defendedPieces &= defendedPieces - 1
	}

	return count
}

func checkIsolatedWhitePawns(whitePawns Bitboard) int16 {
	var isolated int16 = 0
	for position := 0; position < 64; position++ {
		if whitePawns&(1<<position) != 0 { // Verifica si hay un peón blanco en la posición
			pawnPosition := chess.Square(position)
			if isPawnIsolated(pawnPosition, whitePawns) {
				isolated++
			}
		}
	}
	return isolated
}

// Reutiliza la función isPawnIsolated definida previamente
func isPawnIsolated(pawnPosition chess.Square, whitePawns Bitboard) bool {
	positionIndex := squareToBitboardIndex(pawnPosition)

	var leftMask, rightMask Bitboard
	if positionIndex%8 > 0 {
		leftMask = 1 << (positionIndex - 1)
	}
	if positionIndex%8 < 7 {
		rightMask = 1 << (positionIndex + 1)
	}

	isIsolated := (whitePawns&leftMask) == 0 && (whitePawns&rightMask) == 0

	return isIsolated
}

func checkIsolatedBlackPawns(blackPawns Bitboard) int16 {
	var isolated int16 = 0
	for position := 0; position < 64; position++ {
		if blackPawns&(1<<position) != 0 { // Verifica si hay
			pawnPosition := chess.Square(position)
			if isPawnIsolated(pawnPosition, blackPawns) {
				isolated++
			}
		}
	}
	return isolated
}

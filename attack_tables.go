package main

var bishopAttackTable [64]Bitboard

func initBishopAttacks() {
	for sq := 0; sq < 64; sq++ {
		bishopAttackTable[sq] = calculateBishopAttacks(sq)
	}
}

func calculateBishopAttacks(square int) Bitboard {
	var attacks Bitboard
	rank, file := square/8, square%8

	// Diagonal hacia arriba y a la derecha
	for r, f := rank+1, file+1; r <= 7 && f <= 7; r, f = r+1, f+1 {
		attacks |= (1 << (r*8 + f))
	}
	// Diagonal hacia arriba y a la izquierda
	for r, f := rank+1, file-1; r <= 7 && f >= 0; r, f = r+1, f-1 {
		attacks |= (1 << (r*8 + f))
	}
	// Diagonal hacia abajo y a la derecha
	for r, f := rank-1, file+1; r >= 0 && f <= 7; r, f = r-1, f+1 {
		attacks |= (1 << (r*8 + f))
	}
	// Diagonal hacia abajo y a la izquierda
	for r, f := rank-1, file-1; r >= 0 && f >= 0; r, f = r-1, f-1 {
		attacks |= (1 << (r*8 + f))
	}

	return attacks
}

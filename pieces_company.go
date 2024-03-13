package main

func evalPiecesCompany(color int, pieceType int) int16 {
	/*var score int16 = 0
	var penaltys int16 = 0*/
	var bb Bitboard

	if pieceType == 2 && color == 0 {
		bb = whiteKnights
	}
	if pieceType == 3 && color == 0 {
		bb = whiteBishops
	}

	if pieceType == 2 && color == 1 {
		bb = blackKnights
	}
	if pieceType == 3 && color == 1 {
		bb = blackBishops
	}
	rank := 0
	var positions []int
	for i := 0; i < 64; i++ {
		if bb&(1<<i) != 0 {
			if color == 0 {
				rank = i / 8 // Calcula el rango basado en el índice
			} else {
				rank = 7 - i/8 // Calcula el rango basado en el índice
			}
			if color == 0 && rank > 2 || color == 1 && rank < 6 { // Solo evalúa si el caballo está en un rango > 4
				positions = append(positions, i)
				if len(positions) == 2 {
					break
				}
			}
		}
	}

	count := 0
	var friendlyPieces Bitboard
	if color == 0 {
		friendlyPieces = whiteKnights | whiteBishops | whiteRooks | whiteQueens
		// Piezas amigas son las blancas
	} else {
		friendlyPieces = blackKnights | blackBishops | blackRooks | blackQueens // Piezas amigas son las negras
	}

	for p := range positions {
		surroundings := calculateKingSurroundings(p, 4)

		for _, pos := range surroundings {
			if friendlyPieces&(1<<pos) != 0 {
				count++
			}
		}
	}
	bonus := 0
	if count == 1 {
		bonus = 8
	}
	if count == 2 {
		bonus = 16
	}
	if count >= 3 {
		bonus = 20
	}

	if color == 1 {
		bonus = -bonus
	}

	return int16(bonus)
}

package main

func pawnBonus(rank int, file int, isWhitePiece bool) int16 {

	var eval int16 = 0

	if rank == 5 {
		eval += 10
	}

	if rank == 6 {
		eval += 30
	}

	if rank == 7 {
		eval += 80
	}

	if !isWhitePiece {
		eval = eval * -1
	}

	return eval

}

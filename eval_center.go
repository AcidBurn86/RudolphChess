package main

var centerMask Bitboard = 0x0000001818000000

var extendedCenterMask Bitboard = 0x0000003C3C000000

/*
	type gameTurnState struct {
		isWhiteTurn bool
		allWhites   Bitboard
		allBlacks   Bitboard
	}

// Mapas de caché para almacenar los resultados de las funciones.
var (

	cacheEvalCenterControl      = make(map[gameTurnState]int16)
	cacheValidateCenterControl2 = make(map[gameTurnState]int16)

)
*/
func evalCenterControl(isWhiteTurn bool) (result int16) {
	/*state := gameTurnState{isWhiteTurn: isWhiteTurn, allWhites: allWhites, allBlacks: allBlacks}
	if result, found := cacheEvalCenterControl[state]; found {
		return result // Retorna el resultado desde la caché si está presente.
	}*/
	var friendliesBitboard, enemiesBitboard Bitboard
	var friendlies, enemies int16

	if isWhiteTurn {
		friendliesBitboard = allWhites & centerMask
		enemiesBitboard = allBlacks & centerMask
	} else {
		friendliesBitboard = allBlacks & centerMask
		enemiesBitboard = allWhites & centerMask
	}

	friendlies = countBits(friendliesBitboard)
	enemies = countBits(enemiesBitboard)
	result = friendlies + enemies

	//cacheEvalCenterControl[state] = result

	return result
}

//var mutex sync.Mutex

func validateCenterControl2(isWhiteTurn bool) (result int16) {
	/*state := gameTurnState{isWhiteTurn: isWhiteTurn, allWhites: allWhites, allBlacks: allBlacks}

	mutex.Lock()
	result, found := cacheValidateCenterControl2[state]
	mutex.Unlock()

	if found {
		//return result // Retorna el resultado desde la caché si está presente.
	}*/

	var friendliesBitboard, enemiesBitboard Bitboard
	var friendlies, enemies int16

	if isWhiteTurn {
		friendliesBitboard = allWhites & extendedCenterMask
		enemiesBitboard = allBlacks & extendedCenterMask
	} else {
		friendliesBitboard = allBlacks & extendedCenterMask
		enemiesBitboard = allWhites & extendedCenterMask
	}

	friendlies = countBits(friendliesBitboard)
	enemies = countBits(enemiesBitboard)
	result = friendlies + enemies

	/*mutex.Lock()
	cacheValidateCenterControl2[state] = result
	mutex.Unlock()*/

	return result
}

func countBits(b Bitboard) int16 {
	var count int16
	count = 0
	for b != 0 {
		count += int16(b & 1)
		b >>= 1
	}
	return count
}

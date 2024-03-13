package main

import (
	"math/bits"
	"sync"
)

const (
	white = 0
	black = 1
)

var enemyPiecesNearPenalty = []int{30, 45, 60, 75, 90, 95, 105, 115, 150, 150, 150, 150, 150, 150, 150, 150}

func kingSafety(color int) int16 {
	var kingPosition int

	var score int16 = 0
	var penaltys int16 = 0

	if color == 0 {
		kingPosition = bits.TrailingZeros64(uint64(whiteKings))
	} else {
		kingPosition = bits.TrailingZeros64(uint64(blackKings))
	}

	//fmt.Println("King Position: ", kingPosition)

	var bishops Bitboard
	if color == 1 {
		bishops = whiteBishops
	} else {
		bishops = blackBishops
	}

	bishopXrays := attackingBishops(kingPosition, color, bishops)

	if bishopXrays != 0 {
		penaltys += xRayBishopPenalty
	}

	if color == 0 && penaltys != 0 {
		score = score - penaltys
	} else if color == 1 && penaltys != 0 {
		score = score + penaltys
	}

	enemyPiecesNear := enemyPiecesNearKing(kingPosition, color)

	if color == 0 && enemyPiecesNear > 0 {
		score = score - int16(enemyPiecesNearPenalty[enemyPiecesNear-1])
	} else if color == 1 && enemyPiecesNear > 0 {
		score = score + int16(enemyPiecesNearPenalty[enemyPiecesNear-1])
	}

	friendlyPiecesNear := friendlyPiecesNearKing(kingPosition, color)

	if color == 0 && friendlyPiecesNear > 0 {
		score = score + int16(enemyPiecesNearPenalty[friendlyPiecesNear-1])
	} else if color == 1 && friendlyPiecesNear > 0 {

		//	fmt.Println("Score: ", score)
		//	fmt.Println(int16(enemyPiecesNearPenalty[friendlyPiecesNear-1]))
		score = score - int16(enemyPiecesNearPenalty[friendlyPiecesNear-1])
	}

	// check pawn shield
	var pawnShieldScore int16 = 0
	if color == 0 {
		pawnShieldScore = int16(pawnShield(whitePawns, kingPosition))
		if pawnShieldScore == 3 {
			score = score + pawnShieldScore_1
		} else if pawnShieldScore == 2 {
			score = score + pawnShieldScore_2
		} else if pawnShieldScore == 1 {
			score = score - pawnShieldScore_3
		} else if pawnShieldScore == 0 {
			score = score - pawnShieldScore_4
		}
	} else {
		pawnShieldScore = int16(pawnShield(blackPawns, kingPosition))
		if pawnShieldScore == 3 {
			score = score - pawnShieldScore_1
		}
		if pawnShieldScore == 2 {
			score = score - pawnShieldScore_2
		}
		if pawnShieldScore == 1 {
			score = score + pawnShieldScore_3
		}
		if pawnShieldScore == 0 {
			score = score + pawnShieldScore_4
		}

	}
	///fmt.Println("Score: ", score)

	return score
}

func pawnShield(pawns Bitboard, kingPosition int) int {
	var shield int = 0
	var rank int = kingPosition / 8
	var file int = kingPosition % 8

	if rank < 6 {
		return 0
	}

	if file > 0 && pawns&(1<<(kingPosition-1)) != 0 {
		shield++
	}

	if pawns&(1<<kingPosition) != 0 {
		shield++
	}

	if file < 7 && pawns&(1<<(kingPosition+1)) != 0 {
		shield++
	}

	return shield
}

type attackKey struct {
	kingPosition int
	color        int
	bishops      Bitboard
}

// Estructura para memoización con mutex para acceso concurrente seguro.
type memoizeAttacks struct {
	sync.Mutex
	cache map[attackKey]Bitboard
}

var memoizedAttacks = memoizeAttacks{cache: make(map[attackKey]Bitboard)}

// attackingBishops memoizada.
func attackingBishops(kingPosition int, color int, bishops Bitboard) Bitboard {
	key := attackKey{kingPosition: kingPosition, color: color, bishops: bishops}

	memoizedAttacks.Lock()
	defer memoizedAttacks.Unlock()

	// Verificar si el resultado ya está en la caché.
	if result, found := memoizedAttacks.cache[key]; found {
		return result
	}

	var attacks Bitboard

	for bishops != 0 {
		// Encuentra la posición de un alfil.
		bishopPosition := bits.TrailingZeros64(uint64(bishops))

		//fmt.Println("Bishop Position: ", bishopPosition)

		bishopAttacks := bishopAttackTable[bishopPosition]

		// Si el alfil puede atacar la posición del rey, añade su posición a attacks.
		if bishopAttacks&(1<<kingPosition) != 0 {
			attacks |= (1 << bishopPosition)
		}

		// Elimina este alfil del conjunto y continúa con el siguiente.
		bishops &= bishops ^ (1 << bishopPosition)
	}

	return attacks
}

type cacheKey struct {
	kingPosition int
	color        int
}

// Estructuras de caché para cada función
var (
	enemyCache    = make(map[cacheKey]int)
	friendlyCache = make(map[cacheKey]int)
	cacheMutex    sync.Mutex // Para manejar el acceso concurrente a la caché
)

func enemyPiecesNearKing(kingPosition int, color int) int {
	key := cacheKey{kingPosition, color}

	// Revisar si el resultado ya está en la caché
	cacheMutex.Lock()
	result, found := enemyCache[key]
	cacheMutex.Unlock()
	if found {
		return result
	}

	var enemyPieces Bitboard
	if color == 1 {
		enemyPieces = allWhites // Piezas enemigas son las blancas
	} else {
		enemyPieces = allBlacks // Piezas enemigas son las negras
	}

	// Calcular las casillas alrededor del rey hasta dos casillas de distancia
	kingSurroundings := calculateKingSurroundings(kingPosition, 2)

	// Contar cuántas de estas casillas están ocupadas por piezas enemigas
	count := 0
	for _, pos := range kingSurroundings {
		if enemyPieces&(1<<pos) != 0 {
			count++
		}
	}

	cacheMutex.Lock()
	enemyCache[key] = count
	cacheMutex.Unlock()

	return count
}

func friendlyPiecesNearKing(kingPosition int, color int) int {
	key := cacheKey{kingPosition, color}

	// Revisar si el resultado ya está en la caché
	cacheMutex.Lock()
	result, found := friendlyCache[key]
	cacheMutex.Unlock()
	if found {
		return result
	}

	var friendlyPieces Bitboard
	if color == 0 {
		friendlyPieces = allWhites // Piezas amigas son las blancas
	} else {
		friendlyPieces = allBlacks // Piezas amigas son las negras
	}

	// Calcular las casillas alrededor del rey hasta dos casillas de distancia
	kingSurroundings := calculateKingSurroundings(kingPosition, 2)
	//fmt.Println("Color: ", color)
	// Contar cuántas de estas casillas están ocupadas por piezas amigas
	count := 0
	for _, pos := range kingSurroundings {
		if friendlyPieces&(1<<pos) != 0 {
			count++
		}
	}

	//	fmt.Println("Color: ", color, "Count: ", count)

	cacheMutex.Lock()
	friendlyCache[key] = count
	cacheMutex.Unlock()

	return count
}

type kingPositionKey struct {
	kingPosition int
	howFar       int
}

// Actualizar memoizeMap para usar kingPositionKey como clave.
type memoizeMap struct {
	sync.Mutex
	cache map[kingPositionKey][]int
}

// Inicializar la estructura con el tipo correcto para la mapa.
var memoizedResults = memoizeMap{cache: make(map[kingPositionKey][]int)}

func calculateKingSurroundings(kingPosition int, howFar int) []int {
	key := kingPositionKey{kingPosition, howFar}

	memoizedResults.Lock()
	defer memoizedResults.Unlock()

	// Verificar si el resultado ya está en la caché usando la nueva clave.
	if result, found := memoizedResults.cache[key]; found {
		return result
	}

	// Lógica de la función si el resultado no está en la caché.
	var surroundings []int
	rank, file := kingPosition/8, kingPosition%8

	for r := rank - howFar; r <= rank+howFar; r++ {
		for f := file - howFar; f <= file+howFar; f++ {
			if r >= 0 && r <= 7 && f >= 0 && f <= 7 {
				surroundings = append(surroundings, r*8+f)
			}
		}
	}

	// Almacenar el resultado calculado en la caché antes de retornarlo.
	memoizedResults.cache[key] = surroundings

	return surroundings
}

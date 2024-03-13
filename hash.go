package main

import (
	"hash/crc64"
	"sync"
)

var positions sync.Map // Mapa concurrente para almacenar el score por hash CRC64

// hashFENCRC64 genera un hash CRC64 de una cadena FEN.
func hashFENCRC64(fen string) uint64 {
	table := crc64.MakeTable(crc64.ECMA)
	hasher := crc64.New(table)
	hasher.Write([]byte(fen))
	return hasher.Sum64()
}

// savePositionScore guarda el score de una posición usando el hash CRC64 como clave.
func savePositionScore(fen string, score int) {
	hash := hashFENCRC64(fen)
	positions.Store(hash, score)
}

// getPositionScore recupera el score de una posición por su cadena FEN.
func getPositionScore(fen string) (int, bool) {
	hash := hashFENCRC64(fen)
	if score, ok := positions.Load(hash); ok {
		return score.(int), true
	}
	return 0, false
}

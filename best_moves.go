package main

import (
	"sort"
	"strings"
)

type claveValor struct {
	Clave string
	Valor int16
}

func AgregarElemento(clave string, valor int16) {
	mu.Lock()
	defer mu.Unlock()
	miMapa[clave] = valor
}

func BuscarPorPrefijo(prefijo string) []claveValor {
	mu.Lock()
	defer mu.Unlock()

	var pares []claveValor
	for clave, valor := range miMapa {
		if strings.HasPrefix(clave, prefijo) {
			pares = append(pares, claveValor{Clave: clave, Valor: valor})
		}
	}

	sort.Slice(pares, func(i, j int) bool {
		return pares[i].Valor > pares[j].Valor
	})

	if len(pares) > 3 {
		pares = pares[:3]
	}

	return pares
}

func EncontrarPrefijosUnicos(chunck int) []string {
	mu.Lock()
	defer mu.Unlock()

	chunck = int(chunck)

	prefijos := make(map[string]struct{})
	for clave := range miMapa {
		if len(clave) >= chunck {
			prefijo := clave[:chunck]
			prefijos[prefijo] = struct{}{}
		}
	}

	var resultado []string
	for prefijo := range prefijos {
		resultado = append(resultado, prefijo)
	}

	return resultado
}

type keyValue struct {
	Key   string
	Value int16
}

func reProcessBestLines(score int16, depth int8, isWhiteTurn bool) []keyValue {
	// Convertir el mapa a un slice para poder ordenarlo
	var pares []keyValue
	for clave, valor := range miMapa {
		pares = append(pares, keyValue{clave, valor})
	}

	// Ordenar el slice por valor de manera descendiente
	if isWhiteTurn {
		sort.Slice(pares, func(i, j int) bool {
			return pares[i].Value > pares[j].Value
		})
	} else {
		sort.Slice(pares, func(i, j int) bool {
			return pares[i].Value < pares[j].Value
		})
	}

	// Limpiar el mapa y quedarse solo con los 4 valores más altos
	var miMapa2 []keyValue
	for i, par := range pares {
		if i < 8 { // Tomar solo los 4 primeros elementos
			miMapa2 = append(miMapa2, par)
		} else {
			break
		}
	}

	return miMapa2
}

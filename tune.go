package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

type Test struct {
	FEN           string
	GoCmd         []string
	ExpectedMove  string
	ExpectedNodes int
}

func tune() int {
	tests := []Test{
		/*{FEN: "2r4r/pp5p/2n2pp1/7N/3P2k1/2P5/PP3P1R/R1B1K1N1 w Q - 1 22", GoCmd: []string{"fen", "2r4r/pp5p/2n2pp1/7N/3P2k1/2P5/PP3P1R/R1B1K1N1 w Q - 1 22", "go", "depth", "5"}, ExpectedMove: "f2f3", ExpectedNodes: 559},
		{FEN: "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 11", GoCmd: []string{"fen", "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 10", "go", "depth", "6"}, ExpectedMove: "e2a6", ExpectedNodes: 231816},
		{FEN: "3r1rk1/pp5p/2p3p1/8/2bNR3/P5P1/5P1P/4R1K1 w - - 1 25", GoCmd: []string{"fen", "3r1rk1/pp5p/2p3p1/8/2bNR3/P5P1/5P1P/4R1K1 w - - 1 25", "go", "depth", "6"}, ExpectedMove: "d4c6", ExpectedNodes: 674},
		{FEN: "r2q1rk1/ppp1bppp/8/3p4/b2PpPP1/N1P4N/R3Q2P/2n2RK1 w - - 0 17", GoCmd: []string{"fen", "r2q1rk1/ppp1bppp/8/3p4/b2PpPP1/N1P4N/R3Q2P/2n2RK1 w - - 0 17", "go", "depth", "6"}, ExpectedMove: "f1c1", ExpectedNodes: 163},*/
		/*{FEN: "r1b1k2r/ppppnppp/4p3/1PP5/3Pn1q1/P1N2N2/4PP1P/R2QKB1R w KQkq - 1 12", GoCmd: []string{"fen", "r1b1k2r/ppppnppp/4p3/1PP5/3Pn1q1/P1N2N2/4PP1P/R2QKB1R w KQkq - 1 12", "go", "depth", "5"}, ExpectedMove: "c3e4", ExpectedNodes: 113488},
		{FEN: "r1b1kb1r/2qpp3/5npp/pp6/5n2/P1P5/NP2P1RN/R1B1K3 w Qkq - 0 18", GoCmd: []string{"fen", "r1b1kb1r/2qpp3/5npp/pp6/5n2/P1P5/NP2P1RN/R1B1K3 w Qkq - 0 18", "go", "depth", "5"}, ExpectedMove: "c1f4", ExpectedNodes: 101513},*/
		{FEN: "r1bqk2r/1pp2ppp/p1n2n2/P2Pp3/2BP4/2P5/P4PPP/RNBQK1NR b KQkq - 0 8", GoCmd: []string{"fen", "r1bqk2r/1pp2ppp/p1n2n2/P2Pp3/2BP4/2P5/P4PPP/RNBQK1NR b KQkq - 0 8", "go", "depth", "5"}, ExpectedMove: "c8g4", ExpectedNodes: 144046},

		{FEN: "rnbqk2r/1ppp1ppp/p4n2/P3p3/2B1P3/2P5/P2P1PPP/RNBQK1NR b KQkq - 0 6", GoCmd: []string{"fen", "rnbqk2r/1ppp1ppp/p4n2/P3p3/2B1P3/2P5/P2P1PPP/RNBQK1NR b KQkq - 0 6", "go", "depth", "5"}, ExpectedMove: "b8c6", ExpectedNodes: 144046},
		{FEN: "rnbqk1nr/1ppp1ppp/p7/b3p3/1PB1P3/2P5/P2P1PPP/RNBQK1NR b KQkq - 0 5", GoCmd: []string{"fen", "rnbqk1nr/1ppp1ppp/p7/b3p3/1PB1P3/2P5/P2P1PPP/RNBQK1NR b KQkq - 0 5", "go", "depth", "5"}, ExpectedMove: "a5b6", ExpectedNodes: 144046},
	}

	var results []TestResult
	correctPositions := 0 // Contador para las posiciones correctas
	incorrect := 0

	for i, test := range tests {
		fmt.Printf("## Test %d\n", i+1)
		handleUCINewGame()
		handlePosition([]string{"fen", test.FEN})
		result := handleGo(test.GoCmd, test.ExpectedMove, test.ExpectedNodes)
		results = append(results, result)
		if result.IsCorrect {
			correctPositions++
		} else {
			incorrect++
		}
		if incorrect >= 1 {
			break

		}
	}

	for i, result := range results {
		var nodeDiff string

		if result.NodeCount > result.ExpectedNodes {
			// Si se evaluaron más nodos de los esperados, muestra en rojo
			nodeDiff = fmt.Sprintf("%s%d%s", Red, result.NodeCount, Reset)
		} else if result.NodeCount == result.ExpectedNodes {
			// Si se evaluaron menos nodos de los esperados, muestra en verde
			nodeDiff = fmt.Sprintf("%s%d%s", Green, result.NodeCount, Reset)
		} else if result.NodeCount < result.ExpectedNodes {
			// Si se evaluaron menos nodos de los esperados, muestra en verde
			nodeDiff = fmt.Sprintf("%s%d%s", lightBlue, result.NodeCount, Reset)
		} else {
			// Si es exactamente igual, sin color
			nodeDiff = fmt.Sprintf("%d", result.NodeCount)
		}

		// Decide el color para el movimiento encontrado basado en si es correcto o no
		moveColor := ""
		if !result.IsCorrect {
			moveColor = fmt.Sprintf("%s%s%s", Red, result.FoundMove, Reset)
		} else {
			moveColor = fmt.Sprintf("%s%s%s", Green, result.FoundMove, Reset)
		}

		fmt.Printf("Test #%d: Expected: %s, Found: %s, Correct: %t, Nodes: %s, Expected Nodes: %d, Time: %s\n",
			i+1, result.ExpectedMove, moveColor, result.IsCorrect, nodeDiff, result.ExpectedNodes, result.ElapsedTime)
	}

	return correctPositions

}

func adjustVariablesRandomly() {
	rand.Seed(time.Now().UnixNano())                       // Inicializa la semilla del generador de números aleatorios
	threshold_isMiddleGame = int16(rand.Intn(200)+100) * 2 // Ejemplo: rango 100-300, ajustado para ser par
	threshold_isOpening = int16(rand.Intn(50)+2) * 2       // Ejemplo: rango 50-150, ajustado para ser par
	threshold_isEndgame = int16(rand.Intn(30)+5) * 2       // Ejemplo: rango 20-80, ajustado para ser par
	edgeScore_value = int16(rand.Intn(100)+75) * 2         // Ejemplo: rango 150-350, ajustado para ser par

	pawnShieldScore_1 = int16(rand.Intn(25)+25) * 2 // Ejemplo: rango 60-100, ajustado para ser par
	pawnShieldScore_2 = int16(rand.Intn(20)+10) * 2 // Ejemplo: rango 30-70, ajustado para ser par
	pawnShieldScore_3 = int16(rand.Intn(20)+10) * 2 // Ejemplo: rango 10-30, ajustado para ser par
	pawnShieldScore_4 = int16(rand.Intn(20)+5) * 2  // Ejemplo: rango 30-70, ajustado para ser par

	xRayBishopPenalty = int16(rand.Intn(10)+5) * 2 // Ejemplo: rango 10-30, ajustado para ser par

	RookInOpenRank = int16(rand.Intn(20)+25) * 2 // Ejemplo: rango 20-60, ajustado para ser par

	kingSafetyMulti = int16(rand.Intn(3)+1) * 2 // Ejemplo: rango 2-6, ajustado para ser par, nota que el rango es reducido

	center_multi_1 = int16(rand.Intn(15)+2) * 2 // Ejemplo: rango 4-24, ajustado para ser par
	center_multi_2 = int16(rand.Intn(2)+6) * 2  // Ejemplo: rango 2-6, ajustado para ser par

	queenInStartPosition = int16(rand.Intn(100)+50) * 2 // Ejemplo: rango 100-300, ajustado para ser par

	bishopPairBonus = int16(rand.Intn(20)+5) * 2 // Ejemplo: rango 10-50, ajustado para ser par

	PawnsDefending_Multi = int(rand.Intn(10)+5) * 2 // Ejemplo: rango 10-30, ajustado para ser par

	doubledPawnPenalty = int16(rand.Intn(15)+10) * 2 // Ejemplo: rango 20-70, ajustado para ser par

	mobilityPenalty = int16(rand.Intn(30)+10) * 2 // Ejemplo: rango 20-100, ajustado para ser par
}

func printVariableValues() {
	fmt.Println("threshold_isOpening:", threshold_isOpening)
	fmt.Println("threshold_isEndgame:", threshold_isEndgame)
	fmt.Println("edgeScore_value:", edgeScore_value)
	fmt.Println("pawnShieldScore_1:", pawnShieldScore_1)
	fmt.Println("pawnShieldScore_2:", pawnShieldScore_2)
	fmt.Println("pawnShieldScore_3:", pawnShieldScore_3)
	fmt.Println("pawnShieldScore_4:", pawnShieldScore_4)
	fmt.Println("xRayBishopPenalty:", xRayBishopPenalty)
	fmt.Println("RookInOpenRank:", RookInOpenRank)
	fmt.Println("kingSafetyMulti:", kingSafetyMulti)
	fmt.Println("center_multi_1:", center_multi_1)
	fmt.Println("center_multi_2:", center_multi_2)
	fmt.Println("queenInStartPosition:", queenInStartPosition)
	fmt.Println("bishopPairBonus:", bishopPairBonus)
	fmt.Println("PawnsDefending_Multi:", PawnsDefending_Multi)
	fmt.Println("doubledPawnPenalty:", doubledPawnPenalty)
	fmt.Println("mobilityPenalty:", mobilityPenalty)
}

func startTune() {
	hide_output = true
	correctPositions := 0
	iterationCount := 0
	tests := 3

	for correctPositions < (tests) { // Asumiendo que 'tests' es el slice de pruebas
		adjustVariablesRandomly() // Ajusta las variables aleatoriamente
		correctPositions = tune() // Ejecuta la función de pruebas y obtén el número de posiciones correctas
		iterationCount++

		if correctPositions >= 3 {
			saveVariableValuesToFile(correctPositions, iterationCount) // Guarda los valores de las variables en un archivo
		}
		fmt.Println("Iteración:", iterationCount)

		// Si todas las posiciones son correctas, imprime los valores de las variables
		if correctPositions == (tests) {
			fmt.Printf("Todas las posiciones evaluadas correctamente después de %d iteraciones.\n", iterationCount)
			printVariableValues() // Función que imprime los valores de las variables
			break
		}
	}
}

func saveVariableValuesToFile(correctPositions int, iterationCount int) {
	// Crear el nombre del archivo basado en la cantidad de aciertos y el número de la iteración
	//get unix time
	t := time.Now()

	fileName := fmt.Sprintf("%d_aciertos_iteracion_%d_%d.txt", correctPositions, iterationCount, t.Unix())

	// Crear o abrir el archivo para escritura
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error al crear el archivo:", err)
		return
	}
	defer file.Close()

	// Escribir los valores de las variables en el archivo
	writeVariableToFile(file, "threshold_isOpening", threshold_isOpening)
	writeVariableToFile(file, "threshold_isEndgame", threshold_isEndgame)
	writeVariableToFile(file, "edgeScore_value", edgeScore_value)
	writeVariableToFile(file, "pawnShieldScore_1", pawnShieldScore_1)
	writeVariableToFile(file, "pawnShieldScore_2", pawnShieldScore_2)
	writeVariableToFile(file, "pawnShieldScore_3", pawnShieldScore_3)
	writeVariableToFile(file, "pawnShieldScore_4", pawnShieldScore_4)
	writeVariableToFile(file, "xRayBishopPenalty", xRayBishopPenalty)
	writeVariableToFile(file, "RookInOpenRank", RookInOpenRank)
	writeVariableToFile(file, "kingSafetyMulti", kingSafetyMulti)
	writeVariableToFile(file, "center_multi_1", center_multi_1)
	writeVariableToFile(file, "center_multi_2", center_multi_2)
	writeVariableToFile(file, "queenInStartPosition", queenInStartPosition)
	writeVariableToFile(file, "bishopPairBonus", bishopPairBonus)
	writeVariableToFile(file, "PawnsDefending_Multi", PawnsDefending_Multi)
	writeVariableToFile(file, "doubledPawnPenalty", doubledPawnPenalty)
	writeVariableToFile(file, "mobilityPenalty", mobilityPenalty)

	// Continúa escribiendo el resto de las variables aquí...

	fmt.Printf("Valores guardados en el archivo '%s'\n", fileName)
}

// Función auxiliar para escribir cada variable en el archivo
func writeVariableToFile(file *os.File, variableName string, value interface{}) {
	_, err := file.WriteString(fmt.Sprintf("%s: %v\n", variableName, value))
	if err != nil {
		fmt.Println("Error al escribir en el archivo:", err)
	}
}

package main

import (
	"fmt"
)

const (
	Green     = "\033[32m"
	lightBlue = "\033[36m"
	Red       = "\033[31m"
	Reset     = "\033[0m"
)

func t1() {
	var results []TestResult
	fmt.Println("## Test 1")
	handleUCINewGame()
	handlePosition([]string{"fen", "2r4r/pp5p/2n2pp1/7N/3P2k1/2P5/PP3P1R/R1B1K1N1 w Q - 1 22"})
	results = append(results, handleGo([]string{"fen", "2r4r/pp5p/2n2pp1/7N/3P2k1/2P5/PP3P1R/R1B1K1N1 w Q - 1 22", "go", "depth", "1"}, "f2f3", 559))
	fmt.Println("## Test 2")
	handleUCINewGame()
	handlePosition([]string{"fen", "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 11"})
	results = append(results, handleGo([]string{"fen", "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 10", "go", "depth", "1"}, "c3b5", 231816))

	fmt.Println("## Test 3")
	handleUCINewGame()
	handlePosition([]string{"fen", "3r1rk1/pp5p/2p3p1/8/2bNR3/P5P1/5P1P/4R1K1 w - - 1 25"})
	results = append(results, handleGo([]string{"fen", "3r1rk1/pp5p/2p3p1/8/2bNR3/P5P1/5P1P/4R1K1 w - - 1 25", "go", "depth", "1"}, "d4c6", 557140))

	fmt.Println("## Test 4")
	handleUCINewGame()
	handlePosition([]string{"fen", "3Rr1k1/p4ppp/1p6/2p1q3/2P5/1P4P1/r4P1P/3R1QK1 w - - 0 15"})
	results = append(results, handleGo([]string{"fen", "3Rr1k1/p4ppp/1p6/2p1q3/2P5/1P4P1/r4P1P/3R1QK1 w - - 0 15", "go", "depth", "1"}, "d1e1", 1515729))
	fmt.Println("## Test 5")
	handleUCINewGame()
	handlePosition([]string{"fen", "3r1r2/ppp4k/1bn1Q1pq/8/5P2/1P4P1/PBP3BP/3R3K w - - 3 24"})
	results = append(results, handleGo([]string{"fen", "3r1r2/ppp4k/1bn1Q1pq/8/5P2/1P4P1/PBP3BP/3R3K w - - 3 24", "go", "depth", "1"}, "d1d7", 35529))
	fmt.Println("## Test 6")
	handleUCINewGame()
	handlePosition([]string{"fen", "r1b2rk1/ppp5/1bn1N1pq/8/4QP2/1P4P1/PBP3BP/1R5K b - - 0 21"})
	results = append(results, handleGo([]string{"fen", "r1b2rk1/ppp5/1bn1N1pq/8/4QP2/1P4P1/PBP3BP/1R5K b - - 0 21", "go", "depth", "1"}, "f8f5", 334277))
	fmt.Println("## Test 7")
	handleUCINewGame()
	handlePosition([]string{"fen", "1q2k3/2pn1p1p/p1Qbp1r1/8/2PB4/4NP2/P2P1P1P/4K2R w K - 0 22"})
	results = append(results, handleGo([]string{"fen", "1q2k3/2pn1p1p/p1Qbp1r1/8/2PB4/4NP2/P2P1P1P/4K2R w K - 0 22", "go", "depth", "1"}, "e1e2", 280472))

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

}

func testsPos() {
	var results []TestResult
	fmt.Println("## Test 1")
	handleUCINewGame()
	handlePosition([]string{"fen", "2r4r/pp5p/2n2pp1/7N/3P2k1/2P5/PP3P1R/R1B1K1N1 w Q - 1 22"})
	results = append(results, handleGo([]string{"fen", "2r4r/pp5p/2n2pp1/7N/3P2k1/2P5/PP3P1R/R1B1K1N1 w Q - 1 22", "go", "depth", "5"}, "f2f3", 1032))
	fmt.Println("## Test 2")
	handleUCINewGame()
	handlePosition([]string{"fen", "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 11"})
	results = append(results, handleGo([]string{"fen", "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 10", "go", "depth", "5"}, "e2a6", 61970))

	fmt.Println("## Test 3")
	handleUCINewGame()
	handlePosition([]string{"fen", "3r1rk1/pp5p/2p3p1/8/2bNR3/P5P1/5P1P/4R1K1 w - - 1 25"})
	results = append(results, handleGo([]string{"fen", "3r1rk1/pp5p/2p3p1/8/2bNR3/P5P1/5P1P/4R1K1 w - - 1 25", "go", "depth", "5"}, "d4c6", 92128))

	/*fmt.Println("## Test 4")
	handleUCINewGame()
	handlePosition([]string{"fen", "3Rr1k1/p4ppp/1p6/2p1q3/2P5/1P4P1/r4P1P/3R1QK1 w - - 0 15"})
	results = append(results, handleGo([]string{"fen", "3Rr1k1/p4ppp/1p6/2p1q3/2P5/1P4P1/r4P1P/3R1QK1 w - - 0 15", "go", "depth", "5"}, "d1d5", 29530))
	fmt.Println("## Test 5")
	handleUCINewGame()
	handlePosition([]string{"fen", "3r1r2/ppp4k/1bn1Q1pq/8/5P2/1P4P1/PBP3BP/3R3K w - - 3 24"})
	results = append(results, handleGo([]string{"fen", "3r1r2/ppp4k/1bn1Q1pq/8/5P2/1P4P1/PBP3BP/3R3K w - - 3 24", "go", "depth", "5"}, "e6d7", 28103))
	fmt.Println("## Test 6")
	handleUCINewGame()
	handlePosition([]string{"fen", "r1b2rk1/ppp5/1bn1N1pq/8/4QP2/1P4P1/PBP3BP/1R5K b - - 0 21"})
	results = append(results, handleGo([]string{"fen", "r1b2rk1/ppp5/1bn1N1pq/8/4QP2/1P4P1/PBP3BP/1R5K b - - 0 21", "go", "depth", "5"}, "f8e8", 38831))
	fmt.Println("## Test 7")
	handleUCINewGame()
	handlePosition([]string{"fen", "1q2k3/2pn1p1p/p1Qbp1r1/8/2PB4/4NP2/P2P1P1P/4K2R w K - 0 22"})
	results = append(results, handleGo([]string{"fen", "1q2k3/2pn1p1p/p1Qbp1r1/8/2PB4/4NP2/P2P1P1P/4K2R w K - 0 22", "go", "depth", "5"}, "e1e2", 39427))*/

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

}

func testsPos3() {
	var results []TestResult

	//isTest2 = true

	handleUCINewGame()
	handlePosition([]string{"fen", "2r1kq2/p7/Q2bb1N1/2pNpp2/8/4P3/PPPP1PPP/R1B1K2R b KQ - 2 26"})
	results = append(results, handleGo([]string{"fen", "2r1kq2/p7/Q2bb1N1/2pNpp2/8/4P3/PPPP1PPP/R1B1K2R b KQ - 2 26", "go", "depth", "5"}, "f8g7", 557140))

	handleUCINewGame()
	handlePosition([]string{"fen", "rnbqkb1r/pppp1ppp/5n2/4p3/4P1Q1/1P6/P1PP1PPP/RNB1KBNR b KQkq - 0 3"})
	results = append(results, handleGo([]string{"fen", "rnbqkb1r/pppp1ppp/5n2/4p3/4P1Q1/1P6/P1PP1PPP/RNB1KBNR b KQkq - 0 3", "go", "depth", "5"}, "f6g4", 557140))

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
}

func testsPos2() {
	var results []TestResult

	//isTest2 = true

	handleUCINewGame()
	handlePosition([]string{"fen", "r1bq2kr/1p2ppb1/n1p1p1pn/p5Bp/2PP1Q2/2NB1N2/PP3PPP/R4RK1 w - - 0 13"})
	results = append(results, handleGo([]string{"fen", "r1bq2kr/1p2ppb1/n1p1p1pn/p5Bp/2PP1Q2/2NB1N2/PP3PPP/R4RK1 w - - 0 13", "go", "depth", "5"}, "c1f4", 539860))

	/*handleUCINewGame()
	handlePosition([]string{"fen", "r3kbnr/p1pp1ppp/bp2p3/3Pn1Bq/2P5/5N2/PPB2PPP/RN1Q1RK1 w kq - 0 11"})
	results = append(results, handleGo([]string{"fen", "r3kbnr/p1pp1ppp/bp2p3/3Pn1Bq/2P5/5N2/PPB2PPP/RN1Q1RK1 w kq - 0 11", "go", "depth", "6"}, "f3e5", 557140))

	handleUCINewGame()
	handlePosition([]string{"fen", "3Rr1k1/p4ppp/1p6/2p1q3/2P5/1P4P1/r4P1P/3R1QK1 w - - 0 15"})
	results = append(results, handleGo([]string{"fen", "3Rr1k1/p4ppp/1p6/2p1q3/2P5/1P4P1/r4P1P/3R1QK1 w - - 0 15", "go", "depth", "7"}, "d1d5", 1515729))

	handleUCINewGame()
	handlePosition([]string{"fen", "3r1r2/ppp4k/1bn1Q1pq/8/5P2/1P4P1/PBP3BP/3R3K w - - 3 24"})
	results = append(results, handleGo([]string{"fen", "3r1r2/ppp4k/1bn1Q1pq/8/5P2/1P4P1/PBP3BP/3R3K w - - 3 24", "go", "depth", "6"}, "d1d7", 116712))

	handleUCINewGame()
	handlePosition([]string{"fen", "r1b2rk1/ppp5/1bn1N1pq/8/4QP2/1P4P1/PBP3BP/1R5K b - - 0 21"})
	results = append(results, handleGo([]string{"fen", "r1b2rk1/ppp5/1bn1N1pq/8/4QP2/1P4P1/PBP3BP/1R5K b - - 0 21", "go", "depth", "6"}, "f8e8", 334277))

	handleUCINewGame()
	handlePosition([]string{"fen", "1q2k3/2pn1p1p/p1Qbp1r1/8/2PB4/4NP2/P2P1P1P/4K2R w K - 0 22"})
	results = append(results, handleGo([]string{"fen", "1q2k3/2pn1p1p/p1Qbp1r1/8/2PB4/4NP2/P2P1P1P/4K2R w K - 0 22", "go", "depth", "6"}, "e1e2", 280472))*/

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

}

func testsPosLI() {
	var results []TestResult
	fmt.Println("## Test 1")
	handleUCINewGame()
	handlePosition([]string{"fen", "1r3b1r/p1ppkppp/2p3n1/2P1P3/3KNP2/P7/1P4PP/R1B2R2 b - - 4 19"})
	results = append(results, handleGo([]string{"fen", "1r3b1r/p1ppkppp/2p3n1/2P1P3/3KNP2/P7/1P4PP/R1B2R2 b - - 4 19", "go", "depth", "5"}, "f7f5", 13727))
	fmt.Println("## Test 2")
	handleUCINewGame()
	handlePosition([]string{"fen", "r3kb1r/p1ppqppp/b1p5/3nP3/2P5/8/PP1NQPPP/R1B1KB1R b KQkq - 2 9"})
	results = append(results, handleGo([]string{"fen", "r3kb1r/p1ppqppp/b1p5/3nP3/2P5/8/PP1NQPPP/R1B1KB1R b KQkq - 2 9", "go", "depth", "5"}, "e1c1", 22799))

	fmt.Println("## Test 3")
	handleUCINewGame()
	handlePosition([]string{"fen", "4kb1r/p1pp1ppp/2p3n1/2P1P3/1r2NP2/3K4/PP4PP/R1B2R2 b k - 4 16"})
	results = append(results, handleGo([]string{"fen", "4kb1r/p1pp1ppp/2p3n1/2P1P3/1r2NP2/3K4/PP4PP/R1B2R2 b k - 4 16", "go", "depth", "5"}, "e8e7", 45185))

	fmt.Println("## Test 4")
	handleUCINewGame()
	handlePosition([]string{"fen", "r1bqkb1r/1pp2ppp/p1p5/4p3/P3n3/5N1P/1PPPQPP1/RNB1K2R b KQkq - 1 7"})
	results = append(results, handleGo([]string{"fen", "r1bqkb1r/1pp2ppp/p1p5/4p3/P3n3/5N1P/1PPPQPP1/RNB1K2R b KQkq - 1 7", "go", "depth", "5"}, "e4g5", 29434))
	fmt.Println("## Test 5")
	handleUCINewGame()
	handlePosition([]string{"fen", "r3k2r/1pp2ppp/p1p1b3/q3N3/PbPP4/6PP/1P1BQ1P1/RN3RK1 b kq - 2 13"})
	results = append(results, handleGo([]string{"fen", "r3k2r/1pp2ppp/p1p1b3/q3N3/PbPP4/6PP/1P1BQ1P1/RN3RK1 b kq - 2 13", "go", "depth", "5"}, "e1g1", 158851))
	fmt.Println("## Test 6")
	handleUCINewGame()
	handlePosition([]string{"fen", "r3kbnr/4pppp/p7/1ppq4/3n1Q2/6PN/PPPPRP2/RNB2K2 w kq - 4 18"})
	results = append(results, handleGo([]string{"fen", "r3kbnr/4pppp/p7/1ppq4/3n1Q2/6PN/PPPPRP2/RNB2K2 w kq - 4 18", "go", "depth", "5"}, "b1c3", 15692))
	fmt.Println("## Test 7")
	handleUCINewGame()
	handlePosition([]string{"fen", "3r1b1r/1pNbkppp/p2p1n2/4p3/P6Q/2R5/1PPB1PqP/4K1NR w K - 2 17"})
	results = append(results, handleGo([]string{"fen", "3r1b1r/1pNbkppp/p2p1n2/4p3/P6Q/2R5/1PPB1PqP/4K1NR w K - 2 17", "go", "depth", "5"}, "f2f3", 117672))

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

}
